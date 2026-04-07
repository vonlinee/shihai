package models

type ForumPost struct {
	BaseModel
	UserID     uint64 `json:"userId" gorm:"not null;index;comment:用户ID"`
	User       User   `json:"user,omitempty"`
	Title      string `json:"title" gorm:"not null;size:200;comment:帖子标题"`
	Content    string `json:"content" gorm:"not null;type:text;comment:帖子内容"`
	Views      int    `json:"views" gorm:"default:0;comment:浏览量"`
	ReplyCount int    `json:"replyCount" gorm:"default:0;comment:回复数"`
	IsPinned   bool   `json:"isPinned" gorm:"default:false;comment:是否置顶"`
	IsDeleted  bool   `json:"isDeleted" gorm:"default:false;comment:是否删除"`
}

// TableName 指定表名为单数形式
func (ForumPost) TableName() string {
	return "forum_post"
}

type ForumReply struct {
	BaseModel
	PostID    uint64      `json:"postId" gorm:"not null;index;comment:帖子ID"`
	Post      ForumPost   `json:"post,omitempty"`
	UserID    uint64      `json:"userId" gorm:"not null;index;comment:用户ID"`
	User      User        `json:"user,omitempty"`
	Content   string      `json:"content" gorm:"not null;type:text;comment:回复内容"`
	ParentID  *uint64     `json:"parentId" gorm:"index;comment:父回复ID"`
	Parent    *ForumReply `json:"parent,omitempty"`
	IsDeleted bool        `json:"isDeleted" gorm:"default:false;comment:是否删除"`
}

// TableName 指定表名为单数形式
func (ForumReply) TableName() string {
	return "forum_reply"
}
