package main

import (
	"errors"
	"strings"
	"testing"
)

type sampleEntity struct {
	ID int
}

func TestSplitByBatchSizeSupportsAnyEntity(t *testing.T) {
	entities := []sampleEntity{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
	}

	batches := splitByBatchSize(entities, 2)

	if len(batches) != 3 {
		t.Fatalf("batch count = %d, want 3", len(batches))
	}
	if len(batches[0]) != 2 || len(batches[1]) != 2 || len(batches[2]) != 1 {
		t.Fatalf("batch sizes = [%d %d %d], want [2 2 1]", len(batches[0]), len(batches[1]), len(batches[2]))
	}
}

func TestSplitByBatchSizeUsesAllItemsForInvalidBatchSize(t *testing.T) {
	entities := []sampleEntity{{ID: 1}, {ID: 2}}

	batches := splitByBatchSize(entities, 0)

	if len(batches) != 1 {
		t.Fatalf("batch count = %d, want 1", len(batches))
	}
	if len(batches[0]) != len(entities) {
		t.Fatalf("batch size = %d, want %d", len(batches[0]), len(entities))
	}
}

func TestBatchInsertInTransactionCommitsWhenAllBatchesSucceed(t *testing.T) {
	runner := &fakeBatchInsertRunner{}
	entities := []sampleEntity{{ID: 1}, {ID: 2}, {ID: 3}}

	err := batchInsertInTransaction(entities, 2, runner.insert)

	if err != nil {
		t.Fatalf("batchInsertInTransaction error = %v, want nil", err)
	}
	if runner.calls != 2 {
		t.Fatalf("insert calls = %d, want 2", runner.calls)
	}
	if len(runner.batchSizes) != 2 || runner.batchSizes[0] != 2 || runner.batchSizes[1] != 1 {
		t.Fatalf("batch sizes = %v, want [2 1]", runner.batchSizes)
	}
}

func TestBatchInsertInTransactionReturnsErrorWhenABatchFails(t *testing.T) {
	runner := &fakeBatchInsertRunner{failOnCall: 2}
	entities := []sampleEntity{{ID: 1}, {ID: 2}, {ID: 3}}

	err := batchInsertInTransaction(entities, 2, runner.insert)

	if err == nil {
		t.Fatal("batchInsertInTransaction error = nil, want error")
	}
	if runner.calls != 2 {
		t.Fatalf("insert calls = %d, want 2", runner.calls)
	}
}

func TestDeduplicateAuthorsByNameRecordsDuplicateAuthors(t *testing.T) {
	authors := []author{
		{ID: "1", Name: "李白", Desc: "first"},
		{ID: "2", Name: "杜甫", Desc: "second"},
		{ID: "3", Name: "李白", Desc: "duplicate"},
		{ID: "4", Name: " 李白 ", Desc: "duplicate with spaces"},
	}

	uniqueAuthors, duplicateAuthors := deduplicateAuthorsByName(authors)

	if len(uniqueAuthors) != 2 {
		t.Fatalf("unique author count = %d, want 2", len(uniqueAuthors))
	}
	if uniqueAuthors[0].ID != "1" || uniqueAuthors[1].ID != "2" {
		t.Fatalf("unique author IDs = [%s %s], want [1 2]", uniqueAuthors[0].ID, uniqueAuthors[1].ID)
	}
	if len(duplicateAuthors) != 2 {
		t.Fatalf("duplicate author count = %d, want 2", len(duplicateAuthors))
	}
	if duplicateAuthors[0].ID != "3" || duplicateAuthors[1].ID != "4" {
		t.Fatalf("duplicate author IDs = [%s %s], want [3 4]", duplicateAuthors[0].ID, duplicateAuthors[1].ID)
	}
}

func TestBuildTruncateSQLUsesModelTableNameAndCascade(t *testing.T) {
	sql, err := buildTruncateCascadeSQL(&sampleTruncateModel{})

	if err != nil {
		t.Fatalf("buildTruncateCascadeSQL error = %v, want nil", err)
	}
	if !strings.Contains(sql, `TRUNCATE TABLE`) {
		t.Fatalf("sql = %q, want TRUNCATE TABLE", sql)
	}
	if !strings.Contains(sql, `sample_truncate_model`) {
		t.Fatalf("sql = %q, want sample_truncate_model table", sql)
	}
	if !strings.Contains(sql, `RESTART IDENTITY CASCADE`) {
		t.Fatalf("sql = %q, want RESTART IDENTITY CASCADE", sql)
	}
}

type sampleTruncateModel struct {
	ID uint64
}

func (sampleTruncateModel) TableName() string {
	return "sample_truncate_model"
}

type fakeBatchInsertRunner struct {
	calls      int
	failOnCall int
	batchSizes []int
}

func (r *fakeBatchInsertRunner) insert(batch []sampleEntity) error {
	r.calls++
	r.batchSizes = append(r.batchSizes, len(batch))
	if r.failOnCall == r.calls {
		return errors.New("insert failed")
	}
	return nil
}
