package services

import (
	"errors"
	"testing"

	"shihai/internal/models"
)

type fakeCommentStore struct {
	listAllPage     int
	listAllPageSize int
}

func (s *fakeCommentStore) Create(comment *models.Comment) error {
	return nil
}

func (s *fakeCommentStore) GetByID(id uint64) (*models.Comment, error) {
	return nil, errors.New("not implemented")
}

func (s *fakeCommentStore) ListByPoem(poemID uint64, page, pageSize int) ([]models.Comment, int64, error) {
	return nil, 0, errors.New("not implemented")
}

func (s *fakeCommentStore) ListAll(page, pageSize int) ([]models.Comment, int64, error) {
	s.listAllPage = page
	s.listAllPageSize = pageSize
	return []models.Comment{{PoemID: 34, Content: "test comment"}}, 1, nil
}

func (s *fakeCommentStore) Update(comment *models.Comment) error {
	return nil
}

func (s *fakeCommentStore) Delete(id uint64) error {
	return nil
}

func (s *fakeCommentStore) IncrementLikes(id uint64) error {
	return nil
}

func (s *fakeCommentStore) IncrementDislikes(id uint64) error {
	return nil
}

func (s *fakeCommentStore) IncrementReplyCount(id uint64) error {
	return nil
}

func (s *fakeCommentStore) GetVote(commentID uint64, userID *uint64, visitorID string) (*models.CommentVote, error) {
	return nil, errors.New("not implemented")
}

func (s *fakeCommentStore) CreateVote(vote *models.CommentVote) error {
	return nil
}

func (s *fakeCommentStore) UpdateVote(vote *models.CommentVote) error {
	return nil
}

func TestCommentServiceGetAllCommentsUsesAllCommentsStore(t *testing.T) {
	store := &fakeCommentStore{}
	service := NewCommentService(store)

	comments, total, err := service.GetAllComments(2, 20)

	if err != nil {
		t.Fatalf("GetAllComments error = %v", err)
	}
	if store.listAllPage != 2 || store.listAllPageSize != 20 {
		t.Fatalf("expected ListAll page/pageSize 2/20, got %d/%d", store.listAllPage, store.listAllPageSize)
	}
	if total != 1 || len(comments) != 1 || comments[0].PoemID != 34 {
		t.Fatalf("unexpected comments/total: %#v / %d", comments, total)
	}
}
