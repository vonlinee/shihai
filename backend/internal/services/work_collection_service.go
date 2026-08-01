package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"shihai/internal/dto"
	"shihai/internal/models"
)

type workCollectionStore interface {
	List(page, pageSize int, keyword string, published *bool) ([]models.WorkCollection, int64, error)
	Create(collection *models.WorkCollection) error
	GetByID(id uint64) (*models.WorkCollection, error)
	Update(collection *models.WorkCollection) error
	Delete(id uint64) error
	CreateItem(item *models.WorkCollectionItem) error
	GetItemByID(id uint64) (*models.WorkCollectionItem, error)
	UpdateItem(item *models.WorkCollectionItem) error
	DeleteItem(id uint64) error
	FindItemByWork(collectionID uint64, workType string, workID uint64) (*models.WorkCollectionItem, error)
	RecalculateItemCount(collectionID uint64) (int, error)
}

type poemExistenceStore interface {
	ExistsByID(id uint64) (bool, error)
}

type WorkCollectionService struct {
	collectionRepo workCollectionStore
	poemRepo       poemExistenceStore
}

func NewWorkCollectionService(collectionRepo workCollectionStore, poemRepo poemExistenceStore) *WorkCollectionService {
	return &WorkCollectionService{
		collectionRepo: collectionRepo,
		poemRepo:       poemRepo,
	}
}

func (s *WorkCollectionService) GetCollections(req *dto.WorkCollectionListRequest) ([]dto.WorkCollectionResponse, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	collections, total, err := s.collectionRepo.List(page, pageSize, req.Keyword, req.Published)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.WorkCollectionResponse, 0, len(collections))
	for _, collection := range collections {
		responses = append(responses, *s.toCollectionResponse(&collection))
	}
	return responses, total, nil
}

func (s *WorkCollectionService) GetCollectionByID(id uint64) (*dto.WorkCollectionResponse, error) {
	collection, err := s.collectionRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("work collection not found")
	}
	return s.toCollectionResponse(collection), nil
}

func (s *WorkCollectionService) CreateCollection(req *dto.WorkCollectionCreateRequest) (*dto.WorkCollectionResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("title is required")
	}

	collection := &models.WorkCollection{
		Title:       title,
		Description: strings.TrimSpace(req.Description),
		CoverImage:  strings.TrimSpace(req.CoverImage),
		IsPublished: req.IsPublished,
	}
	if err := s.collectionRepo.Create(collection); err != nil {
		return nil, err
	}
	return s.toCollectionResponse(collection), nil
}

func (s *WorkCollectionService) UpdateCollection(id uint64, req *dto.WorkCollectionUpdateRequest) (*dto.WorkCollectionResponse, error) {
	collection, err := s.collectionRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("work collection not found")
	}

	if strings.TrimSpace(req.Title) != "" {
		collection.Title = strings.TrimSpace(req.Title)
	}
	collection.Description = strings.TrimSpace(req.Description)
	collection.CoverImage = strings.TrimSpace(req.CoverImage)
	collection.IsPublished = req.IsPublished

	if err := s.collectionRepo.Update(collection); err != nil {
		return nil, err
	}
	return s.toCollectionResponse(collection), nil
}

func (s *WorkCollectionService) DeleteCollection(id uint64) error {
	return s.collectionRepo.Delete(id)
}

func (s *WorkCollectionService) AddItem(collectionID uint64, req *dto.WorkCollectionItemCreateRequest) (*dto.WorkCollectionItemResponse, error) {
	if _, err := s.collectionRepo.GetByID(collectionID); err != nil {
		return nil, errors.New("work collection not found")
	}

	workType := strings.TrimSpace(req.WorkType)
	workID := uint64(req.WorkID)
	if err := s.validateWorkReference(workType, workID); err != nil {
		return nil, err
	}

	if existing, err := s.collectionRepo.FindItemByWork(collectionID, workType, workID); err == nil && existing != nil {
		return nil, errors.New("work collection item already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	item := &models.WorkCollectionItem{
		CollectionID: collectionID,
		WorkType:     workType,
		WorkID:       workID,
		SortOrder:    req.SortOrder,
	}
	if err := s.collectionRepo.CreateItem(item); err != nil {
		return nil, err
	}
	if _, err := s.collectionRepo.RecalculateItemCount(collectionID); err != nil {
		return nil, err
	}
	return s.toItemResponse(item), nil
}

func (s *WorkCollectionService) UpdateItem(collectionID uint64, itemID uint64, req *dto.WorkCollectionItemUpdateRequest) (*dto.WorkCollectionItemResponse, error) {
	item, err := s.collectionRepo.GetItemByID(itemID)
	if err != nil || item.CollectionID != collectionID {
		return nil, errors.New("item not found")
	}

	item.SortOrder = req.SortOrder
	if err := s.collectionRepo.UpdateItem(item); err != nil {
		return nil, err
	}
	return s.toItemResponse(item), nil
}

func (s *WorkCollectionService) DeleteItem(collectionID uint64, itemID uint64) error {
	item, err := s.collectionRepo.GetItemByID(itemID)
	if err != nil || item.CollectionID != collectionID {
		return errors.New("item not found")
	}

	if err := s.collectionRepo.DeleteItem(itemID); err != nil {
		return err
	}
	_, err = s.collectionRepo.RecalculateItemCount(collectionID)
	return err
}

func (s *WorkCollectionService) validateWorkReference(workType string, workID uint64) error {
	if workType != models.WorkTypePoem {
		return errors.New("unsupported work type")
	}
	if workID == 0 {
		return errors.New("work id is required")
	}
	exists, err := s.poemRepo.ExistsByID(workID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("work not found")
	}
	return nil
}

func (s *WorkCollectionService) toCollectionResponse(collection *models.WorkCollection) *dto.WorkCollectionResponse {
	items := make([]dto.WorkCollectionItemResponse, 0, len(collection.Items))
	for _, item := range collection.Items {
		items = append(items, *s.toItemResponse(&item))
	}

	return &dto.WorkCollectionResponse{
		ID:          collection.ID,
		Title:       collection.Title,
		Description: collection.Description,
		CoverImage:  collection.CoverImage,
		ItemCount:   collection.ItemCount,
		IsPublished: collection.IsPublished,
		Items:       items,
		CreatedAt:   collection.CreatedAt,
		UpdatedAt:   collection.UpdatedAt,
	}
}

func (s *WorkCollectionService) toItemResponse(item *models.WorkCollectionItem) *dto.WorkCollectionItemResponse {
	return &dto.WorkCollectionItemResponse{
		ID:           item.ID,
		CollectionID: item.CollectionID,
		WorkType:     item.WorkType,
		WorkID:       item.WorkID,
		SortOrder:    item.SortOrder,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}
