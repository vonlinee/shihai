package models

// Role 角色表
type Role struct {
	BaseModel
	Name        string `json:"name" gorm:"not null;uniqueIndex;size:50;comment:角色名称"` // admin, editor, reviewer, user
	Description string `json:"description" gorm:"size:255;comment:角色描述"`
	IsActive    bool   `json:"isActive" gorm:"default:true;comment:是否启用"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}

// Permission 权限表
type Permission struct {
	BaseModel
	Code        string `json:"code" gorm:"not null;uniqueIndex;size:100;comment:权限编码"` // user:create, user:read, user:update, user:delete
	Name        string `json:"name" gorm:"not null;size:100;comment:权限名称"`
	Description string `json:"description" gorm:"size:255;comment:权限描述"`
	Module      string `json:"module" gorm:"size:50;comment:所属模块"` // user, poem, comment, system
	IsActive    bool   `json:"isActive" gorm:"default:true;comment:是否启用"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permission"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	BaseModel
	RoleID       uint64     `json:"roleId" gorm:"not null;index:idx_role_perm,unique;comment:角色ID"`
	Role         Role       `json:"role,omitempty"`
	PermissionID uint64     `json:"permissionId" gorm:"not null;index:idx_role_perm,unique;comment:权限ID"`
	Permission   Permission `json:"permission,omitempty"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permission"
}

// UserRole 用户角色关联表
type UserRole struct {
	BaseModel
	UserID uint64 `json:"userId" gorm:"not null;index:idx_user_role,unique;comment:用户ID"`
	User   User   `json:"user,omitempty"`
	RoleID uint64 `json:"roleId" gorm:"not null;index:idx_user_role,unique;comment:角色ID"`
	Role   Role   `json:"role,omitempty"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_role"
}

// PermissionCheck 权限检查辅助结构
type PermissionCheck struct {
	UserID      uint64   `json:"userId"`
	RoleIDs     []uint64 `json:"roleIds"`
	Permissions []string `json:"permissions"`
}

// DefaultRoles 默认角色定义
var DefaultRoles = []Role{
	{Name: "admin", Description: "超级管理员，拥有所有权限", IsActive: true},
	{Name: "editor", Description: "编辑，可管理诗词内容", IsActive: true},
	{Name: "reviewer", Description: "审核员，可审核纠错申请", IsActive: true},
	{Name: "user", Description: "普通用户，基础权限", IsActive: true},
}

