package dto

import "time"

// CorrectionListRequest 纠错列表请求，支持分页和关键词搜索。
type CorrectionListRequest struct {
	Page     int    `form:"page,default=1"`      // Page 页码，从 1 开始。
	PageSize int    `form:"pageSize,default=10"` // PageSize 每页数量，最大 100。
	Keyword  string `form:"keyword"`             // Keyword 搜索关键词，匹配诗词标题、提交用户、原文或建议内容。
}

// CorrectionCreateRequest 提交诗词纠错申请。
type CorrectionCreateRequest struct {
	PoemID        RequestID `json:"poemId" binding:"required"`        // PoemID 被纠错诗词的 Snowflake ID。
	Type          string    `json:"type" binding:"required,max=20"`   // Type 纠错位置或类型，例如 title、author、dynasty、content、translation、appreciation、annotation 或 other。
	OriginalText  string    `json:"originalText" binding:"required"`  // OriginalText 当前错误或待确认的原文。
	SuggestedText string    `json:"suggestedText" binding:"required"` // SuggestedText 用户建议修改后的内容。
	Reason        string    `json:"reason" binding:"required"`        // Reason 纠错原因、依据和补充说明。
}

// CorrectionStatusUpdateRequest 更新纠错申请状态。
type CorrectionStatusUpdateRequest struct {
	Status string `json:"status" binding:"required,oneof=pending voting processing approved rejected resolved completed"` // Status 目标状态。
}

// CorrectionPoemSummary 纠错列表中的诗词摘要。
type CorrectionPoemSummary struct {
	ID    uint64 `json:"id,string"` // ID 诗词 Snowflake ID，对外序列化为字符串。
	Title string `json:"title"`     // Title 诗词标题。
}

// CorrectionUserSummary 纠错列表中的提交用户摘要。
type CorrectionUserSummary struct {
	ID       uint64 `json:"id,string"` // ID 用户 Snowflake ID，对外序列化为字符串。
	Username string `json:"username"`  // Username 用户名。
	Name     string `json:"name"`      // Name 用户显示名称。
	Avatar   string `json:"avatar"`    // Avatar 用户头像地址。
}

// CorrectionResponse 纠错申请响应。
type CorrectionResponse struct {
	ID            uint64                 `json:"id,string"`      // ID 纠错申请 Snowflake ID，对外序列化为字符串。
	PoemID        uint64                 `json:"poemId,string"`  // PoemID 被纠错诗词的 Snowflake ID。
	Poem          *CorrectionPoemSummary `json:"poem,omitempty"` // Poem 被纠错诗词摘要。
	UserID        uint64                 `json:"userId,string"`  // UserID 提交用户的 Snowflake ID。
	User          *CorrectionUserSummary `json:"user,omitempty"` // User 提交用户摘要。
	Type          string                 `json:"type"`           // Type 纠错类型。
	OriginalText  string                 `json:"originalText"`   // OriginalText 被纠错的原文内容。
	SuggestedText string                 `json:"suggestedText"`  // SuggestedText 用户建议修改后的内容。
	Reason        string                 `json:"reason"`         // Reason 用户提交的纠错理由。
	Status        string                 `json:"status"`         // Status 当前状态：pending、voting、processing、approved、rejected、resolved 或 completed。
	VoteCount     int                    `json:"voteCount"`      // VoteCount 投票总数。
	ApproveCount  int                    `json:"approveCount"`   // ApproveCount 支持票数。
	RejectCount   int                    `json:"rejectCount"`    // RejectCount 反对票数。
	CreatedAt     time.Time              `json:"createdAt"`      // CreatedAt 创建时间。
	UpdatedAt     time.Time              `json:"updatedAt"`      // UpdatedAt 更新时间。
}
