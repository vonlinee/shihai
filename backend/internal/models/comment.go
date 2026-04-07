package models

type Comment struct {
	BaseModel
	PoemID      uint64    `json:"poemId" gorm:"not null;index;comment:诗词ID"`
	Poem        Poem      `json:"poem,omitempty"`
	UserID      *uint64   `json:"userId" gorm:"index;comment:用户ID"`
	User        *User     `json:"user,omitempty"`
	VisitorID   string    `json:"visitorId" gorm:"size:64;index;comment:游客ID"`
	VisitorName string    `json:"visitorName" gorm:"size:50;comment:游客名称"`
	Content     string    `json:"content" gorm:"not null;type:text;comment:评论内容"`
	ParentID    *uint64   `json:"parentId" gorm:"index;comment:父评论ID"`
	Parent      *Comment  `json:"parent,omitempty"`
	Replies     []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	ReplyCount  int       `json:"replyCount" gorm:"default:0;comment:回复数"`
	Likes       int       `json:"likes" gorm:"default:0;comment:点赞数"`
	Dislikes    int       `json:"dislikes" gorm:"default:0;comment:点踩数"`
	IsDeleted   bool      `json:"isDeleted" gorm:"default:false;comment:是否删除"`
}

// TableName 指定表名为单数形式
func (Comment) TableName() string {
	return "comment"
}

type CommentVote struct {
	BaseModel
	CommentID uint64  `json:"commentId" gorm:"not null;index;comment:评论ID"`
	UserID    *uint64 `json:"userId" gorm:"index;comment:用户ID"`
	VisitorID string  `json:"visitorId" gorm:"size:64;index;comment:游客ID"`
	Type      string  `json:"type" gorm:"not null;size:10;comment:投票类型"` // like, dislike
}

// TableName 指定表名为单数形式
func (CommentVote) TableName() string {
	return "comment_vote"
}
