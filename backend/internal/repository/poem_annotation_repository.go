package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type PoemAnnotationRepository struct {
	db *gorm.DB
}

func NewPoemAnnotationRepository(db *gorm.DB) *PoemAnnotationRepository {
	return &PoemAnnotationRepository{db: db}
}

func (r *PoemAnnotationRepository) ListByPoemID(poemID uint64) ([]models.PoemAnnotation, error) {
	var annotations []models.PoemAnnotation
	err := r.db.Where("poem_id = ?", poemID).
		Order("start_line ASC").
		Order("start_offset ASC").
		Order("end_line ASC").
		Order("end_offset ASC").
		Order("id ASC").
		Find(&annotations).Error
	return annotations, err
}

func (r *PoemAnnotationRepository) GetByID(id uint64) (*models.PoemAnnotation, error) {
	var annotation models.PoemAnnotation
	if err := r.db.First(&annotation, id).Error; err != nil {
		return nil, err
	}
	return &annotation, nil
}

func (r *PoemAnnotationRepository) Create(annotation *models.PoemAnnotation) error {
	return r.db.Create(annotation).Error
}

func (r *PoemAnnotationRepository) Update(annotation *models.PoemAnnotation) error {
	return r.db.Save(annotation).Error
}

func (r *PoemAnnotationRepository) Delete(id uint64) error {
	return r.db.Delete(&models.PoemAnnotation{}, id).Error
}
