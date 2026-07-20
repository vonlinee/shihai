package services

import (
	"errors"
	"strings"

	"shihai/internal/dto"
	"shihai/internal/models"
)

var (
	ErrForumNotFound         = errors.New("forum resource not found")
	ErrForumPermissionDenied = errors.New("forum permission denied")
	ErrForumInvalidRequest   = errors.New("forum invalid request")
)

type forumStore interface {
	CreatePost(post *models.ForumPost) error
	GetPostByID(id uint64, includeDeleted bool) (*models.ForumPost, error)
	ListPosts(page, pageSize int, keyword string, includeDeleted bool) ([]models.ForumPost, int64, error)
	UpdatePost(post *models.ForumPost) error
	DeletePost(id uint64) error
	SetPostPinned(id uint64, isPinned bool) error
	IncrementPostViews(id uint64) error
	CreateReply(reply *models.ForumReply) error
	GetReplyByID(id uint64, includeDeleted bool) (*models.ForumReply, error)
	ListReplies(postID uint64, page, pageSize int) ([]models.ForumReply, int64, error)
	DeleteReply(id uint64, postID uint64) error
}

type ForumService struct {
	forumRepo forumStore
}

func NewForumService(forumRepo forumStore) *ForumService {
	return &ForumService{forumRepo: forumRepo}
}

func (s *ForumService) ListPosts(page, pageSize int, keyword string, includeDeleted bool) ([]dto.ForumPostResponse, int64, error) {
	page, pageSize = normalizePage(page, pageSize, 50)
	posts, total, err := s.forumRepo.ListPosts(page, pageSize, strings.TrimSpace(keyword), includeDeleted)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.ForumPostResponse, 0, len(posts))
	for _, post := range posts {
		responses = append(responses, *s.toPostResponse(&post))
	}
	return responses, total, nil
}

func (s *ForumService) GetPostByID(id uint64, incrementViews bool) (*dto.ForumPostResponse, error) {
	post, err := s.forumRepo.GetPostByID(id, false)
	if err != nil {
		return nil, ErrForumNotFound
	}

	if incrementViews {
		_ = s.forumRepo.IncrementPostViews(id)
		post.Views++
	}

	return s.toPostResponse(post), nil
}

func (s *ForumService) CreatePost(userID uint64, req *dto.ForumPostCreateRequest) (*dto.ForumPostResponse, error) {
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if title == "" || content == "" {
		return nil, ErrForumInvalidRequest
	}

	post := &models.ForumPost{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		UserID:  userID,
		Title:   title,
		Content: content,
	}
	if err := s.forumRepo.CreatePost(post); err != nil {
		return nil, err
	}

	return s.toPostResponse(post), nil
}

func (s *ForumService) UpdatePost(id uint64, userID uint64, canModerate bool, req *dto.ForumPostUpdateRequest) (*dto.ForumPostResponse, error) {
	post, err := s.forumRepo.GetPostByID(id, false)
	if err != nil {
		return nil, ErrForumNotFound
	}
	if post.UserID != userID && !canModerate {
		return nil, ErrForumPermissionDenied
	}

	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}
	post.UpdatedBy = userID

	if err := s.forumRepo.UpdatePost(post); err != nil {
		return nil, err
	}

	return s.toPostResponse(post), nil
}

func (s *ForumService) DeletePost(id uint64, userID uint64, canModerate bool) error {
	post, err := s.forumRepo.GetPostByID(id, false)
	if err != nil {
		return ErrForumNotFound
	}
	if post.UserID != userID && !canModerate {
		return ErrForumPermissionDenied
	}

	return s.forumRepo.DeletePost(id)
}

func (s *ForumService) SetPostPinned(id uint64, isPinned bool) error {
	if _, err := s.forumRepo.GetPostByID(id, false); err != nil {
		return ErrForumNotFound
	}
	return s.forumRepo.SetPostPinned(id, isPinned)
}

func (s *ForumService) ListReplies(postID uint64, page, pageSize int) ([]dto.ForumReplyResponse, int64, error) {
	if _, err := s.forumRepo.GetPostByID(postID, false); err != nil {
		return nil, 0, ErrForumNotFound
	}

	page, pageSize = normalizePage(page, pageSize, 100)
	replies, total, err := s.forumRepo.ListReplies(postID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.ForumReplyResponse, 0, len(replies))
	for _, reply := range replies {
		responses = append(responses, *s.toReplyResponse(&reply, true))
	}
	return responses, total, nil
}

func (s *ForumService) CreateReply(postID uint64, userID uint64, req *dto.ForumReplyCreateRequest) (*dto.ForumReplyResponse, error) {
	if _, err := s.forumRepo.GetPostByID(postID, false); err != nil {
		return nil, ErrForumNotFound
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, ErrForumInvalidRequest
	}

	if req.ParentID != nil {
		parent, err := s.forumRepo.GetReplyByID(uint64(*req.ParentID), false)
		if err != nil || parent.PostID != postID {
			return nil, ErrForumInvalidRequest
		}
	}

	reply := &models.ForumReply{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		PostID:   postID,
		UserID:   userID,
		Content:  content,
		ParentID: dto.RequestIDPtrValue(req.ParentID),
	}
	if err := s.forumRepo.CreateReply(reply); err != nil {
		return nil, err
	}

	created, err := s.forumRepo.GetReplyByID(reply.ID, false)
	if err != nil {
		return s.toReplyResponse(reply, false), nil
	}
	return s.toReplyResponse(created, true), nil
}

func (s *ForumService) DeleteReply(id uint64, userID uint64, canModerate bool) error {
	reply, err := s.forumRepo.GetReplyByID(id, false)
	if err != nil {
		return ErrForumNotFound
	}
	if reply.UserID != userID && !canModerate {
		return ErrForumPermissionDenied
	}

	return s.forumRepo.DeleteReply(id, reply.PostID)
}

func (s *ForumService) toPostResponse(post *models.ForumPost) *dto.ForumPostResponse {
	resp := &dto.ForumPostResponse{
		ID:         post.ID,
		UserID:     post.UserID,
		Title:      post.Title,
		Content:    post.Content,
		Views:      post.Views,
		ReplyCount: post.ReplyCount,
		IsPinned:   post.IsPinned,
		IsDeleted:  post.IsDeleted,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
	if post.User.ID != 0 {
		resp.User = toForumUserResponse(&post.User)
	}
	return resp
}

func (s *ForumService) toReplyResponse(reply *models.ForumReply, includeParent bool) *dto.ForumReplyResponse {
	resp := &dto.ForumReplyResponse{
		ID:        reply.ID,
		PostID:    reply.PostID,
		UserID:    reply.UserID,
		Content:   reply.Content,
		ParentID:  reply.ParentID,
		IsDeleted: reply.IsDeleted,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}
	if reply.User.ID != 0 {
		resp.User = toForumUserResponse(&reply.User)
	}
	if includeParent && reply.Parent != nil && !reply.Parent.IsDeleted {
		resp.Parent = s.toReplyResponse(reply.Parent, false)
	}
	return resp
}

func toForumUserResponse(user *models.User) *dto.ForumUserResponse {
	return &dto.ForumUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Avatar:   user.Avatar,
	}
}

func normalizePage(page, pageSize, maxPageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}
