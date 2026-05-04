package models

// Announcement 公告表，存储系统公告和通知
// 支持置顶功能，置顶公告会在列表中优先显示
type Announcement struct {
	BaseModel
	Title     string `json:"title" gorm:"not null;size:200;comment:公告标题"`    // 公告标题
	Content   string `json:"content" gorm:"not null;type:text;comment:公告内容"` // 公告正文内容，支持富文本
	IsPinned  bool   `json:"isPinned" gorm:"default:false;comment:是否置顶"`     // 是否置顶，置顶公告优先显示
	ViewCount int    `json:"viewCount" gorm:"default:0;comment:浏览次数"`        // 浏览次数统计
}

// TableName 指定表名为单数形式
func (Announcement) TableName() string {
	return "announcement"
}

// Feedback 用户反馈表，存储用户和游客的意见反馈
// 反馈类型包括：bug(Bug报告)、feature(功能建议)、content(内容反馈)、other(其他)
// 状态流转：pending(待处理) → processing(处理中) → resolved(已解决)
type Feedback struct {
	BaseModel
	UserID    *uint64 `json:"userId" gorm:"index;comment:用户ID"`                   // 反馈用户ID，已登录用户使用，关联 User 表
	User      *User   `json:"user,omitempty"`                                     // 关联的用户信息，查询时预加载
	VisitorID string  `json:"visitorId" gorm:"size:64;index;comment:游客ID"`        // 反馈游客ID，未登录用户使用
	Type      string  `json:"type" gorm:"not null;size:20;comment:反馈类型"`          // 反馈类型：bug/feature/content/other
	Title     string  `json:"title" gorm:"not null;size:200;comment:标题"`          // 反馈标题
	Content   string  `json:"content" gorm:"not null;type:text;comment:内容"`       // 反馈详细内容
	Contact   string  `json:"contact" gorm:"size:100;comment:联系方式"`               // 联系方式，便于后续跟进
	Status    string  `json:"status" gorm:"default:'pending';size:20;comment:状态"` // 处理状态：pending/processing/resolved
}

// TableName 指定表名为单数形式
func (Feedback) TableName() string {
	return "feedback"
}

// OperationLog 操作日志表，记录管理端的关键操作审计日志
// 用于安全审计和问题追溯，记录操作人、操作类型、操作对象及IP等信息
type OperationLog struct {
	BaseModel
	UserID   uint64 `json:"userId" gorm:"not null;index;comment:用户ID"`   // 操作人ID，关联 User 表
	User     User   `json:"user,omitempty"`                              // 关联的操作人信息，查询时预加载
	Action   string `json:"action" gorm:"not null;size:50;comment:操作类型"` // 操作类型，如 create、update、delete
	Target   string `json:"target" gorm:"size:50;comment:操作对象"`          // 操作对象类型，如 poem、user、role
	TargetID uint64 `json:"targetId" gorm:"comment:对象ID"`                // 操作对象的主键ID
	Detail   string `json:"detail" gorm:"type:text;comment:详情"`          // 操作详细描述，JSON 格式存储变更内容
	IP       string `json:"ip" gorm:"size:50;comment:IP地址"`              // 操作者的IP地址
}

// TableName 指定表名为单数形式
func (OperationLog) TableName() string {
	return "operation_log"
}
