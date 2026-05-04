package dto

import "time"

// CommentCreateRequest 创建评论请求，支持登录用户和游客评论
// 登录用户使用 UserID（从 Token 获取），游客使用 VisitorID + VisitorName
type CommentCreateRequest struct {
	PoemID      uint64  `json:"poemId" binding:"required"`    // 评论所属诗词ID，必填
	Content     string  `json:"content" binding:"required"`   // 评论正文内容，必填
	ParentID    *uint64 `json:"parentId"`                     // 父评论ID，回复评论时传入，顶级评论不传
	VisitorID   string  `json:"visitorId" binding:"max=64"`   // 游客唯一标识，未登录用户使用
	VisitorName string  `json:"visitorName" binding:"max=50"` // 游客显示名称
}

// CommentResponse 评论响应，包含评论信息和嵌套回复
type CommentResponse struct {
	ID          uint64            `json:"id"`                // 评论ID
	PoemID      uint64            `json:"poemId"`            // 所属诗词ID
	UserID      *uint64           `json:"userId"`            // 评论用户ID，已登录用户有值
	User        *UserResponse     `json:"user,omitempty"`    // 评论用户信息，预加载时返回
	VisitorID   string            `json:"visitorId"`         // 游客ID
	VisitorName string            `json:"visitorName"`       // 游客名称
	Content     string            `json:"content"`           // 评论正文
	ParentID    *uint64           `json:"parentId"`          // 父评论ID
	Replies     []CommentResponse `json:"replies,omitempty"` // 子回复列表
	ReplyCount  int               `json:"replyCount"`        // 回复数量
	Likes       int               `json:"likes"`             // 点赞数
	Dislikes    int               `json:"dislikes"`          // 点踩数
	IsDeleted   bool              `json:"isDeleted"`         // 是否已删除
	CreatedAt   time.Time         `json:"createdAt"`         // 创建时间
}

// CommentListRequest 评论列表请求，按诗词ID筛选并支持分页
type CommentListRequest struct {
	Page     int    `form:"page,default=1"`            // 页码，默认第1页
	PageSize int    `form:"pageSize,default=10"`       // 每页数量，默认10条
	PoemID   uint64 `form:"poemId" binding:"required"` // 诗词ID，必填，查询该诗词下的评论
}

// CommentVoteRequest 评论投票请求，支持点赞和点踩
type CommentVoteRequest struct {
	CommentID uint64 `json:"commentId" binding:"required"`               // 被投票的评论ID
	Type      string `json:"type" binding:"required,oneof=like dislike"` // 投票类型：like 或 dislike
}
