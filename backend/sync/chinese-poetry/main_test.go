package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"shihai/internal/models"
	"shihai/pkg/utils"
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

func TestPoemAuthorForeignKeyNeedsRepairWhenReferencedTableIsPoet(t *testing.T) {
	const definition = `FOREIGN KEY (author_id) REFERENCES poet(id)`

	if !poemAuthorForeignKeyNeedsRepair(definition) {
		t.Fatalf("poemAuthorForeignKeyNeedsRepair(%q) = false, want true", definition)
	}
}

func TestPoemAuthorForeignKeyNeedsRepairReturnsFalseForAuthorReference(t *testing.T) {
	const definition = `FOREIGN KEY (author_id) REFERENCES author(id)`

	if poemAuthorForeignKeyNeedsRepair(definition) {
		t.Fatalf("poemAuthorForeignKeyNeedsRepair(%q) = true, want false", definition)
	}
}

func TestBuildPoemAuthorForeignKeyRepairSQLReferencesAuthorTable(t *testing.T) {
	sqlStatements := buildPoemAuthorForeignKeyRepairSQL()
	sql := strings.Join(sqlStatements, "\n")

	for _, want := range []string{
		`ALTER TABLE "poem" DROP CONSTRAINT IF EXISTS "fk_poem_author"`,
		`ADD CONSTRAINT "fk_poem_author" FOREIGN KEY ("author_id") REFERENCES "author"("id")`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("repair sql %q does not contain %q", sql, want)
		}
	}
}

func TestBuildPoemModelsResolvesAuthorAndDynastyIDs(t *testing.T) {
	rawPoems := []poem{
		{Author: " 李白 ", Title: " 静夜思 ", Paragraphs: []string{"床前明月光，", "疑是地上霜。"}},
	}
	authorIDs := map[string]uint64{"李白": 101}

	poemModels, skippedPoems := buildPoemModels(rawPoems, authorIDs, 202, nil)

	if len(skippedPoems) != 0 {
		t.Fatalf("skipped poem count = %d, want 0", len(skippedPoems))
	}
	if len(poemModels) != 1 {
		t.Fatalf("poem model count = %d, want 1", len(poemModels))
	}
	if poemModels[0].AuthorID != 101 {
		t.Fatalf("author ID = %d, want 101", poemModels[0].AuthorID)
	}
	if poemModels[0].DynastyID != 202 {
		t.Fatalf("dynasty ID = %d, want 202", poemModels[0].DynastyID)
	}
	if poemModels[0].Title != "静夜思" {
		t.Fatalf("title = %q, want 静夜思", poemModels[0].Title)
	}
	assertStringSliceEqual(t, poemModels[0].Content, []string{"床前明月光，疑是地上霜。"})
	assertStringSliceEqual(t, poemModels[0].Pingze, []string{"平平平仄平，平仄仄仄平。"})
}

func TestBuildPoemModelsStoresJoinedParagraphsAsSingleContentItem(t *testing.T) {
	rawPoems := []poem{
		{Author: "李白", Title: "静夜思", Paragraphs: []string{"床前明月光，", "疑是地上霜。", "举头望明月，", "低头思故乡。"}},
	}
	authorIDs := map[string]uint64{"李白": 101}

	poemModels, skippedPoems := buildPoemModels(rawPoems, authorIDs, 202, nil)

	if len(skippedPoems) != 0 {
		t.Fatalf("skipped poem count = %d, want 0", len(skippedPoems))
	}
	if len(poemModels) != 1 {
		t.Fatalf("poem model count = %d, want 1", len(poemModels))
	}
	assertStringSliceEqual(t, poemModels[0].Content, []string{"床前明月光，疑是地上霜。举头望明月，低头思故乡。"})
	assertStringSliceEqual(t, poemModels[0].Pingze, []string{"平平平仄平，平仄仄仄平。仄平仄平仄，平平平仄平。"})
}

