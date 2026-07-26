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

type sqlCaptureLogger struct {
	logger.Interface
	sql []string
}

func (l *sqlCaptureLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	l.sql = append(l.sql, sql)
}

func TestDistinctGenresUsesPoemTypeReferenceData(t *testing.T) {
	capture := &sqlCaptureLogger{Interface: logger.Discard}
	db, err := gorm.Open(postgres.Open("host=localhost user=test dbname=test sslmode=disable"), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		Logger:               capture,
	})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}
	repo := NewPoemRepository(db)

	if _, err := repo.DistinctGenres(); err != nil {
		t.Fatalf("DistinctGenres error = %v, want nil", err)
	}

	got := strings.Join(capture.sql, "\n")
	if !strings.Contains(got, `"poem_type"`) {
		t.Fatalf("DistinctGenres SQL = %q, want query poem_type reference table", got)
	}
	if strings.Contains(got, `"poem"`) {
		t.Fatalf("DistinctGenres SQL = %q, should not depend on poem rows", got)
	}
}

func TestListPoemTypesUsesPoemTypeReferenceData(t *testing.T) {
	capture := &sqlCaptureLogger{Interface: logger.Discard}
	db, err := gorm.Open(postgres.Open("host=localhost user=test dbname=test sslmode=disable"), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		Logger:               capture,
	})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}
	repo := NewPoemRepository(db)

	if _, err := repo.ListPoemTypes(); err != nil {
		t.Fatalf("ListPoemTypes error = %v, want nil", err)
	}

	got := strings.Join(capture.sql, "\n")
	if !strings.Contains(got, `"poem_type"`) {
		t.Fatalf("ListPoemTypes SQL = %q, want query poem_type reference table", got)
	}
	if strings.Contains(got, `"poem"`) {
		t.Fatalf("ListPoemTypes SQL = %q, should not depend on poem rows", got)
	}
}
