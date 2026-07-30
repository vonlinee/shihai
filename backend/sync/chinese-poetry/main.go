package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgconn"

	"shihai/internal/config"
	"shihai/internal/models"
	"shihai/internal/poetry"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var root = "D:\\Develop\\Code\\Github\\chinese-poetry"
var tangDynastyName = "唐"

func openDBWithGorm() (*gorm.DB, error) {
	cfg := config.Load("config.json")
	return config.InitDB(&cfg.Database)
}

type poem struct {
	Author     string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Title      string   `json:"title"`
	ID         string   `json:"id"`
}

type ci struct {
	Author     string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Rhythmic   string   `json:"rhythmic"`
}

func main() {
	db, err := openDBWithGorm()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = truncateModelCascade(*db, models.Author{})
	if err != nil {
		return
	}
	synAuthor(*db, root+"/全唐诗/authors.tang.json", 500)
	synAuthor(*db, root+"/全唐诗/authors.song.json", 500)

	synSongCiAuthor(*db, root+"/宋词/author.song.json", 500)

	if err := ensurePoemAuthorForeignKey(*db); err != nil {
		log.Fatalf("修复诗歌作者外键失败: %s", err)
	}

	err = truncateModelCascade(*db, models.Poem{})
	if err != nil {
		return
	}
	syncAllPoem(*db, root+"/全唐诗")
	syncAllPoemCi(*db, root+"/宋词")
	// TODO 宋词三百首.json

	if err := config.MigrateMissingPoetsFromPoems(db); err != nil {
		log.Fatalf("补齐诗人数据失败: %s", err)
	}
}