func TestBuildPoemModelsSkipsPoemsWithoutKnownAuthor(t *testing.T) {
	rawPoems := []poem{
		{Author: "李白", Title: "静夜思", Paragraphs: []string{"床前明月光。"}},
		{Author: "未知作者", Title: "无作者诗", Paragraphs: []string{"内容。"}},
		{Author: " ", Title: "空作者诗", Paragraphs: []string{"内容。"}},
	}
	authorIDs := map[string]uint64{"李白": 101}

	poemModels, skippedPoems := buildPoemModels(rawPoems, authorIDs, 202, nil)

	if len(poemModels) != 1 {
		t.Fatalf("poem model count = %d, want 1", len(poemModels))
	}
	if poemModels[0].AuthorID == 0 {
		t.Fatal("author ID = 0, want valid author ID")
	}
	if len(skippedPoems) != 2 {
		t.Fatalf("skipped poem count = %d, want 2", len(skippedPoems))
	}
	if skippedPoems[0].Title != "无作者诗" || skippedPoems[1].Title != "空作者诗" {
		t.Fatalf("skipped poem titles = [%q %q], want [无作者诗 空作者诗]", skippedPoems[0].Title, skippedPoems[1].Title)
	}
}

func TestIsPoemDataFileMatchesConfiguredSource(t *testing.T) {
	tangSource := poemSyncSource{filePrefix: "poet", fileDynasty: "tang"}
	songSource := poemSyncSource{filePrefix: "ci", fileDynasty: "song"}
	tests := []struct {
		name        string
		fileName    string
		source      poemSyncSource
		wantMatched bool
	}{
		{name: "唐诗文件", fileName: "poet.tang.0.json", source: tangSource, wantMatched: true},
		{name: "宋词文件", fileName: "ci.song.1000.json", source: songSource, wantMatched: true},
		{name: "宋词不能按唐诗前缀匹配", fileName: "ci.song.0.json", source: tangSource, wantMatched: false},
		{name: "朝代不匹配", fileName: "poet.song.0.json", source: tangSource, wantMatched: false},
		{name: "缺少数据分片", fileName: "ci.song.json", source: songSource, wantMatched: false},
		{name: "非 JSON 文件", fileName: "ci.song.0.txt", source: songSource, wantMatched: false},
		{name: "短文件名", fileName: "ci", source: songSource, wantMatched: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPoemDataFile(tt.fileName, tt.source); got != tt.wantMatched {
				t.Errorf("isPoemDataFile(%q) = %v, want %v", tt.fileName, got, tt.wantMatched)
			}
		})
	}
}

func TestCollectAuthorNamesTrimsDeduplicatesAndSorts(t *testing.T) {
	poems := []poem{
		{Author: " 苏轼 "},
		{Author: "李白"},
		{Author: "苏轼"},
		{Author: " "},
	}

	names := collectAuthorNames(poems)

	assertStringSliceEqual(t, names, []string{"李白", "苏轼"})
}

func TestBuildMissingAuthorsExcludesExistingNames(t *testing.T) {
	names := []string{"李白", "杜甫", "苏轼"}
	existingIDs := map[string]uint64{"杜甫": 101}

	missing := buildMissingAuthors(names, existingIDs)

	if len(missing) != 2 {
		t.Fatalf("missing author count = %d, want 2", len(missing))
	}
	if missing[0].Name != "李白" || missing[1].Name != "苏轼" {
		t.Fatalf("missing author names = [%q %q], want [李白 苏轼]", missing[0].Name, missing[1].Name)
	}
}

func TestMergeAuthorIDsUsesNormalizedValidRecords(t *testing.T) {
	authorIDs := map[string]uint64{"已有作者": 1}
	mergeAuthorIDs(authorIDs, []models.Author{
		{BaseModel: models.BaseModel{ID: 101}, Name: " 苏轼 "},
		{Name: "无 ID 作者"},
		{BaseModel: models.BaseModel{ID: 102}, Name: " "},
	})

	if authorIDs["已有作者"] != 1 || authorIDs["苏轼"] != 101 {
		t.Fatalf("author IDs = %v, want existing and normalized IDs", authorIDs)
	}
	if _, exists := authorIDs["无 ID 作者"]; exists {
		t.Fatal("author IDs contains record without ID")
	}
}

