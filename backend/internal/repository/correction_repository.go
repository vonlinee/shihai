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
