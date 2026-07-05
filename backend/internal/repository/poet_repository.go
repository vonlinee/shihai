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

// Create creates a poet extension.
func (r *PoetRepository) Create(poet *models.Poet) error {
	return r.db.Create(poet).Error
}

// GetByID returns a poet extension by ID.
func (r *PoetRepository) GetByID(id uint64) (*models.Poet, error) {
	var poet models.Poet
	err := r.db.Preload("Author").Preload("Dynasty").First(&poet, id).Error
	if err != nil {
		return nil, err
	}
	return &poet, nil
}

// GetByAuthorID returns a poet extension by author ID.
func (r *PoetRepository) GetByAuthorID(authorID uint64) (*models.Poet, error) {
	var poet models.Poet
	err := r.db.Preload("Author").Preload("Dynasty").Where("author_id = ?", authorID).First(&poet).Error
	if err != nil {
		return nil, err
	}
	return &poet, nil
}

// List returns poet extensions, optionally filtered by author name.
func (r *PoetRepository) List(keyword string) ([]models.Poet, error) {
	var poets []models.Poet
	query := r.db.Preload("Author").Preload("Dynasty").Order("poet.id ASC")
	if keyword != "" {
		query = query.Joins("JOIN author ON poet.author_id = author.id").Where("author.name LIKE ?", "%"+keyword+"%")
	}
	err := query.Find(&poets).Error
	return poets, err
}

// Update updates a poet extension.
func (r *PoetRepository) Update(poet *models.Poet) error {
	return r.db.Save(poet).Error
}

// Delete deletes a poet extension.
func (r *PoetRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Poet{}, id).Error
}
