package services

import (
	"testing"

	"shihai/internal/dto"
	"shihai/internal/models"
)

type fakeCorrectionStore struct {
	page     int
	pageSize int
	userID   uint64
	keyword  string
	items    []models.CorrectionRequest
	created  *models.CorrectionRequest
	item     *models.CorrectionRequest
	status   string
	total    int64
	err      error
}

func (s *fakeCorrectionStore) Create(correction *models.CorrectionRequest) error {
	s.created = correction
	return s.err
}

func (s *fakeCorrectionStore) List(page, pageSize int, keyword string) ([]models.CorrectionRequest, int64, error) {
	s.page = page
	s.pageSize = pageSize
	s.keyword = keyword
	return s.items, s.total, s.err
}

func (s *fakeCorrectionStore) ListByUser(userID uint64, page, pageSize int) ([]models.CorrectionRequest, int64, error) {
	s.userID = userID
	s.page = page
	s.pageSize = pageSize
	return s.items, s.total, s.err
}

func (s *fakeCorrectionStore) GetByID(id uint64) (*models.CorrectionRequest, error) {
	if s.item != nil {
		return s.item, s.err
	}
	for index := range s.items {
		if s.items[index].ID == id {
			return &s.items[index], s.err
		}
	}
	return nil, s.err
}

func (s *fakeCorrectionStore) UpdateStatus(id uint64, status string) (*models.CorrectionRequest, error) {
	s.status = status
	item, err := s.GetByID(id)
	if item == nil || err != nil {
		return item, err
	}
	item.Status = status
	return item, nil
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

func TestCorrectionServiceListUserCorrectionsFiltersByUserAndMapsResponse(t *testing.T) {
	store := &fakeCorrectionStore{
		items: []models.CorrectionRequest{
			{
				BaseModel:     models.BaseModel{ID: 1002},
				PoemID:        2002,
				Poem:          models.Poem{BaseModel: models.BaseModel{ID: 2002}, Title: "test poem"},
				UserID:        3002,
				Type:          "title",
				OriginalText:  "old title",
				SuggestedText: "new title",
				Reason:        "typo",
				Status:        "processing",
			},
		},
		total: 1,
	}
	service := NewCorrectionService(store)

	responses, total, err := service.ListUserCorrections(3002, dto.CorrectionListRequest{
		Page:     0,
		PageSize: 200,
	})
	if err != nil {
		t.Fatalf("ListUserCorrections returned error: %v", err)
	}
	if total != 1 {
		t.Fatalf("total = %d, want 1", total)
	}
	if store.userID != 3002 || store.page != 1 || store.pageSize != 100 {
		t.Fatalf("repository args = (%d, %d, %d), want (3002, 1, 100)", store.userID, store.page, store.pageSize)
	}
	if len(responses) != 1 || responses[0].UserID != 3002 || responses[0].Status != "processing" {
		t.Fatalf("responses = %+v, want one processing correction for user 3002", responses)
	}
}

func TestCorrectionServiceCreateCorrectionCreatesPendingRequest(t *testing.T) {
	store := &fakeCorrectionStore{}
	service := NewCorrectionService(store)

	response, err := service.CreateCorrection(3001, dto.CorrectionCreateRequest{
		PoemID:        dto.RequestID(2001),
		Type:          " content ",
		OriginalText:  " 错误原文 ",
		SuggestedText: " 正确原文 ",
		Reason:        " 错别字 ",
	})

	if err != nil {
		t.Fatalf("CreateCorrection returned error: %v", err)
	}
	if store.created == nil {
		t.Fatal("created correction = nil, want correction request")
	}
	if store.created.PoemID != 2001 || store.created.UserID != 3001 {
		t.Fatalf("created ids = poem %d user %d, want 2001/3001", store.created.PoemID, store.created.UserID)
	}
	if store.created.Type != "content" || store.created.OriginalText != "错误原文" || store.created.SuggestedText != "正确原文" {
		t.Fatalf("created correction text fields = %+v", store.created)
	}
	if store.created.Reason != "错别字" || store.created.Status != "pending" {
		t.Fatalf("created reason/status = %q/%q, want 错别字/pending", store.created.Reason, store.created.Status)
	}
	if response == nil || response.Status != "pending" || response.UserID != 3001 {
		t.Fatalf("response = %+v, want pending response for user 3001", response)
	}
}
