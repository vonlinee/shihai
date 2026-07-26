package config

import (
	"encoding/json"
	"fmt"
	"shihai/internal/models"
	"shihai/internal/poetry"
	"shihai/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	// First, try to connect to postgres database to create our target database if it doesn't exist
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Port, cfg.SSLMode)

	postgresDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	// Create database if it doesn't exist
	createDBSQL := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
	if err := postgresDB.Exec(createDBSQL).Error; err != nil {
		// Database might already exist, which is fine
		// We can ignore this error or log it for debugging
	}

	// Close the postgres connection
	sqlDB, err := postgresDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	// Now connect to our target database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// AutoMigrateDatabaseModel Auto migrate models
func AutoMigrateDatabaseModel(db *gorm.DB, err error) error {
	if err := db.AutoMigrate(&models.Author{}); err != nil {
		return err
	}
	if err := migratePoemContentToJSONB(db); err != nil {
		return err
	}
	if err := migrateLegacyPoetSchema(db); err != nil {
		return err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Dynasty{},
		&models.PoemType{},
		&models.Author{},
		&models.Poet{},
		&models.Poem{},
		&models.PoemAnnotation{},
		&models.Comment{},
		&models.CommentVote{},
		&models.Quiz{},
		&models.QuizRecord{},
		&models.ForumPost{},
		&models.ForumReply{},
		&models.CorrectionRequest{},
		&models.CorrectionVote{},
		&models.Announcement{},
		&models.Feedback{},
		&models.OperationLog{},
		&models.WorkCollection{},
		&models.WorkCollectionItem{},
		// RBAC models
		&models.Role{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.Permission{},
	)
	if err != nil {
		return err
	}
	if err := backfillPoemPingze(db); err != nil {
		return err
	}
	if err := MigrateMissingPoetsFromPoems(db); err != nil {
		return err
	}

	// Set table comments
	if err := setTableComments(db); err != nil {
		return err
	}
	return err
}

// backfillPoemPingze 为历史诗词补齐与正文行数对应的平仄数组。
//
// db 数据库连接。新增 pingze 列后，历史诗词没有人工维护的平仄数据；这里只补空标记，
// 不按现代普通话自动推断，避免多音字和古入声导致错误数据入库。
func backfillPoemPingze(db *gorm.DB) error {
	for _, sql := range buildPoemPingzeBackfillSQL("poem") {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return backfillMissingPoemPingzeValues(db)
}

// buildPoemPingzeBackfillSQL 构造历史诗词平仄字段的回填 SQL。
//
// tableName 待回填的表名。返回 SQL 会在 pingze 为空或不是 JSON 数组时按 content 数组长度补空字符串。
func buildPoemPingzeBackfillSQL(tableName string) []string {
	return []string{
		fmt.Sprintf(`
UPDATE "%s"
SET "pingze" = COALESCE((
    SELECT jsonb_agg(to_jsonb(''::text) ORDER BY content_lines.ord)
    FROM jsonb_array_elements(COALESCE("content", '[]'::jsonb)) WITH ORDINALITY AS content_lines(line, ord)
), '[]'::jsonb)
WHERE "pingze" IS NULL
   OR jsonb_typeof("pingze") <> 'array';`, tableName),
	}
}

func backfillMissingPoemPingzeValues(db *gorm.DB) error {
	var poems []models.Poem
	return db.Model(&models.Poem{}).
		Select("id", "content", "pingze").
		FindInBatches(&poems, 100, func(tx *gorm.DB, batch int) error {
			for _, poem := range poems {
				if poetry.HasPingzeValue(poem.Pingze) {
					continue
				}
				pingze := poetry.RecognizePingzeLines(poem.Content)
				pingzeJSON, err := marshalPoemPingzeJSON(pingze)
				if err != nil {
					return err
				}
				if err := db.Model(&models.Poem{}).
					Where("id = ?", poem.ID).
					UpdateColumn("pingze", gorm.Expr("?::jsonb", pingzeJSON)).Error; err != nil {
					return err
				}
			}
			return nil
		}).Error
}

func marshalPoemPingzeJSON(pingze []string) (string, error) {
	data, err := json.Marshal(pingze)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// migratePoemContentToJSONB 在 GORM 自动迁移前修复旧库中的诗词正文列类型。
//
// db 数据库连接。旧版本将 poem.content 存为普通文本；直接交给 AutoMigrate 改成 jsonb 时，
// PostgreSQL 会把文本按 JSON 字面量解析，中文正文不是合法 JSON，因此会触发 22P02。
func migratePoemContentToJSONB(db *gorm.DB) error {
	for _, sql := range buildPoemContentJSONBMigrationSQL("poem") {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

// buildPoemContentJSONBMigrationSQL 构造诗词正文列的兼容迁移 SQL。
//
// tableName 待迁移的表名。返回的 SQL 会在列存在且不是 jsonb 时，将旧文本内容包装为单元素 JSON 数组。
func buildPoemContentJSONBMigrationSQL(tableName string) []string {
	return []string{
		fmt.Sprintf(`
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = '%s'
          AND column_name = 'content'
          AND data_type <> 'jsonb'
    ) THEN
        ALTER TABLE "%s"
            ALTER COLUMN content DROP DEFAULT,
            ALTER COLUMN content TYPE JSONB
            USING to_jsonb(ARRAY[content]),
            ALTER COLUMN content SET DEFAULT '[]'::jsonb;
    END IF;
END $$;`, tableName, tableName),
	}
}

// migrateLegacyPoetSchema 兼容旧版本 poet 表中的冗余作者字段。
//
// db 数据库连接。旧表曾直接在 poet 上保存 name、biography、avatar，当前模型已迁移到 author 表；
// 若这些历史列仍带 NOT NULL 约束，补齐 poet 扩展记录时会因未写入旧列而失败。
func migrateLegacyPoetSchema(db *gorm.DB) error {
	for _, sql := range buildLegacyPoetSchemaMigrationSQL("poet") {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

// buildLegacyPoetSchemaMigrationSQL 构造旧 poet 表冗余字段的兼容迁移 SQL。
//
// tableName 待迁移的表名。返回的 SQL 会在历史列存在时解除 NOT NULL 约束，使当前 Poet 模型能够只写扩展字段。
func buildLegacyPoetSchemaMigrationSQL(tableName string) []string {
	return []string{
		fmt.Sprintf(`
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = '%s'
          AND column_name = 'name'
          AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE "%s" ALTER COLUMN "name" DROP NOT NULL;
    END IF;
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = '%s'
          AND column_name = 'biography'
          AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE "%s" ALTER COLUMN "biography" DROP NOT NULL;
    END IF;
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = '%s'
          AND column_name = 'avatar'
          AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE "%s" ALTER COLUMN "avatar" DROP NOT NULL;
    END IF;
END $$;`, tableName, tableName, tableName, tableName, tableName, tableName),
	}
}

// MigrateMissingPoetsFromPoems 为历史诗词作者补齐诗人扩展记录。
//
// db 数据库连接。旧版同步脚本会导入 author 和 poem，但不会同步创建 poet 记录；这会导致后台诗人管理列表为空。
// 该迁移只为已经被 poem.author_id 引用且尚无 poet 记录的作者创建诗人扩展，不会把未被诗词引用的普通作者误加入诗人列表。
func MigrateMissingPoetsFromPoems(db *gorm.DB) error {
	for _, sql := range buildMissingPoetMigrationSQL() {
		if err := db.Exec(sql, utils.GenerateID()).Error; err != nil {
			return err
		}
	}
	return nil
}

// buildMissingPoetMigrationSQL 构造从诗词作者补齐诗人扩展记录的迁移 SQL。
//
// 返回的 SQL 会按 poem.author_id 去重，并跳过已经存在 poet 记录的作者。
func buildMissingPoetMigrationSQL() []string {
	return []string{
		`
INSERT INTO "poet" ("id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "author_id", "dynasty_id", "birth_year", "death_year")
SELECT
    ? + ROW_NUMBER() OVER (ORDER BY missing_author.author_id) AS id,
    NOW(),
    NOW(),
    NULL,
    0,
    0,
    missing_author.author_id,
    missing_author.dynasty_id,
    0,
    0
FROM (
    SELECT DISTINCT p.author_id, MIN(p.dynasty_id) OVER (PARTITION BY p.author_id) AS dynasty_id
    FROM "poem" p
    LEFT JOIN "poet" existing_poet ON existing_poet.author_id = p.author_id
        AND existing_poet.deleted_at IS NULL
    WHERE p.author_id > 0
      AND p.deleted_at IS NULL
      AND existing_poet.id IS NULL
) missing_author`,
	}
}

// setTableComments 设置数据库表注释
func setTableComments(db *gorm.DB) error {
	tableComments := map[string]string{
		"work_collection":      "作品集表 - 存储作品集信息",
		"work_collection_item": "作品集条目表 - 存储作品集与不同类型作品的关联关系",
		"user":                 "用户表 - 存储系统用户信息",
		"dynasty":              "朝代表 - 存储历史朝代信息",
		"poem_type":            "诗词类型表 - 存储诗词体裁分类信息",
		"poet":                 "诗人表 - 存储诗人信息",
		"author":               "作者表 - 存储作者信息",
		"poem":                 "诗词表 - 存储古诗词内容",
		"poem_annotation":      "诗词标注表 - 存储诗词正文选区标注",
		"comment":              "评论表 - 存储诗词评论",
		"comment_vote":         "评论投票表 - 存储评论点赞/点踩记录",
		"quiz":                 "测验题表 - 存储诗词测验题目",
		"quiz_record":          "测验记录表 - 存储用户答题记录",
		"forum_post":           "论坛帖子表 - 存储社区帖子",
		"forum_reply":          "论坛回复表 - 存储帖子回复",
		"correction_request":   "纠错申请表 - 存储诗词纠错申请",
		"correction_vote":      "纠错投票表 - 存储纠错投票记录",
		"announcement":         "公告表 - 存储系统公告",
		"feedback":             "反馈表 - 存储用户反馈",
		"operation_log":        "操作日志表 - 存储系统操作日志",
		"role":                 "角色表 - 存储RBAC角色",
		"role_permission":      "角色权限关联表 - 存储角色与权限编码的关联关系",
		"permission":           "权限信息表 - 存储每个权限点信息",
		"user_role":            "用户角色关联表 - 存储用户与角色的关联关系",
	}

	for tableName, comment := range tableComments {
		sql := fmt.Sprintf("COMMENT ON TABLE \"%s\" IS '%s'", tableName, comment)
		if err := db.Exec(sql).Error; err != nil {
			// 忽略错误，表可能不存在或已设置注释
			continue
		}
	}

	return nil
}
