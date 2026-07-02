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
	poetRepo := &fakePoetRepository{}
	service := NewPoemService(poemRepo, dynastyRepo, poetRepo)

	_, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:      "测试",
		Content:    "1111",
		DynastyID:  706545104525070300,
		AuthorName: "李白",
	})

	if err == nil || !strings.Contains(err.Error(), "dynasty not found") {
		t.Fatalf("CreatePoem error = %v, want dynasty not found", err)
	}
	if poetRepo.createCalls != 0 {
		t.Fatalf("poet create calls = %d, want 0", poetRepo.createCalls)
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
	poetRepo := &fakePoetRepository{}
	service := NewPoemService(poemRepo, dynastyRepo, poetRepo)

	_, err := service.CreatePoem(&dto.PoemCreateRequest{
		Title:       "test",
		Content:     "1111",
		DynastyID:   706545104525070300,
		DynastyName: "Tang",
		AuthorName:  "Li Bai",
	})

	if err != nil {
		t.Fatalf("CreatePoem error = %v, want nil", err)
	}
	if poetRepo.createdDynastyID != 1 {
		t.Fatalf("poet dynasty ID = %d, want 1", poetRepo.createdDynastyID)
	}
	if poemRepo.createdDynastyID != 1 {
		t.Fatalf("poem dynasty ID = %d, want 1", poemRepo.createdDynastyID)
	}
}

type fakePoemRepository struct {
	createCalls      int
	createdDynastyID uint64
}

func (r *fakePoemRepository) Create(poem *models.Poem) error {
	r.createCalls++
	r.createdDynastyID = poem.DynastyID
	poem.ID = 1
	return nil
}

func (r *fakePoemRepository) GetByID(id uint64) (*models.Poem, error) {
	return &models.Poem{BaseModel: models.BaseModel{ID: id}}, nil
}

func (r *fakePoemRepository) Update(poem *models.Poem) error {
	return nil
}

func (r *fakePoemRepository) Delete(id uint64) error {
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
	existingByID   map[uint64]*models.Dynasty
	existingByName map[string]*models.Dynasty
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

func (r *fakeDynastyRepository) List() ([]models.Dynasty, error) {
	return nil, nil
}

type fakePoetRepository struct {
	createCalls      int
	createdDynastyID uint64
}

func (r *fakePoetRepository) Create(poet *models.Poet) error {
	r.createCalls++
	r.createdDynastyID = poet.DynastyID
	poet.ID = 1
	return nil
}

func (r *fakePoetRepository) GetByID(id uint64) (*models.Poet, error) {
	return &models.Poet{BaseModel: models.BaseModel{ID: id}}, nil
}

func (r *fakePoetRepository) GetByName(name string) (*models.Poet, error) {
	return nil, errors.New("not found")
}

func (r *fakePoetRepository) List(keyword string) ([]models.Poet, error) {
	return nil, nil
}

func (r *fakePoetRepository) Update(poet *models.Poet) error {
	return nil
}

func (r *fakePoetRepository) Delete(id uint64) error {
	return nil
}
