package config

import (
	"strings"
	"testing"
)

func TestBuildPoemContentJSONBMigrationSQLWrapsTextContent(t *testing.T) {
	sqlStatements := buildPoemContentJSONBMigrationSQL("poem")
	sql := strings.Join(sqlStatements, "\n")

	for _, want := range []string{
		`WHERE table_name = 'poem'`,
		`AND data_type <> 'jsonb'`,
		`ALTER TABLE "poem"`,
		`USING to_jsonb(ARRAY[content])`,
		`ALTER COLUMN content SET DEFAULT '[]'::jsonb`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("migration sql %q does not contain %q", sql, want)
		}
	}
}
