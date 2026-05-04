package models

// ForumPost 论坛帖子表，存储用户发表的论坛讨论帖
// 支持置顶功能，置顶帖子会在列表中优先显示
type ForumPost struct {
	BaseModel
	UserID     uint64 `json:"userId" gorm:"not null;index;comment:用户ID"`      // 帖子作者ID，关联 User 表
	User       User   `json:"user,omitempty"`                                 // 关联的作者信息，查询时预加载
	Title      string `json:"title" gorm:"not null;size:200;comment:帖子标题"`    // 帖子标题
	Content    string `json:"content" gorm:"not null;type:text;comment:帖子内容"` // 帖子正文内容
	Views      int    `json:"views" gorm:"default:0;comment:浏览量"`             // 浏览次数
	ReplyCount int    `json:"replyCount" gorm:"default:0;comment:回复数"`        // 回复数量统计
	IsPinned   bool   `json:"isPinned" gorm:"default:false;comment:是否置顶"`     // 是否置顶，置顶帖子优先显示
	IsDeleted  bool   `json:"isDeleted" gorm:"default:false;comment:是否删除"`    // 软删除标记，true 表示帖子已被删除
}

// TableName 指定表名为单数形式
func (ForumPost) TableName() string {
	return "forum_post"
}

// ForumReply 论坛回复表，存储用户对帖子的回复
// 支持嵌套回复（通过 ParentID 关联父回复），实现楼中楼功能
type ForumReply struct {
	BaseModel
	PostID    uint64      `json:"postId" gorm:"not null;index;comment:帖子ID"`      // 所属帖子ID，关联 ForumPost 表
	Post      ForumPost   `json:"post,omitempty"`                                 // 关联的帖子信息，查询时预加载
	UserID    uint64      `json:"userId" gorm:"not null;index;comment:用户ID"`      // 回复者ID，关联 User 表
	User      User        `json:"user,omitempty"`                                 // 关联的回复者信息，查询时预加载
	Content   string      `json:"content" gorm:"not null;type:text;comment:回复内容"` // 回复正文内容
	ParentID  *uint64     `json:"parentId" gorm:"index;comment:父回复ID"`            // 父回复ID，顶级回复为 nil
	Parent    *ForumReply `json:"parent,omitempty"`                               // 关联的父回复信息
	IsDeleted bool        `json:"isDeleted" gorm:"default:false;comment:是否删除"`    // 软删除标记，true 表示回复已被删除
}

// TableName 指定表名为单数形式
func (ForumReply) TableName() string {
	return "forum_reply"
}
