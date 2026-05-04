package models

// Comment 评论表，存储用户和游客对诗词的评论
// 支持嵌套回复（通过 ParentID 关联父评论）
// 用户和游客二选一标识：已登录用户用 UserID，未登录游客用 VisitorID
type Comment struct {
	BaseModel
	PoemID      uint64    `json:"poemId" gorm:"not null;index;comment:诗词ID"`      // 评论所属诗词ID，关联 Poem 表
	Poem        Poem      `json:"poem,omitempty"`                                 // 关联的诗词信息，查询时预加载
	UserID      *uint64   `json:"userId" gorm:"index;comment:用户ID"`               // 评论用户ID，已登录用户使用，关联 User 表
	User        *User     `json:"user,omitempty"`                                 // 关联的用户信息，查询时预加载
	VisitorID   string    `json:"visitorId" gorm:"size:64;index;comment:游客ID"`    // 游客唯一标识，未登录用户使用
	VisitorName string    `json:"visitorName" gorm:"size:50;comment:游客名称"`        // 游客显示名称
	Content     string    `json:"content" gorm:"not null;type:text;comment:评论内容"` // 评论正文内容
	ParentID    *uint64   `json:"parentId" gorm:"index;comment:父评论ID"`            // 父评论ID，顶级评论为 nil
	Parent      *Comment  `json:"parent,omitempty"`                               // 关联的父评论信息
	Replies     []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`   // 子回复列表，通过 ParentID 外键关联
	ReplyCount  int       `json:"replyCount" gorm:"default:0;comment:回复数"`        // 回复数量统计
	Likes       int       `json:"likes" gorm:"default:0;comment:点赞数"`             // 点赞数量
	Dislikes    int       `json:"dislikes" gorm:"default:0;comment:点踩数"`          // 点踩数量
	IsDeleted   bool      `json:"isDeleted" gorm:"default:false;comment:是否删除"`    // 软删除标记，true 表示评论已被删除
}

// TableName 指定表名为单数形式
func (Comment) TableName() string {
	return "comment"
}

// CommentVote 评论投票表，记录用户/游客对评论的投票行为
// 每个用户/游客对同一评论只能投一次票（like 或 dislike）
type CommentVote struct {
	BaseModel
	CommentID uint64  `json:"commentId" gorm:"not null;index;comment:评论ID"` // 被投票的评论ID，关联 Comment 表
	UserID    *uint64 `json:"userId" gorm:"index;comment:用户ID"`             // 投票用户ID，已登录用户使用
	VisitorID string  `json:"visitorId" gorm:"size:64;index;comment:游客ID"`  // 投票游客ID，未登录用户使用
	Type      string  `json:"type" gorm:"not null;size:10;comment:投票类型"`    // 投票类型：like(点赞) 或 dislike(点踩)
}

// TableName 指定表名为单数形式
func (CommentVote) TableName() string {
	return "comment_vote"
}
