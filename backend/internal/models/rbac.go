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

// RolePermission 角色权限关联表 - 存储角色拥有的权限编码
type RolePermission struct {
	BaseModel
	RoleID uint64 `json:"roleId" gorm:"not null;index:idx_role_perm,unique;comment:角色ID"`
	Role   Role   `json:"role,omitempty"`
	Code   string `json:"code" gorm:"not null;index:idx_role_perm,unique;size:100;comment:权限编码"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permission"
}

// Permission 权限表 - 存储权限定义
type Permission struct {
	BaseModel
	Code        string `json:"code" gorm:"not null;uniqueIndex;size:100;comment:权限编码"`
	Name        string `json:"name" gorm:"not null;size:100;comment:权限名称"`
	Description string `json:"description" gorm:"size:255;comment:权限描述"`
	Module      string `json:"module" gorm:"not null;size:50;index;comment:所属模块"`
	IsActive    bool   `json:"isActive" gorm:"default:true;comment:是否启用"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permission"
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

// DefaultRoles 默认角色定义
var DefaultRoles = []Role{
	{Name: "admin", Description: "超级管理员，拥有所有权限", IsActive: true},
	{Name: "editor", Description: "编辑，可管理诗词内容", IsActive: true},
	{Name: "reviewer", Description: "审核员，可审核纠错申请", IsActive: true},
	{Name: "user", Description: "普通用户，基础权限", IsActive: true},
}

// DefaultPermissions 默认权限定义
var DefaultPermissions = []Permission{
	// user module
	{Code: "user:create", Name: "创建用户", Module: "user", Description: "创建新用户"},
	{Code: "user:read", Name: "查看用户", Module: "user", Description: "查看用户信息"},
	{Code: "user:update", Name: "更新用户", Module: "user", Description: "更新用户信息"},
	{Code: "user:delete", Name: "删除用户", Module: "user", Description: "删除用户"},
	{Code: "user:list", Name: "用户列表", Module: "user", Description: "查看用户列表"},
	// role module
	{Code: "role:create", Name: "创建角色", Module: "role", Description: "创建新角色"},
	{Code: "role:read", Name: "查看角色", Module: "role", Description: "查看角色信息"},
	{Code: "role:update", Name: "更新角色", Module: "role", Description: "更新角色信息"},
	{Code: "role:delete", Name: "删除角色", Module: "role", Description: "删除角色"},
	{Code: "role:list", Name: "角色列表", Module: "role", Description: "查看角色列表"},
	{Code: "role:assign", Name: "分配角色", Module: "role", Description: "为用户分配角色"},
	// permission module
	{Code: "permission:list", Name: "权限列表", Module: "permission", Description: "查看权限列表"},
	{Code: "permission:read", Name: "查看权限", Module: "permission", Description: "查看权限详情"},
	{Code: "permission:create", Name: "创建权限", Module: "permission", Description: "创建新权限"},
	{Code: "permission:update", Name: "更新权限", Module: "permission", Description: "更新权限信息"},
	{Code: "permission:delete", Name: "删除权限", Module: "permission", Description: "删除权限"},
	// poem module
	{Code: "poem:create", Name: "创建诗词", Module: "poem", Description: "添加新诗词"},
	{Code: "poem:read", Name: "查看诗词", Module: "poem", Description: "查看诗词详情"},
	{Code: "poem:update", Name: "更新诗词", Module: "poem", Description: "编辑诗词信息"},
	{Code: "poem:delete", Name: "删除诗词", Module: "poem", Description: "删除诗词"},
	{Code: "poem:list", Name: "诗词列表", Module: "poem", Description: "查看诗词列表"},
	// comment module
	{Code: "comment:create", Name: "创建评论", Module: "comment", Description: "发表评论"},
	{Code: "comment:read", Name: "查看评论", Module: "comment", Description: "查看评论"},
	{Code: "comment:update", Name: "更新评论", Module: "comment", Description: "更新评论"},
	{Code: "comment:delete", Name: "删除评论", Module: "comment", Description: "删除评论"},
	{Code: "comment:list", Name: "评论列表", Module: "comment", Description: "查看评论列表"},
	{Code: "comment:moderate", Name: "审核评论", Module: "comment", Description: "审核管理评论"},
	// correction module
	{Code: "correction:create", Name: "提交纠错", Module: "correction", Description: "提交纠错申请"},
	{Code: "correction:read", Name: "查看纠错", Module: "correction", Description: "查看纠错详情"},
	{Code: "correction:update", Name: "更新纠错", Module: "correction", Description: "更新纠错信息"},
	{Code: "correction:review", Name: "审核纠错", Module: "correction", Description: "审核纠错申请"},
	{Code: "correction:list", Name: "纠错列表", Module: "correction", Description: "查看纠错列表"},
	// announcement module
	{Code: "announcement:create", Name: "创建公告", Module: "announcement", Description: "创建公告"},
	{Code: "announcement:read", Name: "查看公告", Module: "announcement", Description: "查看公告"},
	{Code: "announcement:update", Name: "更新公告", Module: "announcement", Description: "更新公告"},
	{Code: "announcement:delete", Name: "删除公告", Module: "announcement", Description: "删除公告"},
	{Code: "announcement:list", Name: "公告列表", Module: "announcement", Description: "查看公告列表"},
	// system module
	{Code: "system:config", Name: "系统配置", Module: "system", Description: "系统配置管理"},
	{Code: "system:log", Name: "系统日志", Module: "system", Description: "查看系统日志"},
	{Code: "system:backup", Name: "系统备份", Module: "system", Description: "系统备份管理"},
}

// RolePermissionMap 默认角色权限映射
var RolePermissionMap = map[string][]string{
	"admin": {
		"user:create", "user:read", "user:update", "user:delete", "user:list",
		"role:create", "role:read", "role:update", "role:delete", "role:list", "role:assign",
		"permission:list", "permission:read", "permission:create", "permission:update", "permission:delete",
		"poem:create", "poem:read", "poem:update", "poem:delete", "poem:list",
		"comment:create", "comment:read", "comment:update", "comment:delete", "comment:list", "comment:moderate",
		"correction:create", "correction:read", "correction:update", "correction:review", "correction:list",
		"announcement:create", "announcement:read", "announcement:update", "announcement:delete", "announcement:list",
		"system:config", "system:log", "system:backup",
	},
	"editor": {
		"permission:list", "permission:read",
		"poem:create", "poem:read", "poem:update", "poem:delete", "poem:list",
		"comment:read", "comment:delete", "comment:list", "comment:moderate",
		"correction:read", "correction:review", "correction:list",
		"announcement:create", "announcement:read", "announcement:update", "announcement:list",
	},
	"reviewer": {
		"permission:list", "permission:read",
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
