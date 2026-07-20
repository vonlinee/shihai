package dto

import "time"

// ForumPostCreateRequest 创建论坛帖子请求。
type ForumPostCreateRequest struct {
	Title   string `json:"title" binding:"required,min=2,max=200"`     // Title 帖子标题，2-200 个字符。
	Content string `json:"content" binding:"required,min=2,max=10000"` // Content 帖子正文，最多 10000 个字符。
}

// ForumPostUpdateRequest 更新论坛帖子请求。
type ForumPostUpdateRequest struct {
	Title   string `json:"title" binding:"omitempty,min=2,max=200"`     // Title 帖子标题；为空时不更新。
	Content string `json:"content" binding:"omitempty,min=2,max=10000"` // Content 帖子正文；为空时不更新。
}

// ForumReplyCreateRequest 创建论坛回复请求。
type ForumReplyCreateRequest struct {
	Content  string     `json:"content" binding:"required,min=1,max=5000"` // Content 回复正文，最多 5000 个字符。
	ParentID *RequestID `json:"parentId"`                                  // ParentID 父回复 Snowflake ID，楼中楼回复时传入。
}

// ForumPostPinRequest 设置帖子置顶状态请求。
type ForumPostPinRequest struct {
	IsPinned bool `json:"isPinned"` // IsPinned 是否置顶，置顶帖列表优先展示。
}

// ForumUserResponse 论坛用户摘要。
type ForumUserResponse struct {
	ID       uint64 `json:"id,string"` // ID 用户 Snowflake 主键，对外序列化为字符串。
	Username string `json:"username"`  // Username 用户名。
	Name     string `json:"name"`      // Name 用户显示名称。
	Avatar   string `json:"avatar"`    // Avatar 用户头像 URL。
}

// ForumPostResponse 论坛帖子响应。
type ForumPostResponse struct {
	ID         uint64             `json:"id,string"`      // ID 帖子 Snowflake 主键，对外序列化为字符串。
	UserID     uint64             `json:"userId,string"`  // UserID 发帖用户 Snowflake ID。
	User       *ForumUserResponse `json:"user,omitempty"` // User 发帖用户摘要，预加载时返回。
	Title      string             `json:"title"`          // Title 帖子标题。
	Content    string             `json:"content"`        // Content 帖子正文。
	Views      int                `json:"views"`          // Views 浏览量。
	ReplyCount int                `json:"replyCount"`     // ReplyCount 回复数量。
	IsPinned   bool               `json:"isPinned"`       // IsPinned 是否置顶。
	IsDeleted  bool               `json:"isDeleted"`      // IsDeleted 是否已被业务删除。
	CreatedAt  time.Time          `json:"createdAt"`      // CreatedAt 创建时间。
	UpdatedAt  time.Time          `json:"updatedAt"`      // UpdatedAt 更新时间。
}

// ForumReplyResponse 论坛回复响应。
type ForumReplyResponse struct {
	ID        uint64              `json:"id,string"`        // ID 回复 Snowflake 主键，对外序列化为字符串。
	PostID    uint64              `json:"postId,string"`    // PostID 所属帖子 Snowflake ID。
	UserID    uint64              `json:"userId,string"`    // UserID 回复用户 Snowflake ID。
	User      *ForumUserResponse  `json:"user,omitempty"`   // User 回复用户摘要，预加载时返回。
	Content   string              `json:"content"`          // Content 回复正文。
	ParentID  *uint64             `json:"parentId,string"`  // ParentID 父回复 Snowflake ID，顶级回复为空。
	Parent    *ForumReplyResponse `json:"parent,omitempty"` // Parent 父回复摘要，楼中楼回复时返回。
	IsDeleted bool                `json:"isDeleted"`        // IsDeleted 是否已被业务删除。
	CreatedAt time.Time           `json:"createdAt"`        // CreatedAt 创建时间。
	UpdatedAt time.Time           `json:"updatedAt"`        // UpdatedAt 更新时间。
}
