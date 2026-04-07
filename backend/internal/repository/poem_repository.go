package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type PoemRepository struct {
	db *gorm.DB
}

func NewPoemRepository(db *gorm.DB) *PoemRepository {
	return &PoemRepository{db: db}
}

// Create 创建诗词
func (r *PoemRepository) Create(poem *models.Poem) error {
	return r.db.Create(poem).Error
}

// GetByID 根据ID获取诗词
func (r *PoemRepository) GetByID(id uint64) (*models.Poem, error) {
	var poem models.Poem
	err := r.db.Preload("Author").Preload("Dynasty").First(&poem, id).Error
	if err != nil {
		return nil, err
	}
	return &poem, nil
}

// Update 更新诗词
func (r *PoemRepository) Update(poem *models.Poem) error {
	return r.db.Save(poem).Error
}

// Delete 删除诗词
func (r *PoemRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Poem{}, id).Error
}

// List 获取诗词列表
func (r *PoemRepository) List(page, pageSize int, keyword, dynasty, author, genre string) ([]models.Poem, int64, error) {
	var poems []models.Poem
	var total int64

	query := r.db.Model(&models.Poem{}).Preload("Author").Preload("Dynasty")

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if dynasty != "" {
		query = query.Joins("JOIN dynasties ON poems.dynasty_id = dynasties.id").
			Where("dynasties.name = ?", dynasty)
	}
	if author != "" {
		query = query.Joins("JOIN poets ON poems.author_id = poets.id").
			Where("poets.name = ?", author)
	}
	if genre != "" {
		query = query.Where("genre = ?", genre)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&poems).Error
	if err != nil {
		return nil, 0, err
	}

	return poems, total, nil
}

// IncrementViews 增加浏览量
func (r *PoemRepository) IncrementViews(id uint64) error {
	return r.db.Model(&models.Poem{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + 1")).Error
}

// IncrementLikes 增加点赞
func (r *PoemRepository) IncrementLikes(id uint64) error {
	return r.db.Model(&models.Poem{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + 1")).Error
}

// IncrementFavorites 增加收藏
func (r *PoemRepository) IncrementFavorites(id uint64) error {
	return r.db.Model(&models.Poem{}).Where("id = ?", id).UpdateColumn("favorites", gorm.Expr("favorites + 1")).Error
}

// GetRandom 随机获取诗词
func (r *PoemRepository) GetRandom(limit int) ([]models.Poem, error) {
	var poems []models.Poem
	err := r.db.Order("RANDOM()").Limit(limit).Preload("Author").Preload("Dynasty").Find(&poems).Error
	return poems, err
}