// DefaultPermissions 默认权限定义
var DefaultPermissions = []Permission{
	// 用户管理权限
	{Code: "user:create", Name: "创建用户", Description: "创建新用户", Module: "user"},
	{Code: "user:read", Name: "查看用户", Description: "查看用户信息", Module: "user"},
	{Code: "user:update", Name: "更新用户", Description: "更新用户信息", Module: "user"},
	{Code: "user:delete", Name: "删除用户", Description: "删除用户", Module: "user"},
	{Code: "user:list", Name: "用户列表", Description: "查看用户列表", Module: "user"},

	// 角色管理权限
	{Code: "role:create", Name: "创建角色", Description: "创建新角色", Module: "role"},
	{Code: "role:read", Name: "查看角色", Description: "查看角色信息", Module: "role"},
	{Code: "role:update", Name: "更新角色", Description: "更新角色信息", Module: "role"},
	{Code: "role:delete", Name: "删除角色", Description: "删除角色", Module: "role"},
	{Code: "role:list", Name: "角色列表", Description: "查看角色列表", Module: "role"},
	{Code: "role:assign", Name: "分配角色", Description: "为用户分配角色", Module: "role"},

	// 权限管理权限
	{Code: "permission:create", Name: "创建权限", Description: "创建新权限", Module: "permission"},
	{Code: "permission:read", Name: "查看权限", Description: "查看权限信息", Module: "permission"},
	{Code: "permission:update", Name: "更新权限", Description: "更新权限信息", Module: "permission"},
	{Code: "permission:delete", Name: "删除权限", Description: "删除权限", Module: "permission"},
	{Code: "permission:list", Name: "权限列表", Description: "查看权限列表", Module: "permission"},
	{Code: "permission:assign", Name: "分配权限", Description: "为角色分配权限", Module: "permission"},

	// 诗词管理权限
	{Code: "poem:create", Name: "创建诗词", Description: "创建新诗词", Module: "poem"},
	{Code: "poem:read", Name: "查看诗词", Description: "查看诗词信息", Module: "poem"},
	{Code: "poem:update", Name: "更新诗词", Description: "更新诗词信息", Module: "poem"},
	{Code: "poem:delete", Name: "删除诗词", Description: "删除诗词", Module: "poem"},
	{Code: "poem:list", Name: "诗词列表", Description: "查看诗词列表", Module: "poem"},

	// 评论管理权限
	{Code: "comment:create", Name: "创建评论", Description: "创建新评论", Module: "comment"},
	{Code: "comment:read", Name: "查看评论", Description: "查看评论信息", Module: "comment"},
	{Code: "comment:update", Name: "更新评论", Description: "更新评论信息", Module: "comment"},
	{Code: "comment:delete", Name: "删除评论", Description: "删除评论", Module: "comment"},
	{Code: "comment:list", Name: "评论列表", Description: "查看评论列表", Module: "comment"},
	{Code: "comment:moderate", Name: "审核评论", Description: "审核评论内容", Module: "comment"},

	// 纠错管理权限
	{Code: "correction:create", Name: "创建纠错", Description: "提交纠错申请", Module: "correction"},
	{Code: "correction:read", Name: "查看纠错", Description: "查看纠错申请", Module: "correction"},
	{Code: "correction:update", Name: "更新纠错", Description: "更新纠错申请", Module: "correction"},
	{Code: "correction:review", Name: "审核纠错", Description: "审核纠错申请", Module: "correction"},
	{Code: "correction:list", Name: "纠错列表", Description: "查看纠错列表", Module: "correction"},

	// 公告管理权限
	{Code: "announcement:create", Name: "创建公告", Description: "创建新公告", Module: "announcement"},
	{Code: "announcement:read", Name: "查看公告", Description: "查看公告信息", Module: "announcement"},
	{Code: "announcement:update", Name: "更新公告", Description: "更新公告信息", Module: "announcement"},
	{Code: "announcement:delete", Name: "删除公告", Description: "删除公告", Module: "announcement"},
	{Code: "announcement:list", Name: "公告列表", Description: "查看公告列表", Module: "announcement"},

	// 系统管理权限
	{Code: "system:config", Name: "系统配置", Description: "管理系统配置", Module: "system"},
	{Code: "system:log", Name: "系统日志", Description: "查看系统日志", Module: "system"},
	{Code: "system:backup", Name: "数据备份", Description: "备份系统数据", Module: "system"},
}

// RolePermissionMap 默认角色权限映射
var RolePermissionMap = map[string][]string{
	"admin": {
		"user:create", "user:read", "user:update", "user:delete", "user:list",
		"role:create", "role:read", "role:update", "role:delete", "role:list", "role:assign",
		"permission:create", "permission:read", "permission:update", "permission:delete", "permission:list", "permission:assign",
		"poem:create", "poem:read", "poem:update", "poem:delete", "poem:list",
		"comment:create", "comment:read", "comment:update", "comment:delete", "comment:list", "comment:moderate",
		"correction:create", "correction:read", "correction:update", "correction:review", "correction:list",
		"announcement:create", "announcement:read", "announcement:update", "announcement:delete", "announcement:list",
		"system:config", "system:log", "system:backup",
	},
	"editor": {
		"poem:create", "poem:read", "poem:update", "poem:delete", "poem:list",
		"comment:read", "comment:delete", "comment:list", "comment:moderate",
		"correction:read", "correction:review", "correction:list",
		"announcement:create", "announcement:read", "announcement:update", "announcement:list",
	},
	"reviewer": {
		"poem:read", "poem:list",
		"comment:read", "comment:list", "comment:moderate",
		"correction:read", "correction:review", "correction:list",
	},
	"user": {
		"user:read", "user:update",
		"poem:read", "poem:list",
		"comment:create", "comment:read", "comment:update", "comment:list",
		"correction:create", "correction:read",
		"announcement:read", "announcement:list",
	},
}
