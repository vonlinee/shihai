package services

import (
	"errors"
	"strings"
	"testing"

	"shihai/internal/dto"
	"shihai/internal/models"
)

func TestWorkCollectionServiceCreateCollectionRequiresTitle(t *testing.T) {
	service := NewWorkCollectionService(newFakeWorkCollectionStore(), newFakePoemStore())

	_, err := service.CreateCollection(&dto.WorkCollectionCreateRequest{Title: "   "})

	if err == nil || !strings.Contains(err.Error(), "title is required") {
		t.Fatalf("CreateCollection error = %v, want title required", err)
	}
}

func TestWorkCollectionServiceAddPoemItem(t *testing.T) {
	collections := newFakeWorkCollectionStore()
	poems := newFakePoemStore()
	poems.existing[10] = true
	collection := collections.addCollection("唐诗选读")
	service := NewWorkCollectionService(collections, poems)

	item, err := service.AddItem(collection.ID, &dto.WorkCollectionItemCreateRequest{
		WorkType:  "poem",
		WorkID:    10,
		SortOrder: 3,
	})

	if err != nil {
		t.Fatalf("AddItem error = %v", err)
	}
	if item.WorkType != "poem" || item.WorkID != 10 || item.SortOrder != 3 {
		t.Fatalf("item = %#v, want poem/10/sort 3", item)
	}
	if got := collections.collections[collection.ID].ItemCount; got != 1 {
		t.Fatalf("ItemCount = %d, want 1", got)
	}
}

func TestWorkCollectionServiceRejectsUnsupportedWorkType(t *testing.T) {
	collections := newFakeWorkCollectionStore()
	collection := collections.addCollection("小说合集")
	service := NewWorkCollectionService(collections, newFakePoemStore())

	_, err := service.AddItem(collection.ID, &dto.WorkCollectionItemCreateRequest{
		WorkType: "novel",
		WorkID:   1,
	})

	if err == nil || !strings.Contains(err.Error(), "unsupported work type") {
		t.Fatalf("AddItem error = %v, want unsupported work type", err)
	}
}

func TestWorkCollectionServiceRejectsMissingPoem(t *testing.T) {
	collections := newFakeWorkCollectionStore()
	collection := collections.addCollection("不存在的诗")
	service := NewWorkCollectionService(collections, newFakePoemStore())

	_, err := service.AddItem(collection.ID, &dto.WorkCollectionItemCreateRequest{
		WorkType: "poem",
		WorkID:   99,
	})

	if err == nil || !strings.Contains(err.Error(), "work not found") {
		t.Fatalf("AddItem error = %v, want work not found", err)
	}
}

func TestWorkCollectionServiceRejectsDuplicateItem(t *testing.T) {
	collections := newFakeWorkCollectionStore()
	poems := newFakePoemStore()
	poems.existing[7] = true
	collection := collections.addCollection("重复测试")
	service := NewWorkCollectionService(collections, poems)

	_, err := service.AddItem(collection.ID, &dto.WorkCollectionItemCreateRequest{WorkType: "poem", WorkID: 7})
	if err != nil {
		t.Fatalf("first AddItem error = %v", err)
	}

	_, err = service.AddItem(collection.ID, &dto.WorkCollectionItemCreateRequest{WorkType: "poem", WorkID: 7})

	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("second AddItem error = %v, want already exists", err)
	}
}

func TestWorkCollectionServiceUpdateItemSortRequiresCollectionOwnership(t *testing.T) {
	collections := newFakeWorkCollectionStore()
	poems := newFakePoemStore()
	poems.existing[1] = true
	first := collections.addCollection("合集一")
	second := collections.addCollection("合集二")
	service := NewWorkCollectionService(collections, poems)
	item, err := service.AddItem(first.ID, &dto.WorkCollectionItemCreateRequest{WorkType: "poem", WorkID: 1})
	if err != nil {
		t.Fatalf("AddItem error = %v", err)
	}

	_, err = service.UpdateItem(second.ID, item.ID, &dto.WorkCollectionItemUpdateRequest{SortOrder: 9})

	if err == nil || !strings.Contains(err.Error(), "item not found") {
		t.Fatalf("UpdateItem error = %v, want item not found", err)
	}
}

