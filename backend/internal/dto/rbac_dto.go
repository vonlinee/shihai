package dto

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=255"`
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	Name        string `json:"name" binding:"max=50"`
	Description string `json:"description" binding:"max=255"`
	IsActive    *bool  `json:"isActive"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsActive    bool     `json:"isActive"`
	Permissions []string `json:"permissions,omitempty"`
}

// RoleListRequest 角色列表请求
type RoleListRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=10"`
}

// PermissionCreateRequest 创建权限请求
type PermissionCreateRequest struct {
	Code        string `json:"code" binding:"required,max=100"`
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=255"`
	Module      string `json:"module" binding:"max=50"`
}

// PermissionUpdateRequest 更新权限请求
type PermissionUpdateRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Description string `json:"description" binding:"max=255"`
	Module      string `json:"module" binding:"max=50"`
	IsActive    *bool  `json:"isActive"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Module      string `json:"module"`
	IsActive    bool   `json:"isActive"`
}

// PermissionListRequest 权限列表请求
type PermissionListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Module   string `form:"module"`
}

// AssignPermissionRequest 分配权限请求
type AssignPermissionRequest struct {
	PermissionIDs []uint64 `json:"permissionIds" binding:"required"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	RoleIDs []uint64 `json:"roleIds" binding:"required"`
}

// UserRoleResponse 用户角色响应
type UserRoleResponse struct {
	UserID uint64         `json:"userId"`
	Roles  []RoleResponse `json:"roles"`
}

// RolePermissionResponse 角色权限响应
type RolePermissionResponse struct {
	RoleID      uint64               `json:"roleId"`
	Permissions []PermissionResponse `json:"permissions"`
}

// CheckPermissionRequest 检查权限请求
type CheckPermissionRequest struct {
	Permission string `json:"permission" binding:"required"`
}

// ModulePermissions 模块权限分组
type ModulePermissions struct {
	Module      string               `json:"module"`
	Permissions []PermissionResponse `json:"permissions"`
}
