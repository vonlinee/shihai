package models

type Announcement struct {
	BaseModel
	Title     string `json:"title" gorm:"not null;size:200;comment:公告标题"`
	Content   string `json:"content" gorm:"not null;type:text;comment:公告内容"`
	IsPinned  bool   `json:"isPinned" gorm:"default:false;comment:是否置顶"`
	ViewCount int    `json:"viewCount" gorm:"default:0;comment:浏览次数"`
}

// TableName 指定表名为单数形式
func (Announcement) TableName() string {
	return "announcement"
}

type Feedback struct {
	BaseModel
	UserID    *uint64 `json:"userId" gorm:"index;comment:用户ID"`
	User      *User   `json:"user,omitempty"`
	VisitorID string  `json:"visitorId" gorm:"size:64;index;comment:游客ID"`
	Type      string  `json:"type" gorm:"not null;size:20;comment:反馈类型"` // bug, feature, content, other
	Title     string  `json:"title" gorm:"not null;size:200;comment:标题"`
	Content   string  `json:"content" gorm:"not null;type:text;comment:内容"`
	Contact   string  `json:"contact" gorm:"size:100;comment:联系方式"`
	Status    string  `json:"status" gorm:"default:'pending';size:20;comment:状态"` // pending, processing, resolved
}

// TableName 指定表名为单数形式
func (Feedback) TableName() string {
	return "feedback"
}

type OperationLog struct {
	BaseModel
	UserID   uint64 `json:"userId" gorm:"not null;index;comment:用户ID"`
	User     User   `json:"user,omitempty"`
	Action   string `json:"action" gorm:"not null;size:50;comment:操作类型"`
	Target   string `json:"target" gorm:"size:50;comment:操作对象"`
	TargetID uint64 `json:"targetId" gorm:"comment:对象ID"`
	Detail   string `json:"detail" gorm:"type:text;comment:详情"`
	IP       string `json:"ip" gorm:"size:50;comment:IP地址"`
}

// TableName 指定表名为单数形式
func (OperationLog) TableName() string {
	return "operation_log"
}