func TestResolveAuthorIDsUsesBatchQueriesAndBatchInserts(t *testing.T) {
	names := []string{"李白", "杜甫", "白居易", "苏轼", "辛弃疾"}
	listCalls := 0
	insertBatchSizes := make([]int, 0)

	authorIDs, err := resolveAuthorIDs(
		names,
		2,
		func(gotNames []string) ([]models.Author, error) {
			listCalls++
			if len(gotNames) != len(names) {
				t.Fatalf("batch query name count = %d, want %d", len(gotNames), len(names))
			}
			if listCalls == 1 {
				return []models.Author{{BaseModel: models.BaseModel{ID: 102}, Name: "杜甫"}}, nil
			}
			return []models.Author{
				{BaseModel: models.BaseModel{ID: 101}, Name: "李白"},
				{BaseModel: models.BaseModel{ID: 102}, Name: "杜甫"},
				{BaseModel: models.BaseModel{ID: 103}, Name: "白居易"},
				{BaseModel: models.BaseModel{ID: 104}, Name: "苏轼"},
				{BaseModel: models.BaseModel{ID: 105}, Name: "辛弃疾"},
			}, nil
		},
		func(batch []models.Author) error {
			insertBatchSizes = append(insertBatchSizes, len(batch))
			return nil
		},
	)

	if err != nil {
		t.Fatalf("resolveAuthorIDs error = %v, want nil", err)
	}
	if listCalls != 2 {
		t.Errorf("batch query calls = %d, want 2", listCalls)
	}
	if len(insertBatchSizes) != 2 || insertBatchSizes[0] != 2 || insertBatchSizes[1] != 2 {
		t.Errorf("insert batch sizes = %v, want [2 2]", insertBatchSizes)
	}
	for name, wantID := range map[string]uint64{
		"李白": 101, "杜甫": 102, "白居易": 103, "苏轼": 104, "辛弃疾": 105,
	} {
		if authorIDs[name] != wantID {
			t.Errorf("author %q ID = %d, want %d", name, authorIDs[name], wantID)
		}
	}
}

func TestCollectCiTuneNamesTrimsDeduplicatesAndSorts(t *testing.T) {
	poems := []poem{
		{Rhythmic: " 念奴娇 "},
		{Rhythmic: "水调歌头"},
		{Rhythmic: "念奴娇"},
		{Title: "静夜思"},
	}

	names := collectCiTuneNames(poems)

	assertStringSliceEqual(t, names, []string{"念奴娇", "水调歌头"})
}

func TestBuildMissingCiTunesExcludesExistingNames(t *testing.T) {
	names := []string{"念奴娇", "水调歌头", "蝶恋花"}
	existingIDs := map[string]uint64{"水调歌头": 101}

	missing := buildMissingCiTunes(names, existingIDs)

	if len(missing) != 2 {
		t.Fatalf("missing ci tune count = %d, want 2", len(missing))
	}
	if missing[0].Name != "念奴娇" || missing[1].Name != "蝶恋花" {
		t.Fatalf("missing ci tune names = [%q %q], want [念奴娇 蝶恋花]", missing[0].Name, missing[1].Name)
	}
}

func TestMergeCiTuneIDsUsesNormalizedValidRecords(t *testing.T) {
	ciTuneIDs := map[string]uint64{"已有词牌": 1}
	mergeCiTuneIDs(ciTuneIDs, []models.CiTune{
		{BaseModel: models.BaseModel{ID: 101}, Name: " 念奴娇 "},
		{Name: "无 ID 词牌"},
		{BaseModel: models.BaseModel{ID: 102}, Name: " "},
	})

	if ciTuneIDs["已有词牌"] != 1 || ciTuneIDs["念奴娇"] != 101 {
		t.Fatalf("ci tune IDs = %v, want existing and normalized IDs", ciTuneIDs)
	}
	if _, exists := ciTuneIDs["无 ID 词牌"]; exists {
		t.Fatal("ci tune IDs contains record without ID")
	}
}

