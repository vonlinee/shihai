package services

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"shihai/internal/dto"
	"shihai/internal/models"
)

func TestCreatePoemRejectsMissingDynastyBeforeCreatingPoet(t *testing.T) {
	poemRepo := &fakePoemRepository{}
	dynastyRepo := &fakeDynastyRepository{}
	authorRepo := &fakeAuthorRepository{}
	poetRepo := &fakePoetRepository{}
	service := NewPoemService(poemRepo, dynastyRepo, authorRepo, poetRepo)

	_, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:      "测试",
		Content:    []string{"1111"},
		DynastyID:  706545104525070300,
		AuthorName: "李白",
	})

	if err == nil || !strings.Contains(err.Error(), "dynasty not found") {
		t.Fatalf("CreatePoem error = %v, want dynasty not found", err)
	}
	if poetRepo.createCalls != 0 {
		t.Fatalf("poet create calls = %d, want 0", poetRepo.createCalls)
	}
	if authorRepo.createCalls != 0 {
		t.Fatalf("author create calls = %d, want 0", authorRepo.createCalls)
	}
	if poemRepo.createCalls != 0 {
		t.Fatalf("poem create calls = %d, want 0", poemRepo.createCalls)
	}
}

func TestCreatePoemPrefersDynastyNameOverStaleDynastyID(t *testing.T) {
	poemRepo := &fakePoemRepository{}
	dynastyRepo := &fakeDynastyRepository{
		existingByName: map[string]*models.Dynasty{
			"Tang": {BaseModel: models.BaseModel{ID: 1}, Name: "Tang"},
		},
	}
	authorRepo := &fakeAuthorRepository{}
	poetRepo := &fakePoetRepository{}
	service := NewPoemService(poemRepo, dynastyRepo, authorRepo, poetRepo)

	_, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:       "test",
		Content:     []string{"1111", "2222"},
		DynastyID:   706545104525070300,
		DynastyName: "Tang",
		AuthorName:  "Li Bai",
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	if authorRepo.createdName != "Li Bai" {
		t.Fatalf("created author name = %q, want Li Bai", authorRepo.createdName)
	}
	if poetRepo.createdAuthorID != 1 {
		t.Fatalf("poet author ID = %d, want 1", poetRepo.createdAuthorID)
	}
	if poetRepo.createdDynastyID != 1 {
		t.Fatalf("poet dynasty ID = %d, want 1", poetRepo.createdDynastyID)
	}
	if poemRepo.createdAuthorID != 1 {
		t.Fatalf("poem author ID = %d, want 1", poemRepo.createdAuthorID)
	}
	if poemRepo.createdDynastyID != 1 {
		t.Fatalf("poem dynasty ID = %d, want 1", poemRepo.createdDynastyID)
	}
	assertStringSliceEqual(t, poemRepo.createdContent, []string{"1111", "2222"})
}

func TestGetGenreCategoriesGroupsPoemTypesByCategory(t *testing.T) {
	poemRepo := &fakePoemRepository{
		poemTypes: []models.PoemType{
			{Name: "五言绝句", Category: "诗", Lines: intPtrForTest(4), CharsPerLine: intPtrForTest(5), Description: "四句，每句五字"},
			{Name: "小令", Category: "词", Description: "篇幅较短的词调"},
			{Name: "散文", Category: "文", Description: "不受韵律严格约束"},
			{Name: "未分类", Description: "缺少分类"},
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	categories, err := service.GetGenreCategories()

	if err != nil {
		t.Fatalf("GetGenreCategories error = %v, want nil", err)
	}
	if len(categories) != 4 {
		t.Fatalf("category count = %d, want 4", len(categories))
	}
	assertGenreCategory(t, categories[0], "诗", []string{"五言绝句"})
	assertGenreCategory(t, categories[1], "词", []string{"小令"})
	assertGenreCategory(t, categories[2], "文", []string{"散文"})
	assertGenreCategory(t, categories[3], "其他", []string{"未分类"})
	if categories[0].Genres[0].Lines == nil || *categories[0].Genres[0].Lines != 4 {
		t.Fatalf("poem genre lines = %v, want 4", categories[0].Genres[0].Lines)
	}
	if categories[0].Genres[0].CharsPerLine == nil || *categories[0].Genres[0].CharsPerLine != 5 {
		t.Fatalf("poem genre chars per line = %v, want 5", categories[0].Genres[0].CharsPerLine)
	}
}

func TestCreatePoemTypeStoresReferenceData(t *testing.T) {
	poemRepo := &fakePoemRepository{}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})
	lines := 4
	charsPerLine := 5

	resp, err := service.CreatePoemType(&dto.PoemTypeCreateRequest{
		Name:         "五言绝句",
		Category:     "诗",
		Lines:        &lines,
		CharsPerLine: &charsPerLine,
		Description:  "四句，每句五字",
	})

	if err != nil {
		t.Fatalf("CreatePoemType error = %v, want nil", err)
	}
	if poemRepo.createdPoemType == nil {
		t.Fatal("created poem type = nil, want saved poem type")
	}
	if poemRepo.createdPoemType.Name != "五言绝句" || poemRepo.createdPoemType.Category != "诗" {
		t.Fatalf("created poem type = %#v, want 五言绝句 under 诗", poemRepo.createdPoemType)
	}
	if resp.ID != 1 || resp.Name != "五言绝句" || resp.Category != "诗" {
		t.Fatalf("response poem type = %#v, want created response", resp)
	}
}

