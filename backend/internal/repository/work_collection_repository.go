package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type WorkCollectionRepository struct {
	db *gorm.DB
}

func NewWorkCollectionRepository(db *gorm.DB) *WorkCollectionRepository {
	return &WorkCollectionRepository{db: db}
}

func (r *WorkCollectionRepository) List(page, pageSize int, keyword string, published *bool) ([]models.WorkCollection, int64, error) {
	var collections []models.WorkCollection
	var total int64

	query := r.db.Model(&models.WorkCollection{})
	if keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if published != nil {
		query = query.Where("is_published = ?", *published)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&collections).Error
	if err != nil {
		return nil, 0, err
	}

	return collections, total, nil
}

func (r *WorkCollectionRepository) Create(collection *models.WorkCollection) error {
	return r.db.Create(collection).Error
}

func (r *WorkCollectionRepository) GetByID(id uint64) (*models.WorkCollection, error) {
	var collection models.WorkCollection
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, created_at ASC")
	}).First(&collection, id).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *WorkCollectionRepository) Update(collection *models.WorkCollection) error {
	return r.db.Save(collection).Error
}

func (r *WorkCollectionRepository) Delete(id uint64) error {
	return r.db.Delete(&models.WorkCollection{}, id).Error
}

func (r *WorkCollectionRepository) CreateItem(item *models.WorkCollectionItem) error {
	return r.db.Create(item).Error
}

func (r *WorkCollectionRepository) GetItemByID(id uint64) (*models.WorkCollectionItem, error) {
	var item models.WorkCollectionItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *WorkCollectionRepository) UpdateItem(item *models.WorkCollectionItem) error {
	return r.db.Save(item).Error
}

func (r *WorkCollectionRepository) DeleteItem(id uint64) error {
	return r.db.Delete(&models.WorkCollectionItem{}, id).Error
}

func (r *WorkCollectionRepository) FindItemByWork(collectionID uint64, workType string, workID uint64) (*models.WorkCollectionItem, error) {
	var item models.WorkCollectionItem
	err := r.db.Where("collection_id = ? AND work_type = ? AND work_id = ?", collectionID, workType, workID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *WorkCollectionRepository) RecalculateItemCount(collectionID uint64) (int, error) {
	var count int64
	if err := r.db.Model(&models.WorkCollectionItem{}).Where("collection_id = ?", collectionID).Count(&count).Error; err != nil {
		return 0, err
	}
	if err := r.db.Model(&models.WorkCollection{}).Where("id = ?", collectionID).Update("item_count", count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