func TestWorkCollectionServiceRemoveItemDecrementsCount(t *testing.T) {
	collections := newFakeWorkCollectionStore()
	poems := newFakePoemStore()
	poems.existing[1] = true
	collection := collections.addCollection("删除测试")
	service := NewWorkCollectionService(collections, poems)
	item, err := service.AddItem(collection.ID, &dto.WorkCollectionItemCreateRequest{WorkType: "poem", WorkID: 1})
	if err != nil {
		t.Fatalf("AddItem error = %v", err)
	}

	err = service.DeleteItem(collection.ID, item.ID)

	if err != nil {
		t.Fatalf("DeleteItem error = %v", err)
	}
	if got := collections.collections[collection.ID].ItemCount; got != 0 {
		t.Fatalf("ItemCount = %d, want 0", got)
	}
}

type fakeWorkCollectionStore struct {
	nextCollectionID uint64
	nextItemID       uint64
	collections      map[uint64]*models.WorkCollection
	items            map[uint64]*models.WorkCollectionItem
}

func newFakeWorkCollectionStore() *fakeWorkCollectionStore {
	return &fakeWorkCollectionStore{
		nextCollectionID: 1,
		nextItemID:       1,
		collections:      map[uint64]*models.WorkCollection{},
		items:            map[uint64]*models.WorkCollectionItem{},
	}
}

func (s *fakeWorkCollectionStore) addCollection(title string) *models.WorkCollection {
	collection := &models.WorkCollection{Title: title}
	_ = s.Create(collection)
	return collection
}

func (s *fakeWorkCollectionStore) List(page, pageSize int, keyword string, published *bool) ([]models.WorkCollection, int64, error) {
	result := make([]models.WorkCollection, 0, len(s.collections))
	for _, collection := range s.collections {
		result = append(result, *collection)
	}
	return result, int64(len(result)), nil
}

func (s *fakeWorkCollectionStore) Create(collection *models.WorkCollection) error {
	collection.ID = s.nextCollectionID
	s.nextCollectionID++
	s.collections[collection.ID] = collection
	return nil
}

func (s *fakeWorkCollectionStore) GetByID(id uint64) (*models.WorkCollection, error) {
	collection, ok := s.collections[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := *collection
	for _, item := range s.items {
		if item.CollectionID == id {
			copy.Items = append(copy.Items, *item)
		}
	}
	return &copy, nil
}

func (s *fakeWorkCollectionStore) Update(collection *models.WorkCollection) error {
	s.collections[collection.ID] = collection
	return nil
}

func (s *fakeWorkCollectionStore) Delete(id uint64) error {
	delete(s.collections, id)
	return nil
}

func (s *fakeWorkCollectionStore) CreateItem(item *models.WorkCollectionItem) error {
	item.ID = s.nextItemID
	s.nextItemID++
	s.items[item.ID] = item
	return nil
}

func (s *fakeWorkCollectionStore) GetItemByID(id uint64) (*models.WorkCollectionItem, error) {
	item, ok := s.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := *item
	return &copy, nil
}

func (s *fakeWorkCollectionStore) UpdateItem(item *models.WorkCollectionItem) error {
	s.items[item.ID] = item
	return nil
}

func (s *fakeWorkCollectionStore) DeleteItem(id uint64) error {
	delete(s.items, id)
	return nil
}

func (s *fakeWorkCollectionStore) FindItemByWork(collectionID uint64, workType string, workID uint64) (*models.WorkCollectionItem, error) {
	for _, item := range s.items {
		if item.CollectionID == collectionID && item.WorkType == workType && item.WorkID == workID {
			copy := *item
			return &copy, nil
		}
	}
	return nil, errors.New("not found")
}

func (s *fakeWorkCollectionStore) RecalculateItemCount(collectionID uint64) (int, error) {
	count := 0
	for _, item := range s.items {
		if item.CollectionID == collectionID {
			count++
		}
	}
	if collection, ok := s.collections[collectionID]; ok {
		collection.ItemCount = count
	}
	return count, nil
}

type fakePoemStore struct {
	existing map[uint64]bool
}

func newFakePoemStore() *fakePoemStore {
	return &fakePoemStore{existing: map[uint64]bool{}}
}

func (s *fakePoemStore) ExistsByID(id uint64) (bool, error) {
	return s.existing[id], nil
}
