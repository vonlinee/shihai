package repository

import (
	"strings"

	"shihai/internal/models"

	"gorm.io/gorm"
)

// CorrectionRepository provides persistence access for correction requests.
type CorrectionRepository struct {
	db *gorm.DB
}

// NewCorrectionRepository creates a CorrectionRepository.
func NewCorrectionRepository(db *gorm.DB) *CorrectionRepository {
	return &CorrectionRepository{db: db}
}

// Create stores a new correction request.
func (r *CorrectionRepository) Create(correction *models.CorrectionRequest) error {
	return r.db.Create(correction).Error
}

// GetByID returns a correction request with poem and user summaries.
func (r *CorrectionRepository) GetByID(id uint64) (*models.CorrectionRequest, error) {
	var correction models.CorrectionRequest
	if err := r.db.
		Preload("Poem").
		Preload("User").
		First(&correction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &correction, nil
}

// UpdateStatus updates the review status of a correction request.
func (r *CorrectionRepository) UpdateStatus(id uint64, status string) (*models.CorrectionRequest, error) {
	if err := r.db.Model(&models.CorrectionRequest{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// List returns paginated correction requests with poem and user summaries preloaded.
func (r *CorrectionRepository) List(page, pageSize int, keyword string) ([]models.CorrectionRequest, int64, error) {
	var corrections []models.CorrectionRequest
	var total int64

	baseQuery := func() *gorm.DB {
		query := r.db.Model(&models.CorrectionRequest{}).
			Preload("Poem").
			Preload("User")

		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			return query
		}

		likeKeyword := "%" + keyword + "%"
		return query.
			Joins("LEFT JOIN poem ON poem.id = correction_request.poem_id").
			Joins("LEFT JOIN \"user\" ON \"user\".id = correction_request.user_id").
			Where(
				"poem.title ILIKE ? OR \"user\".username ILIKE ? OR \"user\".name ILIKE ? OR correction_request.original_text ILIKE ? OR correction_request.suggested_text ILIKE ?",
				likeKeyword,
				likeKeyword,
				likeKeyword,
				likeKeyword,
				likeKeyword,
			)
	}

	if err := baseQuery().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery().
		Order("correction_request.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&corrections).Error; err != nil {
		return nil, 0, err
	}

	return corrections, total, nil
}
