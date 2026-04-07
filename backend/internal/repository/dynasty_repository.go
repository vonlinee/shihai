package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type DynastyRepository struct {
	db *gorm.DB
}

func NewDynastyRepository(db *gorm.DB) *DynastyRepository {
	return &DynastyRepository{db: db}
}

// Create 创建朝代
func (r *DynastyRepository) Create(dynasty *models.Dynasty) error {
	return r.db.Create(dynasty).Error
}

// GetByID 根据ID获取朝代
func (r *DynastyRepository) GetByID(id uint64) (*models.Dynasty, error) {
	var dynasty models.Dynasty
	err := r.db.First(&dynasty, id).Error
	if err != nil {
		return nil, err
	}
	return &dynasty, nil
}

// GetByName 根据名称获取朝代
func (r *DynastyRepository) GetByName(name string) (*models.Dynasty, error) {
	var dynasty models.Dynasty
	err := r.db.Where("name = ?", name).First(&dynasty).Error
	if err != nil {
		return nil, err
	}
	return &dynasty, nil
}

// Update 更新朝代
func (r *DynastyRepository) Update(dynasty *models.Dynasty) error {
	return r.db.Save(dynasty).Error
}

// Delete 删除朝代
func (r *DynastyRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Dynasty{}, id).Error
}

// List 获取所有朝代
func (r *DynastyRepository) List() ([]models.Dynasty, error) {
	var dynasties []models.Dynasty
	err := r.db.Find(&dynasties).Error
	return dynasties, err
}
