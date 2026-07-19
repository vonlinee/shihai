package services

import (
	"shihai/internal/dto"
	"shihai/internal/models"
)

type correctionStore interface {
	List(page, pageSize int, keyword string) ([]models.CorrectionRequest, int64, error)
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
