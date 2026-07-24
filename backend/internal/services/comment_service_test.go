package services

import (
	"errors"
	"testing"

	"shihai/internal/dto"
	"shihai/internal/models"
)

type fakeCommentStore struct {
	listAllPage     int
	listAllPageSize int
	comments        map[uint64]*models.Comment
	createdComment  *models.Comment
	incrementedID   uint64
}

func (s *fakeCommentStore) Create(comment *models.Comment) error {
	if s.comments == nil {
		s.comments = make(map[uint64]*models.Comment)
	}
	if comment.ID == 0 {
		comment.ID = 100
	}
	s.createdComment = comment
	s.comments[comment.ID] = comment
	return nil
}

func (s *fakeCommentStore) GetByID(id uint64) (*models.Comment, error) {
	comment, ok := s.comments[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return comment, nil
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
	s.incrementedID = id
	if comment, ok := s.comments[id]; ok {
		comment.ReplyCount++
	}
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

func TestCommentServiceCreateReplyRequiresParentInSamePoem(t *testing.T) {
	parentID := dto.RequestID(10)
	store := &fakeCommentStore{
		comments: map[uint64]*models.Comment{
			10: {PoemID: 2, Content: "parent"},
		},
	}
	service := NewCommentService(store)

	_, err := service.CreateComment(nil, &dto.CommentCreateRequest{
		PoemID:   dto.RequestID(1),
		Content:  "reply",
		ParentID: &parentID,
	})

	if err == nil {
		t.Fatalf("expected invalid parent comment error")
	}
	if store.createdComment != nil {
		t.Fatalf("reply should not be created with parent from another poem")
	}
}

func TestCommentServiceCreateReplyIncrementsParentReplyCount(t *testing.T) {
	parentID := dto.RequestID(10)
	store := &fakeCommentStore{
		comments: map[uint64]*models.Comment{
			10: {PoemID: 1, Content: "parent"},
		},
	}
	service := NewCommentService(store)

	reply, err := service.CreateComment(nil, &dto.CommentCreateRequest{
		PoemID:   dto.RequestID(1),
		Content:  "reply",
		ParentID: &parentID,
	})

	if err != nil {
		t.Fatalf("CreateComment error = %v", err)
	}
	if reply.ParentID == nil || *reply.ParentID != 10 {
		t.Fatalf("expected reply parent id 10, got %#v", reply.ParentID)
	}
	if store.incrementedID != 10 || store.comments[10].ReplyCount != 1 {
		t.Fatalf("expected parent reply count incremented, got id %d count %d", store.incrementedID, store.comments[10].ReplyCount)
	}
}

func TestCommentServiceCreateReplyAllowsReplyToReply(t *testing.T) {
	replyID := dto.RequestID(20)
	store := &fakeCommentStore{
		comments: map[uint64]*models.Comment{
			10: {PoemID: 1, Content: "parent"},
			20: {PoemID: 1, ParentID: uint64Ptr(10), Content: "reply"},
		},
	}
	service := NewCommentService(store)

	nestedReply, err := service.CreateComment(nil, &dto.CommentCreateRequest{
		PoemID:   dto.RequestID(1),
		Content:  "nested reply",
		ParentID: &replyID,
	})

	if err != nil {
		t.Fatalf("CreateComment error = %v", err)
	}
	if nestedReply.ParentID == nil || *nestedReply.ParentID != 20 {
		t.Fatalf("expected nested reply parent id 20, got %#v", nestedReply.ParentID)
	}
	if store.incrementedID != 20 || store.comments[20].ReplyCount != 1 {
		t.Fatalf("expected reply count incremented on reply 20, got id %d count %d", store.incrementedID, store.comments[20].ReplyCount)
	}
}

func uint64Ptr(value uint64) *uint64 {
	return &value
}
