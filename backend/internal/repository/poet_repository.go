package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type PoetRepository struct {
	db *gorm.DB
}

func NewPoetRepository(db *gorm.DB) *PoetRepository {
	return &PoetRepository{db: db}
}

// Create 创建诗人
func (r *PoetRepository) Create(poet *models.Poet) error {
	return r.db.Create(poet).Error
}

// GetByID 根据ID获取诗人
func (r *PoetRepository) GetByID(id uint64) (*models.Poet, error) {
	var poet models.Poet
	err := r.db.Preload("Dynasty").First(&poet, id).Error
	if err != nil {
		return nil, err
	}
	return &poet, nil
}

// GetByName 根据名称获取诗人
func (r *PoetRepository) GetByName(name string) (*models.Poet, error) {
	var poet models.Poet
	err := r.db.Where("name = ?", name).First(&poet).Error
	if err != nil {
		return nil, err
	}
	return &poet, nil
}

// List 获取诗人列表
func (r *PoetRepository) List(keyword string) ([]models.Poet, error) {
	var poets []models.Poet
	query := r.db.Preload("Dynasty").Order("id ASC")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := query.Find(&poets).Error
	return poets, err
}

// Update 更新诗人
func (r *PoetRepository) Update(poet *models.Poet) error {
	return r.db.Save(poet).Error
}

// Delete 删除诗人
func (r *PoetRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Poet{}, id).Error
}
