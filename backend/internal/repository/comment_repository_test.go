package repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type commentSQLCaptureLogger struct {
	logger.Interface
	sql []string
}

func (l *commentSQLCaptureLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	l.sql = append(l.sql, sql)
}

func TestCommentRepositoryListAllDoesNotRequirePoemFilter(t *testing.T) {
	capture := &commentSQLCaptureLogger{Interface: logger.Discard}
	db, err := gorm.Open(postgres.Open("host=localhost user=test dbname=test sslmode=disable"), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		Logger:               capture,
	})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}
	repo := NewCommentRepository(db)

	if _, _, err := repo.ListAll(2, 20); err != nil {
		t.Fatalf("ListAll error = %v, want nil", err)
	}

	got := strings.Join(capture.sql, "\n")
	if strings.Contains(got, "poem_id =") {
		t.Fatalf("ListAll SQL = %q, should not filter by poem_id", got)
	}
	if !strings.Contains(got, `ORDER BY created_at DESC`) {
		t.Fatalf("ListAll SQL = %q, want newest comments first", got)
	}
}
