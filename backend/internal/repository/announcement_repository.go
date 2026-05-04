package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type AnnouncementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) *AnnouncementRepository {
	return &AnnouncementRepository{db: db}
}

// Create 创建公告
func (r *AnnouncementRepository) Create(announcement *models.Announcement) error {
	return r.db.Create(announcement).Error
}

// GetByID 根据ID获取公告
func (r *AnnouncementRepository) GetByID(id uint64) (*models.Announcement, error) {
	var announcement models.Announcement
	err := r.db.First(&announcement, id).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

// Update 更新公告
func (r *AnnouncementRepository) Update(announcement *models.Announcement) error {
	return r.db.Save(announcement).Error
}

// Delete 删除公告
func (r *AnnouncementRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Announcement{}, id).Error
}

// List 获取公告列表
func (r *AnnouncementRepository) List(page, pageSize int, onlyPinned bool) ([]models.Announcement, int64, error) {
	var announcements []models.Announcement
	var total int64

	query := r.db.Model(&models.Announcement{})

	if onlyPinned {
		query = query.Where("is_pinned = ?", true)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("is_pinned DESC, created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&announcements).Error
	if err != nil {
		return nil, 0, err
	}

	return announcements, total, nil
}

// IncrementViewCount 增加浏览量
func (r *AnnouncementRepository) IncrementViewCount(id uint64) error {
	return r.db.Model(&models.Announcement{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
