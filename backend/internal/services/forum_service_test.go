package services

import (
	"errors"
	"testing"

	"shihai/internal/dto"
	"shihai/internal/models"
)

type fakeForumStore struct {
	posts          map[uint64]*models.ForumPost
	replies        map[uint64]*models.ForumReply
	createdReply   *models.ForumReply
	deletedPostID  uint64
	deletedReplyID uint64
}

func newFakeForumStore() *fakeForumStore {
	return &fakeForumStore{
		posts:   make(map[uint64]*models.ForumPost),
		replies: make(map[uint64]*models.ForumReply),
	}
}

func (s *fakeForumStore) CreatePost(post *models.ForumPost) error {
	if post.ID == 0 {
		post.ID = 100
	}
	s.posts[post.ID] = post
	return nil
}

func (s *fakeForumStore) GetPostByID(id uint64, includeDeleted bool) (*models.ForumPost, error) {
	post, ok := s.posts[id]
	if !ok || (!includeDeleted && post.IsDeleted) {
		return nil, errors.New("not found")
	}
	return post, nil
}

func (s *fakeForumStore) ListPosts(page, pageSize int, keyword string, includeDeleted bool) ([]models.ForumPost, int64, error) {
	return nil, 0, errors.New("not implemented")
}

func (s *fakeForumStore) UpdatePost(post *models.ForumPost) error {
	s.posts[post.ID] = post
	return nil
}

func (s *fakeForumStore) DeletePost(id uint64) error {
	s.deletedPostID = id
	if post, ok := s.posts[id]; ok {
		post.IsDeleted = true
	}
	return nil
}

func (s *fakeForumStore) SetPostPinned(id uint64, isPinned bool) error {
	if post, ok := s.posts[id]; ok {
		post.IsPinned = isPinned
	}
	return nil
}

func (s *fakeForumStore) IncrementPostViews(id uint64) error {
	if post, ok := s.posts[id]; ok {
		post.Views++
	}
	return nil
}

func (s *fakeForumStore) CreateReply(reply *models.ForumReply) error {
	if reply.ID == 0 {
		reply.ID = 200
	}
	s.createdReply = reply
	s.replies[reply.ID] = reply
	if post, ok := s.posts[reply.PostID]; ok {
		post.ReplyCount++
	}
	return nil
}

func (s *fakeForumStore) GetReplyByID(id uint64, includeDeleted bool) (*models.ForumReply, error) {
	reply, ok := s.replies[id]
	if !ok || (!includeDeleted && reply.IsDeleted) {
		return nil, errors.New("not found")
	}
	return reply, nil
}

func (s *fakeForumStore) ListReplies(postID uint64, page, pageSize int) ([]models.ForumReply, int64, error) {
	return nil, 0, errors.New("not implemented")
}

func (s *fakeForumStore) DeleteReply(id uint64, postID uint64) error {
	s.deletedReplyID = id
	if reply, ok := s.replies[id]; ok {
		reply.IsDeleted = true
	}
	if post, ok := s.posts[postID]; ok && post.ReplyCount > 0 {
		post.ReplyCount--
	}
	return nil
}

func TestForumServiceDeletePostRequiresOwnerOrModerator(t *testing.T) {
	store := newFakeForumStore()
	store.posts[1] = &models.ForumPost{UserID: 10, Title: "title", Content: "content"}
	service := NewForumService(store)

	err := service.DeletePost(1, 20, false)

	if !errors.Is(err, ErrForumPermissionDenied) {
		t.Fatalf("expected permission denied, got %v", err)
	}
	if store.deletedPostID != 0 {
		t.Fatalf("post should not be deleted by non-owner, got deleted id %d", store.deletedPostID)
	}

	if err := service.DeletePost(1, 20, true); err != nil {
		t.Fatalf("moderator delete post error = %v", err)
	}
	if store.deletedPostID != 1 {
		t.Fatalf("expected deleted post id 1, got %d", store.deletedPostID)
	}
}

func TestForumServiceCreateReplyRejectsParentFromAnotherPost(t *testing.T) {
	parentID := dto.RequestID(20)
	store := newFakeForumStore()
	store.posts[1] = &models.ForumPost{UserID: 10, Title: "title", Content: "content"}
	store.replies[20] = &models.ForumReply{PostID: 2, UserID: 11, Content: "parent"}
	service := NewForumService(store)

	_, err := service.CreateReply(1, 10, &dto.ForumReplyCreateRequest{
		Content:  "reply",
		ParentID: &parentID,
	})

	if !errors.Is(err, ErrForumInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}
	if store.createdReply != nil {
		t.Fatalf("reply should not be created with parent from another post")
	}
}

func TestForumServiceCreateReplyTrimsContentAndCreatesReply(t *testing.T) {
	store := newFakeForumStore()
	store.posts[1] = &models.ForumPost{UserID: 10, Title: "title", Content: "content"}
	service := NewForumService(store)

	reply, err := service.CreateReply(1, 10, &dto.ForumReplyCreateRequest{Content: "  好帖  "})

	if err != nil {
		t.Fatalf("CreateReply error = %v", err)
	}
	if reply.Content != "好帖" {
		t.Fatalf("expected trimmed content, got %q", reply.Content)
	}
	if store.posts[1].ReplyCount != 1 {
		t.Fatalf("expected reply count 1, got %d", store.posts[1].ReplyCount)
	}
	if store.createdReply == nil || store.createdReply.CreatedBy != 10 || store.createdReply.UpdatedBy != 10 {
		t.Fatalf("expected audit fields set on created reply, got %#v", store.createdReply)
	}
}
