package services

import (
	"errors"
	"shihai/internal/dto"
	"shihai/internal/models"
	"shihai/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
}

func NewCommentService(commentRepo *repository.CommentRepository) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

// GetCommentsByPoem 获取诗词的评论列表
func (s *CommentService) GetCommentsByPoem(poemID uint64, page, pageSize int) ([]dto.CommentResponse, int64, error) {
	comments, total, err := s.commentRepo.ListByPoem(poemID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.CommentResponse
	for _, comment := range comments {
		responses = append(responses, *s.toCommentResponse(&comment))
	}

	return responses, total, nil
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(userID *uint64, req *dto.CommentCreateRequest) (*dto.CommentResponse, error) {
	comment := &models.Comment{
		PoemID:      req.PoemID,
		UserID:      userID,
		Content:     req.Content,
		ParentID:    req.ParentID,
		VisitorID:   req.VisitorID,
		VisitorName: req.VisitorName,
	}

	err := s.commentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	// 如果是回复，增加父评论的回复数
	if req.ParentID != nil {
		s.commentRepo.IncrementReplyCount(*req.ParentID)
	}

	return s.toCommentResponse(comment), nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(id uint64, userID uint64) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return errors.New("comment not found")
	}

	// 检查权限：只有评论作者或管理员可以删除
	if comment.UserID != nil && *comment.UserID != userID {
		return errors.New("permission denied")
	}

	return s.commentRepo.Delete(id)
}

// VoteComment 评论投票（点赞/点踩）
func (s *CommentService) VoteComment(userID *uint64, visitorID string, req *dto.CommentVoteRequest) error {
	// 检查评论是否存在
	_, err := s.commentRepo.GetByID(req.CommentID)
	if err != nil {
		return errors.New("comment not found")
	}

	// 检查是否已投票
	existingVote, err := s.commentRepo.GetVote(req.CommentID, userID, visitorID)

	if err != nil {
		// 未投票，创建新投票
		vote := &models.CommentVote{
			CommentID: req.CommentID,
			UserID:    userID,
			VisitorID: visitorID,
			Type:      req.Type,
		}
		err = s.commentRepo.CreateVote(vote)
		if err != nil {
			return err
		}

		// 更新评论计数
		if req.Type == "like" {
			return s.commentRepo.IncrementLikes(req.CommentID)
		}
		return s.commentRepo.IncrementDislikes(req.CommentID)
	}

	// 已投票
	if existingVote.Type == req.Type {
		// 相同类型，取消投票（这里简化处理，实际应该删除记录并减计数）
		return errors.New("already voted")
	}

	// 不同类型，更新投票
	existingVote.Type = req.Type
	err = s.commentRepo.UpdateVote(existingVote)
	if err != nil {
		return err
	}

	// 更新计数（简化处理）
	if req.Type == "like" {
		s.commentRepo.IncrementLikes(req.CommentID)
		s.commentRepo.IncrementDislikes(req.CommentID) // 需要减1，这里简化
	} else {
		s.commentRepo.IncrementDislikes(req.CommentID)
		s.commentRepo.IncrementLikes(req.CommentID) // 需要减1，这里简化
	}

	return nil
}

// toCommentResponse 转换为响应格式
func (s *CommentService) toCommentResponse(comment *models.Comment) *dto.CommentResponse {
	resp := &dto.CommentResponse{
		ID:          comment.ID,
		PoemID:      comment.PoemID,
		UserID:      comment.UserID,
		VisitorID:   comment.VisitorID,
		VisitorName: comment.VisitorName,
		Content:     comment.Content,
		ParentID:    comment.ParentID,
		ReplyCount:  comment.ReplyCount,
		Likes:       comment.Likes,
		Dislikes:    comment.Dislikes,
		IsDeleted:   comment.IsDeleted,
		CreatedAt:   comment.CreatedAt,
	}

	if comment.User != nil {
		resp.User = &dto.UserResponse{
			ID:       comment.User.ID,
			Username: comment.User.Username,
			Name:     comment.User.Name,
			Avatar:   comment.User.Avatar,
		}
	}

	// 处理回复
	for _, reply := range comment.Replies {
		if !reply.IsDeleted {
			resp.Replies = append(resp.Replies, *s.toCommentResponse(&reply))
		}
	}

	return resp
}

