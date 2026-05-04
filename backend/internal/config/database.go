package config

import (
	"fmt"
	"shihai/internal/models"

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

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Dynasty{},
		&models.Poet{},
		&models.Poem{},
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
		// RBAC models
		&models.Role{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.Permission{},
	)
	if err != nil {
		return nil, err
	}

	// Set table comments
	if err := setTableComments(db); err != nil {
		return nil, err
	}

	return db, nil
}

// setTableComments 设置数据库表注释
func setTableComments(db *gorm.DB) error {
	tableComments := map[string]string{
		"user":               "用户表 - 存储系统用户信息",
		"dynasty":            "朝代表 - 存储历史朝代信息",
		"poet":               "诗人表 - 存储诗人信息",
		"poem":               "诗词表 - 存储古诗词内容",
		"comment":            "评论表 - 存储诗词评论",
		"comment_vote":       "评论投票表 - 存储评论点赞/点踩记录",
		"quiz":               "测验题表 - 存储诗词测验题目",
		"quiz_record":        "测验记录表 - 存储用户答题记录",
		"forum_post":         "论坛帖子表 - 存储社区帖子",
		"forum_reply":        "论坛回复表 - 存储帖子回复",
		"correction_request": "纠错申请表 - 存储诗词纠错申请",
		"correction_vote":    "纠错投票表 - 存储纠错投票记录",
		"announcement":       "公告表 - 存储系统公告",
		"feedback":           "反馈表 - 存储用户反馈",
		"operation_log":      "操作日志表 - 存储系统操作日志",
		"role":               "角色表 - 存储RBAC角色",
		"role_permission":    "角色权限关联表 - 存储角色与权限编码的关联关系",
		"user_role":          "用户角色关联表 - 存储用户与角色的关联关系",
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
