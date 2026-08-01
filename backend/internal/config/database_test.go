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

func TestBuildLegacyPoetSchemaMigrationSQLDropsOldNotNullColumns(t *testing.T) {
	sqlStatements := buildLegacyPoetSchemaMigrationSQL("poet")
	sql := strings.Join(sqlStatements, "\n")

	for _, want := range []string{
		`WHERE table_name = 'poet'`,
		`AND column_name = 'name'`,
		`ALTER TABLE "poet" ALTER COLUMN "name" DROP NOT NULL`,
		`ALTER TABLE "poet" ALTER COLUMN "biography" DROP NOT NULL`,
		`ALTER TABLE "poet" ALTER COLUMN "avatar" DROP NOT NULL`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("legacy poet migration sql %q does not contain %q", sql, want)
		}
	}
}

func TestBuildMissingPoetMigrationSQLCreatesPoetsForPoemAuthors(t *testing.T) {
	sqlStatements := buildMissingPoetMigrationSQL()
	sql := strings.Join(sqlStatements, "\n")

	for _, want := range []string{
		`INSERT INTO "poet"`,
		`? + ROW_NUMBER()`,
		`SELECT DISTINCT p.author_id`,
		`FROM "poem" p`,
		`LEFT JOIN "poet" existing_poet ON existing_poet.author_id = p.author_id`,
		`WHERE p.author_id > 0`,
		`AND p.deleted_at IS NULL`,
		`AND existing_poet.deleted_at IS NULL`,
		`AND existing_poet.id IS NULL`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("missing poet migration sql %q does not contain %q", sql, want)
		}
	}
	if strings.Contains(sql, `:id`) {
		t.Fatalf("missing poet migration sql %q should not use postgres-invalid named placeholder :id", sql)
	}
}

func TestBuildPoemPingzeBackfillSQLFillsMissingJSONArray(t *testing.T) {
	sqlStatements := buildPoemPingzeBackfillSQL("poem")
	sql := strings.Join(sqlStatements, "\n")

	for _, want := range []string{
		`UPDATE "poem"`,
		`SET "pingze" = COALESCE`,
		`jsonb_agg(to_jsonb(''::text) ORDER BY content_lines.ord)`,
		`WHERE "pingze" IS NULL`,
		`OR jsonb_typeof("pingze") <> 'array'`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("pingze backfill sql %q does not contain %q", sql, want)
		}
	}
}

func TestBuildPoemGenreCategoryBackfillSQLUsesPoemTypes(t *testing.T) {
	sqlStatements := buildPoemGenreCategoryBackfillSQL("poem")
	sql := strings.Join(sqlStatements, "\n")

	for _, want := range []string{
		`UPDATE "poem" AS p`,
		`SET "genre_category" = pt."category"`,
		`FROM "poem_type" AS pt`,
		`AND p."genre" = pt."name"`,
		`FROM "ci_tune" AS ct`,
		`JOIN "poem_type" AS pt ON pt."id" = ct."poem_type_id"`,
		`AND p."ci_tune_id" = ct."id"`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("genre category backfill sql %q does not contain %q", sql, want)
		}
	}
}

func TestMarshalPoemPingzeJSONReturnsJSONArray(t *testing.T) {
	got, err := marshalPoemPingzeJSON([]string{"平仄?", "仄平"})

	if err != nil {
		t.Fatalf("marshalPoemPingzeJSON error = %v, want nil", err)
	}
	if got != `["平仄?","仄平"]` {
		t.Fatalf("marshalPoemPingzeJSON = %q, want JSON array", got)
	}
}