func TestResolveCiTuneIDsUsesBatchQueriesAndBatchInserts(t *testing.T) {
	names := []string{"卜算子", "念奴娇", "水调歌头", "蝶恋花", "鹧鸪天"}
	listCalls := 0
	insertBatchSizes := make([]int, 0)

	ciTuneIDs, err := resolveCiTuneIDs(
		names,
		2,
		func(gotNames []string) ([]models.CiTune, error) {
			listCalls++
			if len(gotNames) != len(names) {
				t.Fatalf("batch query name count = %d, want %d", len(gotNames), len(names))
			}
			if listCalls == 1 {
				return []models.CiTune{{BaseModel: models.BaseModel{ID: 103}, Name: "水调歌头"}}, nil
			}
			return []models.CiTune{
				{BaseModel: models.BaseModel{ID: 101}, Name: "卜算子"},
				{BaseModel: models.BaseModel{ID: 102}, Name: "念奴娇"},
				{BaseModel: models.BaseModel{ID: 103}, Name: "水调歌头"},
				{BaseModel: models.BaseModel{ID: 104}, Name: "蝶恋花"},
				{BaseModel: models.BaseModel{ID: 105}, Name: "鹧鸪天"},
			}, nil
		},
		func(batch []models.CiTune) error {
			insertBatchSizes = append(insertBatchSizes, len(batch))
			return nil
		},
	)

	if err != nil {
		t.Fatalf("resolveCiTuneIDs error = %v, want nil", err)
	}
	if listCalls != 2 {
		t.Errorf("batch query calls = %d, want 2", listCalls)
	}
	if len(insertBatchSizes) != 2 || insertBatchSizes[0] != 2 || insertBatchSizes[1] != 2 {
		t.Errorf("insert batch sizes = %v, want [2 2]", insertBatchSizes)
	}
	for name, wantID := range map[string]uint64{
		"卜算子": 101, "念奴娇": 102, "水调歌头": 103, "蝶恋花": 104, "鹧鸪天": 105,
	} {
		if ciTuneIDs[name] != wantID {
			t.Errorf("ci tune %q ID = %d, want %d", name, ciTuneIDs[name], wantID)
		}
	}
}

func TestBuildPoemModelsAssociatesCiTuneFromRhythmic(t *testing.T) {
	rawPoems := []poem{
		{Author: "苏轼", Rhythmic: " 念奴娇 ", Paragraphs: []string{"大江东去，浪淘尽。"}},
		{Author: "李白", Title: "静夜思", Paragraphs: []string{"床前明月光。"}},
	}
	authorIDs := map[string]uint64{"苏轼": 11, "李白": 12}
	ciTuneIDs := map[string]uint64{"念奴娇": 21}

	poemModels, skippedPoems := buildPoemModels(rawPoems, authorIDs, 31, ciTuneIDs)

	if len(skippedPoems) != 0 || len(poemModels) != 2 {
		t.Fatalf("models = %d, skipped = %d, want 2 and 0", len(poemModels), len(skippedPoems))
	}
	if poemModels[0].Title != "念奴娇" {
		t.Errorf("ci title = %q, want 念奴娇", poemModels[0].Title)
	}
	if poemModels[0].CiTuneID == nil || *poemModels[0].CiTuneID != 21 {
		t.Errorf("ci tune ID = %v, want 21", poemModels[0].CiTuneID)
	}
	if poemModels[1].CiTuneID != nil {
		t.Errorf("poem ci tune ID = %v, want nil", poemModels[1].CiTuneID)
	}
}

func TestFindPoemsWithMissingAuthorIDs(t *testing.T) {
	poemModels := []models.Poem{
		{Title: "有效诗歌", AuthorID: 101},
		{Title: "缺失作者诗歌", AuthorID: 404},
		{Title: "零作者诗歌", AuthorID: 0},
	}
	existingAuthorIDs := map[uint64]struct{}{101: {}}

	missingPoems := findPoemsWithMissingAuthorIDs(poemModels, existingAuthorIDs)

	if len(missingPoems) != 2 {
		t.Fatalf("missing poem count = %d, want 2", len(missingPoems))
	}
	if missingPoems[0].Title != "缺失作者诗歌" || missingPoems[1].Title != "零作者诗歌" {
		t.Fatalf("missing poem titles = [%q %q], want [缺失作者诗歌 零作者诗歌]", missingPoems[0].Title, missingPoems[1].Title)
	}
}

