package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

// Create creates an author.
func (r *AuthorRepository) Create(author *models.Author) error {
	return r.db.Create(author).Error
}

// GetByID returns an author by ID.
func (r *AuthorRepository) GetByID(id uint64) (*models.Author, error) {
	var author models.Author
	err := r.db.First(&author, id).Error
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// GetByName returns an author by name.
func (r *AuthorRepository) GetByName(name string) (*models.Author, error) {
	var author models.Author
	err := r.db.Where("name = ?", name).First(&author).Error
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// List returns authors ordered by ID.
func (r *AuthorRepository) List(keyword string) ([]models.Author, error) {
	var authors []models.Author
	query := r.db.Order("id ASC")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := query.Find(&authors).Error
	return authors, err
}

// Update updates an author.
func (r *AuthorRepository) Update(author *models.Author) error {
	return r.db.Save(author).Error
}

// Delete deletes an author.
func (r *AuthorRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Author{}, id).Error
}
