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

// RolePermissionMap 默认角色权限映射
var RolePermissionMap = map[string][]string{
	"admin": {
		"user:create", "user:read", "user:update", "user:delete", "user:list",
		"role:create", "role:read", "role:update", "role:delete", "role:list", "role:assign",
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