func TestFormatPgErrorPrintsPostgresFields(t *testing.T) {
	err := &pgconn.PgError{
		Severity:       "ERROR",
		Code:           "23503",
		Message:        "insert or update on table \"poem\" violates foreign key constraint",
		Detail:         "Key (author_id)=(123) is not present in table \"author\".",
		Hint:           "Insert the author first.",
		SchemaName:     "public",
		TableName:      "poem",
		ColumnName:     "author_id",
		ConstraintName: "fk_poem_author",
		Where:          "SQL statement",
	}

	message := formatPgError(err)

	for _, want := range []string{
		"PostgreSQL error",
		"severity=ERROR",
		"code=23503",
		"message=insert or update on table \"poem\" violates foreign key constraint",
		"detail=Key (author_id)=(123) is not present in table \"author\".",
		"hint=Insert the author first.",
		"schema=public",
		"table=poem",
		"column=author_id",
		"constraint=fk_poem_author",
		"where=SQL statement",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("formatted error %q does not contain %q", message, want)
		}
	}
}

func TestFormatPgErrorSupportsWrappedPgError(t *testing.T) {
	err := fmt.Errorf("batch insert failed: %w", &pgconn.PgError{
		Code:           "23505",
		Message:        "duplicate key value violates unique constraint",
		ConstraintName: "author_name_key",
	})

	message := formatPgError(err)

	if !strings.Contains(message, "code=23505") {
		t.Fatalf("formatted error %q does not contain wrapped pg error code", message)
	}
	if !strings.Contains(message, "constraint=author_name_key") {
		t.Fatalf("formatted error %q does not contain wrapped pg error constraint", message)
	}
}

func TestFormatPgErrorFallsBackToRegularError(t *testing.T) {
	err := errors.New("regular error")

	message := formatPgError(err)

	if message != "error: regular error" {
		t.Fatalf("formatted error = %q, want %q", message, "error: regular error")
	}
}

func TestWalkDirectChildFilesVisitsOnlyDirectFiles(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "first.json"), "first")
	mustWriteFile(t, filepath.Join(root, "second.json"), "second")
	mustWriteFile(t, filepath.Join(root, "nested", "ignored.json"), "ignored")

	var visited []string
	err := utils.WalkDirectChildFiles(root, func(path string, entry os.DirEntry) error {
		if entry.IsDir() {
			t.Fatalf("visited directory %q, want files only", path)
		}
		visited = append(visited, filepath.Base(path))
		return nil
	})

	if err != nil {
		t.Fatalf("walkDirectChildFiles error = %v, want nil", err)
	}
	sort.Strings(visited)
	assertStringSliceEqual(t, visited, []string{"first.json", "second.json"})
}

func TestWalkAllFilesVisitsNestedFiles(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "root.json"), "root")
	mustWriteFile(t, filepath.Join(root, "nested", "child.json"), "child")
	mustWriteFile(t, filepath.Join(root, "nested", "deep", "grandchild.json"), "grandchild")

	var visited []string
	err := utils.WalkAllFiles(root, func(path string, entry os.DirEntry) error {
		if entry.IsDir() {
			t.Fatalf("visited directory %q, want files only", path)
		}
		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			t.Fatalf("filepath.Rel error = %v, want nil", err)
		}
		visited = append(visited, filepath.ToSlash(relativePath))
		return nil
	})

	if err != nil {
		t.Fatalf("walkAllFiles error = %v, want nil", err)
	}
	sort.Strings(visited)
	assertStringSliceEqual(t, visited, []string{
		"nested/child.json",
		"nested/deep/grandchild.json",
		"root.json",
	})
}

func TestWalkAllFilesStopsWhenCallbackReturnsError(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "first.json"), "first")
	mustWriteFile(t, filepath.Join(root, "second.json"), "second")
	expectedErr := errors.New("stop walking")

	err := utils.WalkAllFiles(root, func(path string, entry os.DirEntry) error {
		return expectedErr
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("walkAllFiles error = %v, want %v", err, expectedErr)
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

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("os.MkdirAll error = %v, want nil", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile error = %v, want nil", err)
	}
}

func assertStringSliceEqual(t *testing.T, got []string, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d; got %v", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("slice[%d] = %q, want %q; got %v", i, got[i], want[i], got)
		}
	}
}
