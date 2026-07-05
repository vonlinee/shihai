package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"shihai/internal/config"
	"shihai/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var root = "D:\\Develop\\Code\\Github\\chinese-poetry"

func openDBWithGorm() (*gorm.DB, error) {
	cfg := config.Load("config.json")
	return config.InitDB(&cfg.Database)
}

func main() {
	db, err := openDBWithGorm()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Printf("%p", db)

	err = truncateModelCascade(*db, models.Author{})
	if err != nil {
		return
	}
	synAuthor(*db, root+"/全唐诗/authors.tang.json", 500)
	synAuthor(*db, root+"/全唐诗/authors.song.json", 500)
}

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
	return db.Transaction(func(tx *gorm.DB) error {
		return batchInsertInTransaction(entities, batchSize, func(batch []T) error {
			return tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(batch, len(batch)).Error
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