// syncAllPoemCi 同步宋词
// chinese-poetry\宋词 目录下 ci.song.xxx.json 文件
func syncAllPoemCi(db gorm.DB, path string) {
	err := walkDirectChildFiles(path, func(path string, entry os.DirEntry) error {
		if strings.HasPrefix(entry.Name(), "ci") {
			arr := strings.Split(entry.Name(), ".")
			if arr[1] == "song" {
				fmt.Printf(">>> read poem ci from %s\n", path)
				syncCi(db, path, 500)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("同步诗歌失败: %s", err)
	}
}

// 同步宋词信息
func syncCi(db gorm.DB, path string, batchSize int) {
	text, _ := readFileToString(path)
	var cis []ci
	if err := json.Unmarshal([]byte(text), &cis); err != nil {
		log.Fatalf("解析失败: %s", err)
	}

	authorIDs, err := loadAuthorIDsByName(db)
	if err != nil {
		log.Fatalf("查询作者失败: %s", err)
		return
	}
	dynastyID, err := findDynastyIDByName(db, tangDynastyName)
	if err != nil {
		log.Fatalf("查询朝代失败: %s", err)
		return
	}

	poemModels, skippedPoems := buildCiModels(cis, authorIDs, dynastyID)
	for _, p := range skippedPoems {
		log.Printf("skip poem without known author from %s: title=%q author=%q", path, p.Rhythmic, p.Author)
	}
	log.Printf("prepared poems from %s: total=%d insert=%d skipped=%d authorMappings=%d dynastyID=%d", path, len(cis), len(poemModels), len(skippedPoems), len(authorIDs), dynastyID)

	if err := validatePoemAuthorReferences(db, poemModels); err != nil {
		log.Fatalf("校验词作者外键失败: %s", err)
		return
	}
	if err := batchInsert(db, poemModels, batchSize); err != nil {
		printPgError("批量插入词失败:", err)
	}
}

func syncAllPoem(db gorm.DB, path string) {
	err := walkDirectChildFiles(path, func(path string, entry os.DirEntry) error {
		if strings.HasPrefix(entry.Name(), "poet") {
			arr := strings.Split(entry.Name(), ".")
			if arr[1] == "tang" {
				fmt.Printf(">>> read poem from %s\n", path)
				syncPoem(db, path, 500)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("同步诗歌失败: %s", err)
	}
}

// 同步诗歌信息
func syncPoem(db gorm.DB, path string, batchSize int) {
	text, _ := readFileToString(path)
	var poems []poem
	if err := json.Unmarshal([]byte(text), &poems); err != nil {
		log.Fatalf("解析失败: %s", err)
	}

	authorIDs, err := loadAuthorIDsByName(db)
	if err != nil {
		log.Fatalf("查询作者失败: %s", err)
		return
	}
	dynastyID, err := findDynastyIDByName(db, tangDynastyName)
	if err != nil {
		log.Fatalf("查询朝代失败: %s", err)
		return
	}

	poemModels, skippedPoems := buildPoemModels(poems, authorIDs, dynastyID)
	logSkippedPoems(path, skippedPoems)
	log.Printf("prepared poems from %s: total=%d insert=%d skipped=%d authorMappings=%d dynastyID=%d", path, len(poems), len(poemModels), len(skippedPoems), len(authorIDs), dynastyID)

	if err := validatePoemAuthorReferences(db, poemModels); err != nil {
		log.Fatalf("校验诗歌作者外键失败: %s", err)
		return
	}
	if err := batchInsert(db, poemModels, batchSize); err != nil {
		printPgError("批量插入诗歌失败:", err)
	}
}

// printPgError 打印 PostgreSQL 错误的关键诊断字段。
//
// err 待打印的错误；当错误链中包含 pgconn.PgError 时，输出数据库错误码、表名、列名和约束名等信息。
func printPgError(msg string, err error) {
	log.Print(msg, formatPgError(err))
}

// formatPgError 格式化 PostgreSQL 错误的关键诊断字段。
//
// err 待格式化的错误；当错误链中包含 pgconn.PgError 时，返回 PostgreSQL 诊断信息，否则返回普通错误文本。
func formatPgError(err error) string {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fmt.Sprintf("error: %v", err)
	}

	fields := []string{
		"PostgreSQL error",
		fmt.Sprintf("severity=%s", pgErr.Severity),
		fmt.Sprintf("code=%s", pgErr.Code),
		fmt.Sprintf("message=%s", pgErr.Message),
		fmt.Sprintf("detail=%s", pgErr.Detail),
		fmt.Sprintf("hint=%s", pgErr.Hint),
		fmt.Sprintf("schema=%s", pgErr.SchemaName),
		fmt.Sprintf("table=%s", pgErr.TableName),
		fmt.Sprintf("column=%s", pgErr.ColumnName),
		fmt.Sprintf("constraint=%s", pgErr.ConstraintName),
		fmt.Sprintf("where=%s", pgErr.Where),
	}
	return strings.Join(fields, "\n")
}

func buildCiModels(poems []ci, authorIDs map[string]uint64, dynastyID uint64) ([]models.Poem, []ci) {
	poemModels := make([]models.Poem, 0, len(poems))
	skippedPoems := make([]ci, 0)

	for _, p := range poems {
		authorName := strings.TrimSpace(p.Author)
		authorID, ok := authorIDs[authorName]
		if !ok || authorID == 0 {
			skippedPoems = append(skippedPoems, p)
			continue
		}

		content := buildPoemContent(p.Paragraphs)
		poemModels = append(poemModels, models.Poem{
			Title:     strings.TrimSpace(p.Rhythmic),
			Content:   content,
			Pingze:    poetry.RecognizePingzeLines(content),
			AuthorID:  authorID,
			DynastyID: dynastyID,
		})
	}
	return poemModels, skippedPoems
}

// buildPoemModels 将解析出的诗歌数据转换为可入库的诗歌模型。
//
// poems 待转换的原始诗歌数据。
// authorIDs 作者名称到作者 ID 的映射，调用方需要保证映射来自已入库作者。
// dynastyID 诗歌所属朝代 ID。
// 返回可入库诗歌模型，以及因作者缺失而跳过的原始诗歌。
func buildPoemModels(poems []poem, authorIDs map[string]uint64, dynastyID uint64) ([]models.Poem, []poem) {
	poemModels := make([]models.Poem, 0, len(poems))
	skippedPoems := make([]poem, 0)

	for _, p := range poems {
		authorName := strings.TrimSpace(p.Author)
		authorID, ok := authorIDs[authorName]
		if !ok || authorID == 0 {
			skippedPoems = append(skippedPoems, p)
			continue
		}

		content := buildPoemContent(p.Paragraphs)
		poemModels = append(poemModels, models.Poem{
			Title:     strings.TrimSpace(p.Title),
			Content:   content,
			Pingze:    poetry.RecognizePingzeLines(content),
			AuthorID:  authorID,
			DynastyID: dynastyID,
		})
	}
	return poemModels, skippedPoems
}

// buildPoemContent 将原始段落合并为单个正文元素。
//
// paragraphs 原始诗词段落数组，通常来自 chinese-poetry 数据源的 paragraphs 字段。
// 返回包含一个合并正文元素的 JSON 数组；空白段落会被跳过，避免存入无意义内容。
func buildPoemContent(paragraphs []string) []string {
	items := make([]string, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		trimmedParagraph := strings.TrimSpace(paragraph)
		if trimmedParagraph == "" {
			continue
		}
		items = append(items, trimmedParagraph)
	}
	if len(items) == 0 {
		return []string{}
	}
	return []string{strings.Join(items, "")}
}

// validatePoemAuthorReferences 校验待插入诗歌的作者外键是否存在。
//
// db 数据库连接。
// poemModels 待插入的诗歌模型。
// 当存在未入库的作者 ID 时返回错误，并在错误中包含样例诗歌信息。
func validatePoemAuthorReferences(db gorm.DB, poemModels []models.Poem) error {
	authorIDs := make([]uint64, 0, len(poemModels))
	seen := make(map[uint64]struct{}, len(poemModels))
	for _, p := range poemModels {
		if _, ok := seen[p.AuthorID]; ok {
			continue
		}
		seen[p.AuthorID] = struct{}{}
		authorIDs = append(authorIDs, p.AuthorID)
	}
	if len(authorIDs) == 0 {
		return nil
	}

	var existingAuthorIDs []uint64
	if err := db.Model(&models.Author{}).Where("id IN ?", authorIDs).Pluck("id", &existingAuthorIDs).Error; err != nil {
		return err
	}

	existingAuthorIDSet := make(map[uint64]struct{}, len(existingAuthorIDs))
	for _, id := range existingAuthorIDs {
		existingAuthorIDSet[id] = struct{}{}
	}

	missingPoems := findPoemsWithMissingAuthorIDs(poemModels, existingAuthorIDSet)
	if len(missingPoems) == 0 {
		return nil
	}
	sample := missingPoems[0]
	return fmt.Errorf("存在 %d 首诗歌引用了不存在的作者 ID，示例: title=%q authorID=%d", len(missingPoems), sample.Title, sample.AuthorID)
}

// findPoemsWithMissingAuthorIDs 找出引用缺失作者 ID 的诗歌。
//
// poemModels 待检查的诗歌模型。
// existingAuthorIDs 已存在的作者 ID 集合。
// 返回作者 ID 为 0 或不在 existingAuthorIDs 中的诗歌。
func findPoemsWithMissingAuthorIDs(poemModels []models.Poem, existingAuthorIDs map[uint64]struct{}) []models.Poem {
	missingPoems := make([]models.Poem, 0)
	for _, p := range poemModels {
		if p.AuthorID == 0 {
			missingPoems = append(missingPoems, p)
			continue
		}
		if _, ok := existingAuthorIDs[p.AuthorID]; !ok {
			missingPoems = append(missingPoems, p)
		}
	}
	return missingPoems
}

// ensurePoemAuthorForeignKey ensures poem.author_id references author.id before syncing poems.
//
// db 数据库连接。该方法会读取当前 fk_poem_author 约束定义；当约束不存在或仍指向旧的 poet 表时，会重建为
// poem.author_id -> author.id。查询或执行修复 SQL 失败时返回错误。
func ensurePoemAuthorForeignKey(db gorm.DB) error {
	const constraintName = "fk_poem_author"

	definition, err := findConstraintDefinition(db, "poem", constraintName)
	if err != nil {
		return err
	}
	if !poemAuthorForeignKeyNeedsRepair(definition) {
		return nil
	}

	log.Printf("repair %s definition: %s", constraintName, definition)
	return db.Transaction(func(tx *gorm.DB) error {
		for _, sql := range buildPoemAuthorForeignKeyRepairSQL() {
			if err := tx.Exec(sql).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// findConstraintDefinition returns a PostgreSQL constraint definition by table and constraint name.
//
// db 数据库连接。
// tableName 约束所在表名。
// constraintName 约束名称。
// 返回 PostgreSQL pg_get_constraintdef 生成的约束定义；未找到约束时返回空字符串。
func findConstraintDefinition(db gorm.DB, tableName string, constraintName string) (string, error) {
	var definition string
	err := db.Raw(`
SELECT COALESCE(pg_get_constraintdef(c.oid), '')
FROM pg_constraint c
JOIN pg_class t ON t.oid = c.conrelid
JOIN pg_namespace n ON n.oid = t.relnamespace
WHERE n.nspname = current_schema()
  AND t.relname = ?
  AND c.conname = ?
`, tableName, constraintName).Scan(&definition).Error
	return definition, err
}

// poemAuthorForeignKeyNeedsRepair reports whether fk_poem_author should be rebuilt.
//
// definition 当前 PostgreSQL 外键定义，通常来自 pg_get_constraintdef；空定义代表约束不存在。
// 返回 true 时调用方应将 poem.author_id 外键重建为引用 author.id。
func poemAuthorForeignKeyNeedsRepair(definition string) bool {
	normalizedDefinition := strings.ToLower(strings.Join(strings.Fields(definition), " "))
	if normalizedDefinition == "" {
		return true
	}
	return !strings.Contains(normalizedDefinition, "foreign key (author_id) references author(id)")
}

// buildPoemAuthorForeignKeyRepairSQL builds SQL statements that rebuild poem.author_id to reference author.id.
//
// 返回的 SQL 语句会先删除旧 fk_poem_author 约束，再创建指向 author(id) 的级联更新、限制删除外键。
func buildPoemAuthorForeignKeyRepairSQL() []string {
	return []string{
		`ALTER TABLE "poem" DROP CONSTRAINT IF EXISTS "fk_poem_author"`,
		`ALTER TABLE "poem" ADD CONSTRAINT "fk_poem_author" FOREIGN KEY ("author_id") REFERENCES "author"("id") ON UPDATE CASCADE ON DELETE RESTRICT`,
	}
}

// loadAuthorIDsByName 查询作者名称到作者 ID 的映射。
//
// db 数据库连接。
// 返回以去除首尾空白后的作者名称为键、作者 ID 为值的映射；查询失败时返回错误。
func loadAuthorIDsByName(db gorm.DB) (map[string]uint64, error) {
	var authors []models.Author
	if err := db.Find(&authors).Error; err != nil {
		return nil, err
	}

	authorIDs := make(map[string]uint64, len(authors))
	for _, a := range authors {
		name := strings.TrimSpace(a.Name)
		if name == "" || a.ID == 0 {
			continue
		}
		authorIDs[name] = a.ID
	}
	return authorIDs, nil
}

// findDynastyIDByName 根据朝代名称查询朝代 ID。
//
// db 数据库连接。
// name 朝代名称，调用方传入前无需自行去除首尾空白。
// 返回匹配朝代的 ID；未找到或查询失败时返回错误。
func findDynastyIDByName(db gorm.DB, name string) (uint64, error) {
	var dynasty models.Dynasty
	if err := db.Where("name = ?", strings.TrimSpace(name)).First(&dynasty).Error; err != nil {
		return 0, err
	}
	return dynasty.ID, nil
}

// 读取文件内容
func readFileToString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type author struct {
	Desc string `json:"desc"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

type authorSongCi struct {
	Desc             string `json:"description"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
}

func synSongCiAuthor(db gorm.DB, path string, batchSize int) {
	authorsContent, err := readFileToString(path)
	if err != nil {
		log.Fatalf("读取文件失败: %s", err)
	}

	var authors []authorSongCi
	if err := json.Unmarshal([]byte(authorsContent), &authors); err != nil {
		log.Fatalf("解析失败: %s", err)
	}
	log.Printf("read %d authors from %s\n", len(authors), path)

	uniqueAuthors, duplicateAuthors := deduplicateSongCiAuthorsByName(authors)
	for _, a := range duplicateAuthors {
		log.Printf("skip duplicate author from %s: id=%s name=%q", path, a.Desc, a.Name)
	}

	authorModels := make([]models.Author, 0, len(uniqueAuthors))
	for _, a := range uniqueAuthors {
		authorModels = append(authorModels, models.Author{
			Name:      strings.TrimSpace(a.Name),
			Biography: strings.TrimSpace(a.Desc),
		})
	}

	if err := batchInsert(db, authorModels, batchSize); err != nil {
		log.Fatalf("批量插入作者失败: %s", err)
	}
}

func deduplicateSongCiAuthorsByName(authors []authorSongCi) ([]authorSongCi, []authorSongCi) {
	seen := make(map[string]struct{}, len(authors))
	uniqueAuthors := make([]authorSongCi, 0, len(authors))
	duplicateAuthors := make([]authorSongCi, 0)

	for _, a := range authors {
		name := strings.TrimSpace(a.Name)
		if name == "" {
			duplicateAuthors = append(duplicateAuthors, a)
			continue
		}
		if _, ok := seen[name]; ok {
			duplicateAuthors = append(duplicateAuthors, a)
			continue
		}

		seen[name] = struct{}{}
		a.Name = name
		uniqueAuthors = append(uniqueAuthors, a)
	}
	return uniqueAuthors, duplicateAuthors
}

func synAuthor(db gorm.DB, path string, batchSize int) {
	authorsContent, err := readFileToString(path)
	if err != nil {
		log.Fatalf("读取文件失败: %s", err)
	}

	var authors []author
	if err := json.Unmarshal([]byte(authorsContent), &authors); err != nil {
		log.Fatalf("解析失败: %s", err)
	}
	log.Printf("read %d authors from %s\n", len(authors), path)

	uniqueAuthors, duplicateAuthors := deduplicateAuthorsByName(authors)
	logDuplicateAuthors(path, duplicateAuthors)

	authorModels := make([]models.Author, 0, len(uniqueAuthors))
	for _, a := range uniqueAuthors {
		authorModels = append(authorModels, models.Author{
			Name:      strings.TrimSpace(a.Name),
			Biography: strings.TrimSpace(a.Desc),
		})
	}

	if err := batchInsert(db, authorModels, batchSize); err != nil {
		log.Fatalf("批量插入作者失败: %s", err)
	}
}

func batchInsert[T any](db gorm.DB, entities []T, batchSize int) error {
	if len(entities) == 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		return batchInsertInTransaction(entities, batchSize, func(batch []T) error {
			return tx.Omit(clause.Associations).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(batch, len(batch)).Error
		})
	})
}

func truncateModelCascade(db gorm.DB, model any) error {
	sql, err := buildTruncateCascadeSQL(model)
	if err != nil {
		return err
	}
	return db.Exec(sql).Error
}

func buildTruncateCascadeSQL(model any) (string, error) {
	modelSchema, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return "", err
	}
	return "TRUNCATE TABLE " + quoteIdentifier(modelSchema.Table) + " RESTART IDENTITY CASCADE", nil
}

func quoteIdentifier(identifier string) string {
	parts := strings.Split(identifier, ".")
	for i, part := range parts {
		parts[i] = `"` + strings.ReplaceAll(part, `"`, `""`) + `"`
	}
	return strings.Join(parts, ".")
}

func batchInsertInTransaction[T any](entities []T, batchSize int, insertBatch func([]T) error) error {
	for _, batch := range splitByBatchSize(entities, batchSize) {
		if err := insertBatch(batch); err != nil {
			return err
		}
	}
	return nil
}

func splitByBatchSize[T any](items []T, batchSize int) [][]T {
	if len(items) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = len(items)
	}

	batches := make([][]T, 0, (len(items)+batchSize-1)/batchSize)
	for start := 0; start < len(items); start += batchSize {
		end := start + batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[start:end])
	}
	return batches
}

func deduplicateAuthorsByName(authors []author) ([]author, []author) {
	seen := make(map[string]struct{}, len(authors))
	uniqueAuthors := make([]author, 0, len(authors))
	duplicateAuthors := make([]author, 0)

	for _, a := range authors {
		name := strings.TrimSpace(a.Name)
		if name == "" {
			duplicateAuthors = append(duplicateAuthors, a)
			continue
		}
		if _, ok := seen[name]; ok {
			duplicateAuthors = append(duplicateAuthors, a)
			continue
		}

		seen[name] = struct{}{}
		a.Name = name
		uniqueAuthors = append(uniqueAuthors, a)
	}
	return uniqueAuthors, duplicateAuthors
}

func logDuplicateAuthors(path string, duplicateAuthors []author) {
	for _, a := range duplicateAuthors {
		log.Printf("skip duplicate author from %s: id=%s name=%q", path, a.ID, a.Name)
	}
}

// logSkippedPoems 记录因缺少作者映射而跳过的诗歌。
//
// path 当前同步的诗歌文件路径。
// skippedPoems 被跳过的原始诗歌列表。
func logSkippedPoems(path string, skippedPoems []poem) {
	for _, p := range skippedPoems {
		log.Printf("skip poem without known author from %s: id=%s title=%q author=%q", path, p.ID, p.Title, p.Author)
	}
}

// walkDirectChildFiles 遍历 rootDir 目录下的直接子文件。
//
// rootDir 待遍历的目录路径。
// handleFile 接收文件完整路径和文件目录项；当回调返回错误时，遍历立即停止并返回该错误。
func walkDirectChildFiles(rootDir string, handleFile func(path string, entry os.DirEntry) error) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if err := handleFile(filepath.Join(rootDir, entry.Name()), entry); err != nil {
			return err
		}
	}
	return nil
}

// walkAllFiles 递归遍历 rootDir 目录下的所有文件。
//
// rootDir 待遍历的根目录路径。
// handleFile 接收文件完整路径和文件目录项；当回调返回错误时，遍历立即停止并返回该错误。
func walkAllFiles(rootDir string, handleFile func(path string, entry os.DirEntry) error) error {
	return filepath.WalkDir(rootDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		return handleFile(path, entry)
	})
}