func TestUpdatePoemTypeCanClearLineConstraints(t *testing.T) {
	poemRepo := &fakePoemRepository{
		poemTypeByID: &models.PoemType{
			BaseModel:    models.BaseModel{ID: 9},
			Name:         "五言绝句",
			Category:     "诗",
			Lines:        intPtrForTest(4),
			CharsPerLine: intPtrForTest(5),
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.UpdatePoemType(9, &dto.PoemTypeUpdateRequest{
		Name:        "古体诗",
		Category:    "诗",
		Description: "不限句数",
	})

	if err != nil {
		t.Fatalf("UpdatePoemType error = %v, want nil", err)
	}
	if poemRepo.updatedPoemType == nil {
		t.Fatal("updated poem type = nil, want saved poem type")
	}
	if poemRepo.updatedPoemType.Lines != nil || poemRepo.updatedPoemType.CharsPerLine != nil {
		t.Fatalf("updated line constraints = %v/%v, want nil/nil", poemRepo.updatedPoemType.Lines, poemRepo.updatedPoemType.CharsPerLine)
	}
	if resp.Name != "古体诗" || resp.Description != "不限句数" {
		t.Fatalf("response poem type = %#v, want updated response", resp)
	}
}

func TestBatchDeletePoemTypesRejectsEmptyIDs(t *testing.T) {
	service := NewPoemService(&fakePoemRepository{}, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	err := service.BatchDeletePoemTypes(nil)

	if err == nil || !strings.Contains(err.Error(), "ids cannot be empty") {
		t.Fatalf("BatchDeletePoemTypes error = %v, want ids cannot be empty", err)
	}
}

func TestGetPoemByIDReturnsContentArray(t *testing.T) {
	poemRepo := &fakePoemRepository{
		existingPoem: &models.Poem{
			BaseModel: models.BaseModel{ID: 1},
			Title:     "静夜思",
			Content:   []string{"床前明月光，", "疑是地上霜。"},
			Pingze:    []string{"平平平仄平", "平仄仄仄平"},
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.GetPoemByID(1)

	if err != nil {
		t.Fatalf("GetPoemByID error = %v, want nil", err)
	}
	assertStringSliceEqual(t, resp.Content, []string{"床前明月光，", "疑是地上霜。"})
	assertStringSliceEqual(t, resp.Pingze, []string{"平平平仄平", "平仄仄仄平"})
}

func TestCreatePoemStoresPingzeByContentLine(t *testing.T) {
	poemRepo := &fakePoemRepository{}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:   "静夜思",
		Content: []string{"床前明月光", "疑是地上霜"},
		Pingze:  []string{"平平平仄平", "平仄仄仄平"},
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	assertStringSliceEqual(t, poemRepo.createdPingze, []string{"平平平仄平", "平仄仄仄平"})
	assertStringSliceEqual(t, resp.Pingze, []string{"平平平仄平", "平仄仄仄平"})
}

func TestCreatePoemStoresGenreCategory(t *testing.T) {
	poemRepo := &fakePoemRepository{}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:         "静夜思",
		Content:       []string{"床前明月光"},
		GenreCategory: "诗",
		Genre:         "五言绝句",
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	if poemRepo.createdGenreCategory != "诗" {
		t.Fatalf("created genre category = %q, want 诗", poemRepo.createdGenreCategory)
	}
	if resp.GenreCategory != "诗" {
		t.Fatalf("response genre category = %q, want 诗", resp.GenreCategory)
	}
}

func TestCreatePoemInfersGenreCategoryFromReferenceData(t *testing.T) {
	poemRepo := &fakePoemRepository{
		poemTypes: []models.PoemType{{Name: "小令", Category: "词"}},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	_, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:   "测试",
		Content: []string{"我爱你"},
		Genre:   "小令",
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	if poemRepo.createdGenreCategory != "词" {
		t.Fatalf("created genre category = %q, want 词", poemRepo.createdGenreCategory)
	}
}

func TestCreatePoemAutoResolvesCiTuneByTitle(t *testing.T) {
	poemRepo := &fakePoemRepository{
		ciTunes: []models.CiTune{
			{BaseModel: models.BaseModel{ID: 11}, Name: "水调歌头"},
			{BaseModel: models.BaseModel{ID: 12}, Name: "水调"},
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:         "水调歌头·明月几时有",
		Content:       []string{"明月几时有"},
		GenreCategory: "词",
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	if poemRepo.createdCiTuneID == nil || *poemRepo.createdCiTuneID != 11 {
		t.Fatalf("created ci tune ID = %v, want 11", poemRepo.createdCiTuneID)
	}
	if resp.CiTuneID == nil || *resp.CiTuneID != 11 {
		t.Fatalf("response ci tune ID = %v, want 11", resp.CiTuneID)
	}
}

func TestParseCiTuneTitlePrefersLongestAlias(t *testing.T) {
	poemRepo := &fakePoemRepository{
		ciTunes: []models.CiTune{
			{BaseModel: models.BaseModel{ID: 11}, Name: "水调"},
			{BaseModel: models.BaseModel{ID: 12}, Name: "水调歌头", Aliases: []string{"元会曲"}},
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.ParseCiTuneTitle("水调歌头 明月几时有")

	if err != nil {
		t.Fatalf("ParseCiTuneTitle error = %v, want nil", err)
	}
	if resp.CiTune == nil || resp.CiTune.ID != 12 {
		t.Fatalf("matched ci tune = %#v, want ID 12", resp.CiTune)
	}
	if resp.MatchedName != "水调歌头" {
		t.Fatalf("matched name = %q, want 水调歌头", resp.MatchedName)
	}
}

func TestCreatePoemClearsCiTuneForNonCiCategory(t *testing.T) {
	poemRepo := &fakePoemRepository{
		ciTuneByID: &models.CiTune{BaseModel: models.BaseModel{ID: 11}, Name: "水调歌头"},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	_, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:         "静夜思",
		Content:       []string{"床前明月光"},
		GenreCategory: "诗",
		CiTuneID:      dto.RequestID(11),
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	if poemRepo.createdCiTuneID != nil {
		t.Fatalf("created ci tune ID = %v, want nil", poemRepo.createdCiTuneID)
	}
}

func TestCreatePoemRejectsInvalidPingzeMark(t *testing.T) {
	service := NewPoemService(&fakePoemRepository{}, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	_, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:   "静夜思",
		Content: []string{"床前明月光"},
		Pingze:  []string{"平平中仄平"},
	})

	if err == nil || !strings.Contains(err.Error(), "pingze can only contain") {
		t.Fatalf("CreatePoem error = %v, want invalid pingze error", err)
	}
}

func TestCreatePoemAutoRecognizesPingzeWhenMissing(t *testing.T) {
	poemRepo := &fakePoemRepository{}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:   "测试",
		Content: []string{"我爱你"},
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	assertStringSliceEqual(t, poemRepo.createdPingze, []string{"仄仄仄"})
	assertStringSliceEqual(t, resp.Pingze, []string{"仄仄仄"})
}

func TestUpdatePoemAlignsPingzeWhenContentChanges(t *testing.T) {
	poemRepo := &fakePoemRepository{
		existingPoem: &models.Poem{
			BaseModel: models.BaseModel{ID: 1},
			Title:     "静夜思",
			Content:   []string{"床前明月光"},
			Pingze:    []string{"平平平仄平"},
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.UpdatePoem(1, &dto.PoemUpdateRequest{
		Content: []string{"我爱你", "中国"},
	})

	if err != nil {
		t.Fatalf("UpdatePoem error = %v, want nil", err)
	}
	if poemRepo.updatedPoem == nil {
		t.Fatal("updated poem = nil, want saved poem")
	}
	assertStringSliceEqual(t, poemRepo.updatedPoem.Pingze, []string{"仄仄仄", "平平"})
	assertStringSliceEqual(t, resp.Pingze, []string{"仄仄仄", "平平"})
}

func TestUpdatePoemStoresGenreCategory(t *testing.T) {
	poemRepo := &fakePoemRepository{
		existingPoem: &models.Poem{
			BaseModel: models.BaseModel{ID: 1},
			Genre:     "五言绝句",
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.UpdatePoem(1, &dto.PoemUpdateRequest{
		GenreCategory: "诗",
		Genre:         "七言律诗",
	})

	if err != nil {
		t.Fatalf("UpdatePoem error = %v, want nil", err)
	}
	if poemRepo.updatedPoem == nil {
		t.Fatal("updated poem = nil, want saved poem")
	}
	if poemRepo.updatedPoem.GenreCategory != "诗" {
		t.Fatalf("updated genre category = %q, want 诗", poemRepo.updatedPoem.GenreCategory)
	}
	if resp.GenreCategory != "诗" {
		t.Fatalf("response genre category = %q, want 诗", resp.GenreCategory)
	}
}

func TestGetPoemByIDReturnsAllPersistedAnnotations(t *testing.T) {
	poemRepo := &fakePoemRepository{
		existingPoem: &models.Poem{
			BaseModel: models.BaseModel{ID: 1},
			Title:     "test",
			Content:   []string{"abcdef"},
		},
	}
	annotationRepo := &fakePoemAnnotationRepository{
		annotations: []models.PoemAnnotation{
			{
				BaseModel:    models.BaseModel{ID: 11},
				PoemID:       1,
				TargetField:  models.PoemAnnotationTargetContent,
				StartLine:    0,
				StartOffset:  1,
				EndLine:      0,
				EndOffset:    3,
				SelectedText: "bc",
				Content:      "persisted annotation",
			},
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})
	service.SetPoemAnnotationRepository(annotationRepo)

	resp, err := service.GetPoemByID(1)

	if err != nil {
		t.Fatalf("GetPoemByID error = %v, want nil", err)
	}
	if len(resp.Annotations) != 1 {
		t.Fatalf("len(resp.Annotations) = %d, want 1", len(resp.Annotations))
	}
	if resp.Annotations[0].SelectedText != "bc" {
		t.Fatalf("annotation = %#v, want selected text bc", resp.Annotations[0])
	}
}

func TestGetPoetListReturnsPaginatedPoets(t *testing.T) {
	poetRepo := &fakePoetRepository{
		listPoets: []models.Poet{
			{
				BaseModel: models.BaseModel{ID: 1},
				AuthorID:  11,
				Author:    models.Author{Name: "Li Bai"},
			},
		},
		listTotal: 12,
	}
	service := NewPoemService(&fakePoemRepository{}, &fakeDynastyRepository{}, &fakeAuthorRepository{}, poetRepo)

	poets, total, err := service.GetPoetList("Li", 9, 2, 3)

	if err != nil {
		t.Fatalf("GetPoetList error = %v, want nil", err)
	}
	if total != 12 {
		t.Fatalf("total = %d, want 12", total)
	}
	if poetRepo.listKeyword != "Li" {
		t.Fatalf("keyword = %q, want Li", poetRepo.listKeyword)
	}
	if poetRepo.listDynastyID != 9 {
		t.Fatalf("dynastyID = %d, want 9", poetRepo.listDynastyID)
	}
	if poetRepo.listPage != 2 {
		t.Fatalf("page = %d, want 2", poetRepo.listPage)
	}
	if poetRepo.listPageSize != 3 {
		t.Fatalf("pageSize = %d, want 3", poetRepo.listPageSize)
	}
	if len(poets) != 1 || poets[0].Name != "Li Bai" {
		t.Fatalf("poets = %#v, want one Li Bai", poets)
	}
}

func TestBatchDeleteResourcesPassesIDsToRepositories(t *testing.T) {
	poemRepo := &fakePoemRepository{}
	dynastyRepo := &fakeDynastyRepository{}
	poetRepo := &fakePoetRepository{}
	service := NewPoemService(poemRepo, dynastyRepo, &fakeAuthorRepository{}, poetRepo)

	if err := service.BatchDeletePoems([]uint64{1, 2, 3}); err != nil {
		t.Fatalf("BatchDeletePoems error = %v, want nil", err)
	}
	assertUint64SliceEqual(t, poemRepo.batchDeletedIDs, []uint64{1, 2, 3})

	if err := service.BatchDeleteDynasties([]uint64{4, 5}); err != nil {
		t.Fatalf("BatchDeleteDynasties error = %v, want nil", err)
	}
	assertUint64SliceEqual(t, dynastyRepo.batchDeletedIDs, []uint64{4, 5})

	if err := service.BatchDeletePoets([]uint64{6, 7}); err != nil {
		t.Fatalf("BatchDeletePoets error = %v, want nil", err)
	}
	assertUint64SliceEqual(t, poetRepo.batchDeletedIDs, []uint64{6, 7})
}

func TestBatchDeleteResourcesRejectsEmptyIDs(t *testing.T) {
	service := NewPoemService(&fakePoemRepository{}, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	if err := service.BatchDeletePoems(nil); err == nil {
		t.Fatal("BatchDeletePoems error = nil, want error")
	}
	if err := service.BatchDeleteDynasties(nil); err == nil {
		t.Fatal("BatchDeleteDynasties error = nil, want error")
	}
	if err := service.BatchDeletePoets(nil); err == nil {
		t.Fatal("BatchDeletePoets error = nil, want error")
	}
}

func TestUpdatePoetReturnsChangedDynasty(t *testing.T) {
	poetRepo := &fakePoetRepository{
		existingPoet: &models.Poet{
			BaseModel: models.BaseModel{ID: 1},
			AuthorID:  11,
			Author:    models.Author{BaseModel: models.BaseModel{ID: 11}, Name: "李白"},
			DynastyID: 100,
			Dynasty:   models.Dynasty{BaseModel: models.BaseModel{ID: 100}, Name: "唐"},
		},
	}
	dynastyRepo := &fakeDynastyRepository{
		existingByID: map[uint64]*models.Dynasty{
			200: {BaseModel: models.BaseModel{ID: 200}, Name: "宋", Period: "960-1279"},
		},
	}
	service := NewPoemService(&fakePoemRepository{}, dynastyRepo, &fakeAuthorRepository{}, poetRepo)

	resp, err := service.UpdatePoet(1, &dto.PoetUpdateRequest{DynastyID: dto.RequestID(200)})

	if err != nil {
		t.Fatalf("UpdatePoet error = %v, want nil", err)
	}
	if poetRepo.updatedPoet == nil {
		t.Fatal("updated poet = nil, want saved poet")
	}
	if poetRepo.updatedPoet.DynastyID != 200 {
		t.Fatalf("updated DynastyID = %d, want 200", poetRepo.updatedPoet.DynastyID)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Marshal response error = %v, want nil", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("Unmarshal response error = %v, want nil", err)
	}
	dynasty, ok := body["dynasty"].(map[string]any)
	if !ok {
		t.Fatalf("dynasty = %#v, want object", body["dynasty"])
	}
	if dynasty["id"] != "200" {
		t.Fatalf("dynasty.id = %#v, want string ID 200", dynasty["id"])
	}
	if dynasty["name"] != "宋" {
		t.Fatalf("dynasty.name = %#v, want 宋", dynasty["name"])
	}
}

func TestUpdatePoemSyncsAnnotationsWhenProvided(t *testing.T) {
	poemRepo := &fakePoemRepository{
		existingPoem: &models.Poem{
			BaseModel: models.BaseModel{ID: 1},
			Title:     "静夜思",
			Content:   []string{"床前明月光", "疑是地上霜"},
		},
	}
	annotationRepo := &fakePoemAnnotationRepository{
		annotations: []models.PoemAnnotation{
			{
				BaseModel:    models.BaseModel{ID: 11},
				PoemID:       1,
				TargetField:  models.PoemAnnotationTargetContent,
				StartLine:    0,
				StartOffset:  0,
				EndLine:      0,
				EndOffset:    1,
				SelectedText: "床",
				Content:      "旧标注",
			},
			{
				BaseModel:    models.BaseModel{ID: 12},
				PoemID:       1,
				TargetField:  models.PoemAnnotationTargetContent,
				StartLine:    1,
				StartOffset:  0,
				EndLine:      1,
				EndOffset:    1,
				SelectedText: "疑",
				Content:      "应被删除",
			},
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})
	service.SetPoemAnnotationRepository(annotationRepo)
	annotations := []dto.PoemAnnotationUpsertRequest{
		{
			ID:           dto.RequestID(11),
			TargetField:  models.PoemAnnotationTargetContent,
			StartLine:    0,
			StartOffset:  2,
			EndLine:      0,
			EndOffset:    4,
			SelectedText: "明月",
			Content:      "更新标注",
		},
		{
			TargetField:  models.PoemAnnotationTargetContent,
			StartLine:    1,
			StartOffset:  2,
			EndLine:      1,
			EndOffset:    4,
			SelectedText: "地上",
			Content:      "新增标注",
		},
	}

	resp, err := service.UpdatePoem(1, &dto.PoemUpdateRequest{
		Content:     []string{"床前明月光", "疑是地上霜"},
		Annotations: &annotations,
	})

	if err != nil {
		t.Fatalf("UpdatePoem error = %v, want nil", err)
	}
	if annotationRepo.updated == nil || annotationRepo.updated.ID != 11 || annotationRepo.updated.SelectedText != "明月" {
		t.Fatalf("updated annotation = %#v, want ID 11 明月", annotationRepo.updated)
	}
	if annotationRepo.created == nil || annotationRepo.created.SelectedText != "地上" {
		t.Fatalf("created annotation = %#v, want 地上", annotationRepo.created)
	}
	if annotationRepo.deletedID != 12 {
		t.Fatalf("deleted annotation ID = %d, want 12", annotationRepo.deletedID)
	}
	if len(resp.Annotations) != 2 {
		t.Fatalf("len(resp.Annotations) = %d, want 2", len(resp.Annotations))
	}
	if resp.Annotations[0].SelectedText != "明月" || resp.Annotations[1].SelectedText != "地上" {
		t.Fatalf("resp annotations = %#v, want updated annotations in text order", resp.Annotations)
	}
}

func TestCreatePoemSyncsAnnotationsWhenProvided(t *testing.T) {
	poemRepo := &fakePoemRepository{}
	annotationRepo := &fakePoemAnnotationRepository{}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})
	service.SetPoemAnnotationRepository(annotationRepo)
	annotations := []dto.PoemAnnotationUpsertRequest{
		{
			TargetField:  models.PoemAnnotationTargetContent,
			StartLine:    0,
			StartOffset:  0,
			EndLine:      0,
			EndOffset:    2,
			SelectedText: "床前",
			Content:      "新建标注",
		},
	}

	resp, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:       "静夜思",
		Content:     []string{"床前明月光"},
		Annotations: &annotations,
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	if annotationRepo.created == nil || annotationRepo.created.SelectedText != "床前" {
		t.Fatalf("created annotation = %#v, want 床前", annotationRepo.created)
	}
	if len(resp.Annotations) != 1 || resp.Annotations[0].SelectedText != "床前" {
		t.Fatalf("resp annotations = %#v, want created annotation", resp.Annotations)
	}
}

type fakePoemRepository struct {
	createCalls           int
	createdAuthorID       uint64
	createdDynastyID      uint64
	createdContent        []string
	createdPingze         []string
	createdGenreCategory  string
	createdCiTuneID       *uint64
	updatedPoem           *models.Poem
	existingPoem          *models.Poem
	poemTypes             []models.PoemType
	createdPoemType       *models.PoemType
	poemTypeByID          *models.PoemType
	updatedPoemType       *models.PoemType
	deletedPoemTypeID     uint64
	batchDeletedTypeIDs   []uint64
	ciTunes               []models.CiTune
	ciTuneByID            *models.CiTune
	createdCiTune         *models.CiTune
	updatedCiTune         *models.CiTune
	deletedCiTuneID       uint64
	batchDeletedCiTuneIDs []uint64
	batchDeletedIDs       []uint64
}

func (r *fakePoemRepository) Create(poem *models.Poem) error {
	r.createCalls++
	r.createdAuthorID = poem.AuthorID
	r.createdDynastyID = poem.DynastyID
	r.createdContent = append([]string(nil), poem.Content...)
	r.createdPingze = append([]string(nil), poem.Pingze...)
	r.createdGenreCategory = poem.GenreCategory
	if poem.CiTuneID != nil {
		ciTuneID := *poem.CiTuneID
		r.createdCiTuneID = &ciTuneID
	} else {
		r.createdCiTuneID = nil
	}
	poem.ID = 1
	return nil
}

func (r *fakePoemRepository) GetByID(id uint64) (*models.Poem, error) {
	if r.existingPoem != nil {
		return r.existingPoem, nil
	}
	return &models.Poem{BaseModel: models.BaseModel{ID: id}}, nil
}

func (r *fakePoemRepository) Update(poem *models.Poem) error {
	copied := *poem
	copied.Content = append([]string(nil), poem.Content...)
	copied.Pingze = append([]string(nil), poem.Pingze...)
	r.updatedPoem = &copied
	return nil
}

func (r *fakePoemRepository) Delete(id uint64) error {
	return nil
}

func (r *fakePoemRepository) BatchDelete(ids []uint64) error {
	r.batchDeletedIDs = append([]uint64(nil), ids...)
	return nil
}

func (r *fakePoemRepository) List(page, pageSize int, keyword, dynasty, author, genre string) ([]models.Poem, int64, error) {
	return nil, 0, nil
}

func (r *fakePoemRepository) IncrementViews(id uint64) error {
	return nil
}

func (r *fakePoemRepository) IncrementLikes(id uint64) error {
	return nil
}

func (r *fakePoemRepository) GetRandom(limit int) ([]models.Poem, error) {
	return nil, nil
}

func (r *fakePoemRepository) DistinctGenres() ([]string, error) {
	return nil, nil
}

func (r *fakePoemRepository) ListPoemTypes() ([]models.PoemType, error) {
	return r.poemTypes, nil
}

func (r *fakePoemRepository) CreatePoemType(poemType *models.PoemType) error {
	copied := *poemType
	copied.ID = 1
	r.createdPoemType = &copied
	poemType.ID = 1
	return nil
}

func (r *fakePoemRepository) GetPoemTypeByID(id uint64) (*models.PoemType, error) {
	if r.poemTypeByID != nil {
		return r.poemTypeByID, nil
	}
	return nil, errors.New("not found")
}

func (r *fakePoemRepository) UpdatePoemType(poemType *models.PoemType) error {
	copied := *poemType
	r.updatedPoemType = &copied
	return nil
}

func (r *fakePoemRepository) DeletePoemType(id uint64) error {
	r.deletedPoemTypeID = id
	return nil
}

func (r *fakePoemRepository) BatchDeletePoemTypes(ids []uint64) error {
	r.batchDeletedTypeIDs = append([]uint64(nil), ids...)
	return nil
}

func (r *fakePoemRepository) ListCiTunes() ([]models.CiTune, error) {
	return r.ciTunes, nil
}

func (r *fakePoemRepository) CreateCiTune(ciTune *models.CiTune) error {
	copied := *ciTune
	copied.ID = 1
	copied.Aliases = append([]string(nil), ciTune.Aliases...)
	r.createdCiTune = &copied
	ciTune.ID = 1
	return nil
}

func (r *fakePoemRepository) GetCiTuneByID(id uint64) (*models.CiTune, error) {
	if r.ciTuneByID != nil {
		return r.ciTuneByID, nil
	}
	for i := range r.ciTunes {
		if r.ciTunes[i].ID == id {
			return &r.ciTunes[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (r *fakePoemRepository) UpdateCiTune(ciTune *models.CiTune) error {
	copied := *ciTune
	copied.Aliases = append([]string(nil), ciTune.Aliases...)
	r.updatedCiTune = &copied
	return nil
}

func (r *fakePoemRepository) DeleteCiTune(id uint64) error {
	r.deletedCiTuneID = id
	return nil
}

func (r *fakePoemRepository) BatchDeleteCiTunes(ids []uint64) error {
	r.batchDeletedCiTuneIDs = append([]uint64(nil), ids...)
	return nil
}

type fakeDynastyRepository struct {
	existingByID    map[uint64]*models.Dynasty
	existingByName  map[string]*models.Dynasty
	batchDeletedIDs []uint64
}

func (r *fakeDynastyRepository) Create(dynasty *models.Dynasty) error {
	dynasty.ID = 1
	return nil
}

func (r *fakeDynastyRepository) GetByID(id uint64) (*models.Dynasty, error) {
	if r.existingByID != nil {
		if dynasty, ok := r.existingByID[id]; ok {
			return dynasty, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *fakeDynastyRepository) GetByName(name string) (*models.Dynasty, error) {
	if r.existingByName != nil {
		if dynasty, ok := r.existingByName[name]; ok {
			return dynasty, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *fakeDynastyRepository) Update(dynasty *models.Dynasty) error {
	return nil
}

func (r *fakeDynastyRepository) Delete(id uint64) error {
	return nil
}

func (r *fakeDynastyRepository) BatchDelete(ids []uint64) error {
	r.batchDeletedIDs = append([]uint64(nil), ids...)
	return nil
}

func (r *fakeDynastyRepository) List() ([]models.Dynasty, error) {
	return nil, nil
}

type fakeAuthorRepository struct {
	createCalls int
	createdName string
}

func (r *fakeAuthorRepository) Create(author *models.Author) error {
	r.createCalls++
	r.createdName = author.Name
	author.ID = 1
	return nil
}

func (r *fakeAuthorRepository) GetByID(id uint64) (*models.Author, error) {
	return &models.Author{BaseModel: models.BaseModel{ID: id}}, nil
}

func (r *fakeAuthorRepository) GetByName(name string) (*models.Author, error) {
	return nil, errors.New("not found")
}

func (r *fakeAuthorRepository) List(keyword string) ([]models.Author, error) {
	return nil, nil
}

func (r *fakeAuthorRepository) Update(author *models.Author) error {
	return nil
}

func (r *fakeAuthorRepository) Delete(id uint64) error {
	return nil
}

type fakePoetRepository struct {
	createCalls      int
	createdAuthorID  uint64
	createdDynastyID uint64
	existingPoet     *models.Poet
	updatedPoet      *models.Poet
	listKeyword      string
	listDynastyID    uint64
	listPage         int
	listPageSize     int
	listPoets        []models.Poet
	listTotal        int64
	batchDeletedIDs  []uint64
}

func (r *fakePoetRepository) Create(poet *models.Poet) error {
	r.createCalls++
	r.createdAuthorID = poet.AuthorID
	r.createdDynastyID = poet.DynastyID
	poet.ID = 1
	return nil
}

func (r *fakePoetRepository) GetByID(id uint64) (*models.Poet, error) {
	if r.existingPoet != nil {
		return r.existingPoet, nil
	}
	return &models.Poet{BaseModel: models.BaseModel{ID: id}}, nil
}

func (r *fakePoetRepository) GetByAuthorID(authorID uint64) (*models.Poet, error) {
	return nil, errors.New("not found")
}

func (r *fakePoetRepository) List(keyword string, dynastyID uint64, page, pageSize int) ([]models.Poet, int64, error) {
	r.listKeyword = keyword
	r.listDynastyID = dynastyID
	r.listPage = page
	r.listPageSize = pageSize
	return r.listPoets, r.listTotal, nil
}

func (r *fakePoetRepository) Update(poet *models.Poet) error {
	copied := *poet
	r.updatedPoet = &copied
	return nil
}

func (r *fakePoetRepository) Delete(id uint64) error {
	return nil
}

func (r *fakePoetRepository) BatchDelete(ids []uint64) error {
	r.batchDeletedIDs = append([]uint64(nil), ids...)
	return nil
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

func assertUint64SliceEqual(t *testing.T, got []uint64, want []uint64) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d; got %v", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("slice[%d] = %d, want %d; got %v", i, got[i], want[i], got)
		}
	}
}

func assertGenreCategory(t *testing.T, got dto.GenreCategoryResponse, wantName string, wantGenres []string) {
	t.Helper()

	if got.Name != wantName {
		t.Fatalf("category name = %q, want %q", got.Name, wantName)
	}
	gotGenres := make([]string, 0, len(got.Genres))
	for _, genre := range got.Genres {
		gotGenres = append(gotGenres, genre.Name)
	}
	assertStringSliceEqual(t, gotGenres, wantGenres)
}

func intPtrForTest(value int) *int {
	return &value
}
