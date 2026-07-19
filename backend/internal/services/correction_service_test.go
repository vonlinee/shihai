package services

import (
	"testing"

	"shihai/internal/dto"
	"shihai/internal/models"
)

type fakeCorrectionStore struct {
	page     int
	pageSize int
	keyword  string
	items    []models.CorrectionRequest
	total    int64
	err      error
}

func (s *fakeCorrectionStore) List(page, pageSize int, keyword string) ([]models.CorrectionRequest, int64, error) {
	s.page = page
	s.pageSize = pageSize
	s.keyword = keyword
	return s.items, s.total, s.err
}

func TestCorrectionServiceListCorrectionsNormalizesPaginationAndMapsResponse(t *testing.T) {
	store := &fakeCorrectionStore{
		items: []models.CorrectionRequest{
			{
				BaseModel:     models.BaseModel{ID: 1001},
				PoemID:        2001,
				Poem:          models.Poem{BaseModel: models.BaseModel{ID: 2001}, Title: "静夜思"},
				UserID:        3001,
				User:          models.User{BaseModel: models.BaseModel{ID: 3001}, Username: "reader", Name: "读者"},
				Type:          "content",
				OriginalText:  "床前明月光",
				SuggestedText: "窗前明月光",
				Reason:        "原文疑似有误",
				Status:        "pending",
			},
		},
		total: 1,
	}
	service := NewCorrectionService(store)

	responses, total, err := service.ListCorrections(dto.CorrectionListRequest{
		Page:     0,
		PageSize: 200,
		Keyword:  "明月",
	})
	if err != nil {
		t.Fatalf("ListCorrections returned error: %v", err)
	}
	if total != 1 {
		t.Fatalf("total = %d, want 1", total)
	}
	if store.page != 1 || store.pageSize != 100 || store.keyword != "明月" {
		t.Fatalf("repository args = (%d, %d, %q), want (1, 100, 明月)", store.page, store.pageSize, store.keyword)
	}
	if len(responses) != 1 {
		t.Fatalf("len(responses) = %d, want 1", len(responses))
	}

	got := responses[0]
	if got.ID != 1001 || got.PoemID != 2001 || got.UserID != 3001 {
		t.Fatalf("unexpected ids: %+v", got)
	}
	if got.Poem == nil || got.Poem.Title != "静夜思" {
		t.Fatalf("poem summary = %+v, want title 静夜思", got.Poem)
	}
	if got.User == nil || got.User.Name != "读者" || got.User.Username != "reader" {
		t.Fatalf("user summary = %+v, want reader/读者", got.User)
	}
	if got.OriginalText != "床前明月光" || got.SuggestedText != "窗前明月光" {
		t.Fatalf("unexpected text fields: %+v", got)
	}
}
