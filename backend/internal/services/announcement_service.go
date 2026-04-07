package services

import (
	"errors"
	"shihai/internal/models"
	"shihai/internal/repository"
	"time"
)

// AnnouncementResponse 公告响应
type AnnouncementResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsPinned  bool      `json:"isPinned"`
	ViewCount int       `json:"viewCount"`
	CreatedAt time.Time `json:"createdAt"`
}

type AnnouncementService struct {
	announcementRepo *repository.AnnouncementRepository
}

func NewAnnouncementService(announcementRepo *repository.AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{announcementRepo: announcementRepo}
}

// GetAnnouncements 获取公告列表
func (s *AnnouncementService) GetAnnouncements(page, pageSize int, onlyPinned bool) ([]AnnouncementResponse, int64, error) {
	announcements, total, err := s.announcementRepo.List(page, pageSize, onlyPinned)
	if err != nil {
		return nil, 0, err
	}

	var responses []AnnouncementResponse
	for _, announcement := range announcements {
		responses = append(responses, *s.toResponse(&announcement))
	}

	return responses, total, nil
}

// GetAnnouncementByID 根据ID获取公告
func (s *AnnouncementService) GetAnnouncementByID(id uint64) (*AnnouncementResponse, error) {
	announcement, err := s.announcementRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("announcement not found")
	}

	// 增加浏览量
	s.announcementRepo.IncrementViewCount(id)

	return s.toResponse(announcement), nil
}

// CreateAnnouncement 创建公告
func (s *AnnouncementService) CreateAnnouncement(title, content string, isPinned bool) (*AnnouncementResponse, error) {
	announcement := &models.Announcement{
		Title:    title,
		Content:  content,
		IsPinned: isPinned,
	}

	err := s.announcementRepo.Create(announcement)
	if err != nil {
		return nil, err
	}

	return s.toResponse(announcement), nil
}

// UpdateAnnouncement 更新公告
func (s *AnnouncementService) UpdateAnnouncement(id uint64, title, content string, isPinned bool) (*AnnouncementResponse, error) {
	announcement, err := s.announcementRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("announcement not found")
	}

	if title != "" {
		announcement.Title = title
	}
	if content != "" {
		announcement.Content = content
	}
	announcement.IsPinned = isPinned

	err = s.announcementRepo.Update(announcement)
	if err != nil {
		return nil, err
	}

	return s.toResponse(announcement), nil
}

// DeleteAnnouncement 删除公告
func (s *AnnouncementService) DeleteAnnouncement(id uint64) error {
	return s.announcementRepo.Delete(id)
}

// toResponse 转换为响应格式
func (s *AnnouncementService) toResponse(announcement *models.Announcement) *AnnouncementResponse {
	return &AnnouncementResponse{
		ID:        announcement.ID,
		Title:     announcement.Title,
		Content:   announcement.Content,
		IsPinned:  announcement.IsPinned,
		ViewCount: announcement.ViewCount,
		CreatedAt: announcement.CreatedAt,
	}
}
