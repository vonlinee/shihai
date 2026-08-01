package services

import (
	"strings"

	"shihai/internal/dto"
	"shihai/internal/models"
)

type correctionStore interface {
	Create(correction *models.CorrectionRequest) error
	GetByID(id uint64) (*models.CorrectionRequest, error)
	List(page, pageSize int, keyword string) ([]models.CorrectionRequest, int64, error)
	UpdateStatus(id uint64, status string) (*models.CorrectionRequest, error)
}

// CorrectionService coordinates correction request queries.
type CorrectionService struct {
	correctionRepo correctionStore
}

// NewCorrectionService creates a CorrectionService.
func NewCorrectionService(correctionRepo correctionStore) *CorrectionService {
	return &CorrectionService{correctionRepo: correctionRepo}
}

// ListCorrections returns paginated correction requests for admin management.
func (s *CorrectionService) ListCorrections(req dto.CorrectionListRequest) ([]dto.CorrectionResponse, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	corrections, total, err := s.correctionRepo.List(page, pageSize, req.Keyword)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.CorrectionResponse, 0, len(corrections))
	for _, correction := range corrections {
		responses = append(responses, toCorrectionResponse(correction))
	}
	return responses, total, nil
}

// CreateCorrection creates a pending correction request for the current user.
func (s *CorrectionService) CreateCorrection(userID uint64, req dto.CorrectionCreateRequest) (*dto.CorrectionResponse, error) {
	correction := &models.CorrectionRequest{
		PoemID:        uint64(req.PoemID),
		UserID:        userID,
		Type:          strings.TrimSpace(req.Type),
		OriginalText:  strings.TrimSpace(req.OriginalText),
		SuggestedText: strings.TrimSpace(req.SuggestedText),
		Reason:        strings.TrimSpace(req.Reason),
		Status:        "pending",
	}

	if err := s.correctionRepo.Create(correction); err != nil {
		return nil, err
	}

	resp := toCorrectionResponse(*correction)
	return &resp, nil
}

// GetCorrection returns a correction request by ID.
func (s *CorrectionService) GetCorrection(id uint64) (*dto.CorrectionResponse, error) {
	correction, err := s.correctionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	resp := toCorrectionResponse(*correction)
	return &resp, nil
}

// UpdateCorrectionStatus updates a correction request review status.
func (s *CorrectionService) UpdateCorrectionStatus(id uint64, req dto.CorrectionStatusUpdateRequest) (*dto.CorrectionResponse, error) {
	correction, err := s.correctionRepo.UpdateStatus(id, strings.TrimSpace(req.Status))
	if err != nil {
		return nil, err
	}

	resp := toCorrectionResponse(*correction)
	return &resp, nil
}

func toCorrectionResponse(correction models.CorrectionRequest) dto.CorrectionResponse {
	resp := dto.CorrectionResponse{
		ID:            correction.ID,
		PoemID:        correction.PoemID,
		UserID:        correction.UserID,
		Type:          correction.Type,
		OriginalText:  correction.OriginalText,
		SuggestedText: correction.SuggestedText,
		Reason:        correction.Reason,
		Status:        correction.Status,
		VoteCount:     correction.VoteCount,
		ApproveCount:  correction.ApproveCount,
		RejectCount:   correction.RejectCount,
		CreatedAt:     correction.CreatedAt,
		UpdatedAt:     correction.UpdatedAt,
	}

	if correction.Poem.ID != 0 {
		resp.Poem = &dto.CorrectionPoemSummary{
			ID:    correction.Poem.ID,
			Title: correction.Poem.Title,
		}
	}
	if correction.User.ID != 0 {
		resp.User = &dto.CorrectionUserSummary{
			ID:       correction.User.ID,
			Username: correction.User.Username,
			Name:     correction.User.Name,
			Avatar:   correction.User.Avatar,
		}
	}

	return resp
}
