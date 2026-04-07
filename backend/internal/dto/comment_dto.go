package dto

import "time"

// CommentCreateRequest 创建评论请求
type CommentCreateRequest struct {
	PoemID      uint64  `json:"poemId" binding:"required"`
	Content     string  `json:"content" binding:"required"`
	ParentID    *uint64 `json:"parentId"`
	VisitorID   string  `json:"visitorId" binding:"max=64"`
	VisitorName string  `json:"visitorName" binding:"max=50"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID          uint64            `json:"id"`
	PoemID      uint64            `json:"poemId"`
	UserID      *uint64           `json:"userId"`
	User        *UserResponse     `json:"user,omitempty"`
	VisitorID   string            `json:"visitorId"`
	VisitorName string            `json:"visitorName"`
	Content     string            `json:"content"`
	ParentID    *uint64           `json:"parentId"`
	Replies     []CommentResponse `json:"replies,omitempty"`
	ReplyCount  int               `json:"replyCount"`
	Likes       int               `json:"likes"`
	Dislikes    int               `json:"dislikes"`
	IsDeleted   bool              `json:"isDeleted"`
	CreatedAt   time.Time         `json:"createdAt"`
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	PoemID   uint64 `form:"poemId" binding:"required"`
}

// CommentVoteRequest 评论投票请求
type CommentVoteRequest struct {
	CommentID uint64 `json:"commentId" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=like dislike"`
}
