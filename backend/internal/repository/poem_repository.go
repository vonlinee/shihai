package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

const poemTypeDisplayOrder = "CASE category WHEN '诗' THEN 1 WHEN '词' THEN 2 WHEN '曲' THEN 3 WHEN '文' THEN 4 WHEN '其他' THEN 99 ELSE 90 END ASC, id ASC"

type PoemRepository struct {
	db *gorm.DB
}

func NewPoemRepository(db *gorm.DB) *PoemRepository {
	return &PoemRepository{db: db}
}

// DB 返回底层数据库连接
func (r *PoemRepository) DB() *gorm.DB {
	return r.db
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

func (r *PoemRepository) ExistsByID(id uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Poem{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *PoemRepository) DistinctGenres() ([]string, error) {
	var genres []string
	err := r.db.Model(&models.PoemType{}).
		Order(poemTypeDisplayOrder).
		Pluck("name", &genres).Error
	return genres, err
}

// ListPoemTypes returns poem type reference data in configured display order.
func (r *PoemRepository) ListPoemTypes() ([]models.PoemType, error) {
	var poemTypes []models.PoemType
	err := r.db.Model(&models.PoemType{}).
		Order(poemTypeDisplayOrder).
		Find(&poemTypes).Error
	return poemTypes, err
}

// CreatePoemType creates poem type reference data.
func (r *PoemRepository) CreatePoemType(poemType *models.PoemType) error {
	return r.db.Create(poemType).Error
}

// GetPoemTypeByID returns poem type reference data by ID.
func (r *PoemRepository) GetPoemTypeByID(id uint64) (*models.PoemType, error) {
	var poemType models.PoemType
	if err := r.db.First(&poemType, id).Error; err != nil {
		return nil, err
	}
	return &poemType, nil
}

// UpdatePoemType updates poem type reference data.
func (r *PoemRepository) UpdatePoemType(poemType *models.PoemType) error {
	return r.db.Save(poemType).Error
}

// DeletePoemType deletes poem type reference data by ID.
func (r *PoemRepository) DeletePoemType(id uint64) error {
	return r.db.Delete(&models.PoemType{}, id).Error
}

// BatchDeletePoemTypes deletes poem type reference data by IDs.
func (r *PoemRepository) BatchDeletePoemTypes(ids []uint64) error {
	return r.db.Where("id IN ?", ids).Delete(&models.PoemType{}).Error
}

// Update 更新诗词
func (r *PoemRepository) Update(poem *models.Poem) error {
	return r.db.Save(poem).Error
}

// Delete 删除诗词
func (r *PoemRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Poem{}, id).Error
}

func (r *PoemRepository) BatchDelete(ids []uint64) error {
	return r.db.Where("id IN ?", ids).Delete(&models.Poem{}).Error
}

// List 获取诗词列表
func (r *PoemRepository) List(page, pageSize int, keyword, dynasty, author, genre string) ([]models.Poem, int64, error) {
	var poems []models.Poem
	var total int64

	query := r.db.Model(&models.Poem{}).Preload("Author").Preload("Dynasty")

	if keyword != "" {
		query = query.Where("title LIKE ? OR content::text LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if dynasty != "" {
		query = query.Joins("JOIN dynasty ON poem.dynasty_id = dynasty.id").
			Where("dynasty.name = ?", dynasty)
	}
	if author != "" {
		query = query.Joins("JOIN author ON poem.author_id = author.id").
			Where("author.name = ?", author)
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
