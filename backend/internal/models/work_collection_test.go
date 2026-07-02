package models

import (
	"sync"
	"testing"

	"gorm.io/gorm/schema"
)

func TestWorkCollectionTableNames(t *testing.T) {
	if got := (WorkCollection{}).TableName(); got != "work_collection" {
		t.Fatalf("WorkCollection table name = %q, want %q", got, "work_collection")
	}

	if got := (WorkCollectionItem{}).TableName(); got != "work_collection_item" {
		t.Fatalf("WorkCollectionItem table name = %q, want %q", got, "work_collection_item")
	}
}

func TestWorkCollectionItemUniqueWorkReferenceIndex(t *testing.T) {
	var cache sync.Map
	parsed, err := schema.Parse(&WorkCollectionItem{}, &cache, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("parse WorkCollectionItem schema: %v", err)
	}

	indexes := parsed.ParseIndexes()
	index, ok := indexes["idx_collection_work"]
	if !ok {
		t.Fatalf("missing idx_collection_work index; indexes: %#v", indexes)
	}

	if index.Class != "UNIQUE" {
		t.Fatalf("idx_collection_work class = %q, want UNIQUE", index.Class)
	}

	gotColumns := make([]string, 0, len(index.Fields))
	for _, field := range index.Fields {
		gotColumns = append(gotColumns, field.DBName)
	}

	wantColumns := []string{"collection_id", "work_type", "work_id"}
	if len(gotColumns) != len(wantColumns) {
		t.Fatalf("idx_collection_work columns = %#v, want %#v", gotColumns, wantColumns)
	}
	for i := range wantColumns {
		if gotColumns[i] != wantColumns[i] {
			t.Fatalf("idx_collection_work columns = %#v, want %#v", gotColumns, wantColumns)
		}
	}
}
