package services

import (
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

func TestGetPoemByIDReturnsContentArray(t *testing.T) {
	poemRepo := &fakePoemRepository{
		existingPoem: &models.Poem{
			BaseModel: models.BaseModel{ID: 1},
			Title:     "静夜思",
			Content:   []string{"床前明月光，", "疑是地上霜。"},
		},
	}
	service := NewPoemService(poemRepo, &fakeDynastyRepository{}, &fakeAuthorRepository{}, &fakePoetRepository{})

	resp, err := service.GetPoemByID(1)

	if err != nil {
		t.Fatalf("GetPoemByID error = %v, want nil", err)
	}
	assertStringSliceEqual(t, resp.Content, []string{"床前明月光，", "疑是地上霜。"})
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

	poets, total, err := service.GetPoetList("Li", 2, 3)

	if err != nil {
		t.Fatalf("GetPoetList error = %v, want nil", err)
	}
	if total != 12 {
		t.Fatalf("total = %d, want 12", total)
	}
	if poetRepo.listKeyword != "Li" {
		t.Fatalf("keyword = %q, want Li", poetRepo.listKeyword)
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

type fakePoemRepository struct {
	createCalls      int
	createdAuthorID  uint64
	createdDynastyID uint64
	createdContent   []string
	existingPoem     *models.Poem
	batchDeletedIDs  []uint64
}

func (r *fakePoemRepository) Create(poem *models.Poem) error {
	r.createCalls++
	r.createdAuthorID = poem.AuthorID
	r.createdDynastyID = poem.DynastyID
	r.createdContent = append([]string(nil), poem.Content...)
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
	listKeyword      string
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
	return &models.Poet{BaseModel: models.BaseModel{ID: id}}, nil
}

func (r *fakePoetRepository) GetByAuthorID(authorID uint64) (*models.Poet, error) {
	return nil, errors.New("not found")
}

func (r *fakePoetRepository) List(keyword string, page, pageSize int) ([]models.Poet, int64, error) {
	r.listKeyword = keyword
	r.listPage = page
	r.listPageSize = pageSize
	return r.listPoets, r.listTotal, nil
}

func (r *fakePoetRepository) Update(poet *models.Poet) error {
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
