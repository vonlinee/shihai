package models

// Role 角色表，存储系统角色定义
// 系统预设四种角色：admin(超管)、editor(编辑)、reviewer(审核员)、user(普通用户)
// 角色通过 RolePermission 关联权限编码，通过 UserRole 关联用户
type Role struct {
	BaseModel
	Name        string `json:"name" gorm:"not null;uniqueIndex;size:50;comment:角色名称"` // 角色名称，唯一索引，如 admin、editor
	Description string `json:"description" gorm:"size:255;comment:角色描述"`              // 角色功能描述
	IsActive    bool   `json:"isActive" gorm:"default:true;comment:是否启用"`             // 是否启用，false 表示角色被禁用
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}

// RolePermission 角色权限关联表 - 存储角色拥有的权限编码
// 采用 (RoleID, Code) 联合唯一索引，同一角色下权限编码不可重复
// Code 字段直接存储权限编码字符串（如 "user:create"），而非权限表ID
type RolePermission struct {
	BaseModel
	RoleID uint64 `json:"roleId" gorm:"not null;index:idx_role_perm,unique;comment:角色ID"`        // 所属角色ID，关联 Role 表
	Role   Role   `json:"role,omitempty"`                                                        // 关联的角色信息，查询时预加载
	Code   string `json:"code" gorm:"not null;index:idx_role_perm,unique;size:100;comment:权限编码"` // 权限编码，如 "user:create"、"poem:read"
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permission"
}

// Permission 权限表 - 存储权限定义
// 权限编码格式为 "模块:操作"，如 "user:create"、"poem:list"
// 权限按模块(Module)分组，方便管理和展示
type Permission struct {
	BaseModel
	Code        string `json:"code" gorm:"not null;uniqueIndex;size:100;comment:权限编码"` // 权限编码，唯一索引，格式 "模块:操作"
	Name        string `json:"name" gorm:"not null;size:100;comment:权限名称"`             // 权限中文显示名称，如 "创建用户"
	Description string `json:"description" gorm:"size:255;comment:权限描述"`               // 权限详细描述
	Module      string `json:"module" gorm:"not null;size:50;index;comment:所属模块"`      // 所属业务模块，如 user、poem、comment
	IsActive    bool   `json:"isActive" gorm:"default:true;comment:是否启用"`              // 是否启用，false 表示权限被禁用
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permission"
}

// UserRole 用户角色关联表
// 采用 (UserID, RoleID) 联合唯一索引，同一用户不可重复分配同一角色
// 一个用户可拥有多个角色，角色权限取并集
type UserRole struct {
	BaseModel
	UserID uint64 `json:"userId" gorm:"not null;index:idx_user_role,unique;comment:用户ID"` // 用户ID，关联 User 表
	User   User   `json:"user,omitempty"`                                                 // 关联的用户信息，查询时预加载
	RoleID uint64 `json:"roleId" gorm:"not null;index:idx_user_role,unique;comment:角色ID"` // 角色ID，关联 Role 表
	Role   Role   `json:"role,omitempty"`                                                 // 关联的角色信息，查询时预加载
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_role"
}

// RoleAdmin/Editor/Reviewer/User 默认角色名常量
// 避免在业务代码中使用字符串字面量，同时为 DefaultRoles 提供唯一来源
const (
	RoleAdmin    = "admin"    // 超级管理员
	RoleEditor   = "editor"   // 编辑
	RoleReviewer = "reviewer" // 审核员
	RoleUser     = "user"     // 普通用户
)

// DefaultRoles 默认角色定义
// 系统启动时自动初始化，确保基础角色存在
//   - admin: 超级管理员，启动时自动获取 AllPermissions 中的全部权限
//   - editor: 编辑，可管理诗词内容和评论
//   - reviewer: 审核员，可审核纠错申请
//   - user: 普通用户，仅拥有基础权限
var DefaultRoles = []Role{
	{Name: RoleAdmin, Description: "超级管理员，拥有所有权限", IsActive: true},
	{Name: RoleEditor, Description: "编辑，可管理诗词内容", IsActive: true},
	{Name: RoleReviewer, Description: "审核员，可审核纠错申请", IsActive: true},
	{Name: RoleUser, Description: "普通用户，基础权限", IsActive: true},
}

// RolePermissionMap 默认角色权限映射
// 系统启动时根据此映射自动为非 admin 角色分配权限
// 每次启动会重新同步，确保角色拥有最新的权限配置
//
// 重要说明：
//   - admin 角色不在此映射中出现，他会在 InitDefaultRolesAndPermissions 中
//     自动获取 AllPermissions 中的全部权限，新增的接口权限会自动同步。
//   - 本映射中使用的所有权限编码都必须是 permission_codes.go 中定义的常量，
//     禁止使用 "user:create" 这样的字面量。
//
// 权限划分：
//   - editor: 诗词管理 + 评论审核 + 纠错审核 + 公告管理 + 论坛管理 + 题目管理 + 权限查看
//   - reviewer: 诗词查看 + 评论审核 + 纠错审核 + 论坛查看 + 题目查看 + 权限查看
//   - user: 基础查看 + 评论 + 纠错提交 + 公告查看 + 论坛 + 题目 + 反馈
var RolePermissionMap = map[string][]string{
	RoleEditor: {
		PermPermissionList, PermPermissionRead,
		PermPoemCreate, PermPoemRead, PermPoemUpdate, PermPoemDelete, PermPoemList,
		PermWorkCollectionCreate, PermWorkCollectionRead, PermWorkCollectionUpdate, PermWorkCollectionDelete, PermWorkCollectionList, PermWorkCollectionItemManage,
		PermCommentRead, PermCommentDelete, PermCommentList, PermCommentModerate,
		PermCorrectionRead, PermCorrectionReview, PermCorrectionList,
		PermAnnouncementCreate, PermAnnouncementRead, PermAnnouncementUpdate, PermAnnouncementList,
		PermForumRead, PermForumUpdate, PermForumDelete, PermForumList, PermForumModerate,
		PermQuizCreate, PermQuizRead, PermQuizUpdate, PermQuizDelete, PermQuizList,
		PermFeedbackRead, PermFeedbackList,
	},
	RoleReviewer: {
		PermPermissionList, PermPermissionRead,
		PermPoemRead, PermPoemList,
		PermCommentRead, PermCommentList, PermCommentModerate,
		PermCorrectionRead, PermCorrectionReview, PermCorrectionList,
		PermForumRead, PermForumList, PermForumModerate,
		PermQuizRead, PermQuizList,
		PermFeedbackRead, PermFeedbackList,
	},
	RoleUser: {
		PermUserRead, PermUserUpdate,
		PermPoemRead, PermPoemList,
		PermCommentCreate, PermCommentRead, PermCommentUpdate, PermCommentList,
		PermCorrectionCreate, PermCorrectionRead,
		PermAnnouncementRead, PermAnnouncementList,
		PermForumCreate, PermForumRead, PermForumList,
		PermQuizRead, PermQuizList,
		PermFeedbackCreate, PermFeedbackRead,
	},
}

// 权限点定义请参见 permission_codes.go。
// 启动同步逻辑在 services/rbac_service.go 的 InitDefaultRolesAndPermissions 中，
// 会从 AllPermissions 读取所有权限点并 UPSERT 到 permission 表。
