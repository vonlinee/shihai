package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"shihai/internal/dto"

	"github.com/gin-gonic/gin"
)

type fakeCommentService struct {
	allPage     int
	allPageSize int
}

func (s *fakeCommentService) GetCommentsByPoem(poemID uint64, page, pageSize int) ([]dto.CommentResponse, int64, error) {
	return nil, 0, errors.New("not implemented")
}

func (s *fakeCommentService) GetAllComments(page, pageSize int) ([]dto.CommentResponse, int64, error) {
	s.allPage = page
	s.allPageSize = pageSize
	return []dto.CommentResponse{{ID: 12, PoemID: 34, Content: "test comment"}}, 1, nil
}

func (s *fakeCommentService) CreateComment(userID *uint64, req *dto.CommentCreateRequest) (*dto.CommentResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *fakeCommentService) DeleteComment(id uint64, userID uint64) error {
	return errors.New("not implemented")
}

func (s *fakeCommentService) VoteComment(userID *uint64, visitorID string, req *dto.CommentVoteRequest) error {
	return errors.New("not implemented")
}

func TestCommentHandlerGetAllCommentsDoesNotRequirePoemID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeCommentService{}
	handler := NewCommentHandler(service)
	router := gin.New()
	router.GET("/api/admin/comments/all", handler.GetAllComments)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/comments/all?page=2&pageSize=20", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.allPage != 2 || service.allPageSize != 20 {
		t.Fatalf("expected page/pageSize 2/20, got %d/%d", service.allPage, service.allPageSize)
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			List []dto.CommentResponse `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Code != 200 || len(response.Data.List) != 1 || response.Data.List[0].ID != 12 {
		t.Fatalf("unexpected response: %#v", response)
	}
}
