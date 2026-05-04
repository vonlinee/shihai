package dto

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"` // 角色名称，必填，如 admin、editor
	Description string `json:"description" binding:"max=255"`  // 角色功能描述
}

// RoleUpdateRequest 更新角色请求，IsActive 使用指针以区分未提供和 false
type RoleUpdateRequest struct {
	Name        string `json:"name" binding:"max=50"`         // 角色名称
	Description string `json:"description" binding:"max=255"` // 角色描述
	IsActive    *bool  `json:"isActive"`                      // 是否启用，指针类型以区分未提供和 false
}

// RoleResponse 角色响应，包含角色信息和关联的权限编码列表
type RoleResponse struct {
	ID          uint64   `json:"id"`                    // 角色ID
	Name        string   `json:"name"`                  // 角色名称
	Description string   `json:"description"`           // 角色描述
	IsActive    bool     `json:"isActive"`              // 是否启用
	Permissions []string `json:"permissions,omitempty"` // 权限编码列表，如 ["user:create", "poem:read"]
}

// RoleListRequest 角色列表请求，支持分页
type RoleListRequest struct {
	Page     int `form:"page,default=1"`      // 页码，默认第1页
	PageSize int `form:"pageSize,default=10"` // 每页数量，默认10条
}

// AssignPermissionRequest 为角色分配权限请求，传入权限编码列表
// 会替换该角色的所有权限，而非追加
type AssignPermissionRequest struct {
	Permissions []string `json:"permissions" binding:"required"` // 权限编码列表，如 ["user:create", "poem:read"]
}

// AssignRoleRequest 分配角色请求，为用户指定角色ID列表
// 会替换该用户的所有角色，而非追加
type AssignRoleRequest struct {
	RoleIDs []uint64 `json:"roleIds" binding:"required"` // 角色ID列表
}

// UserRoleResponse 用户角色响应，返回用户的所有角色信息
type UserRoleResponse struct {
	UserID uint64         `json:"userId"` // 用户ID
	Roles  []RoleResponse `json:"roles"`  // 角色列表
}

// CheckPermissionRequest 检查权限请求，用于验证当前用户是否拥有指定权限
type CheckPermissionRequest struct {
	Permission string `json:"permission" binding:"required"` // 权限编码，如 "user:create"
}

// PermissionCreateRequest 创建权限请求
type PermissionCreateRequest struct {
	Code        string `json:"code" binding:"required,max=100"`  // 权限编码，必填，格式 "模块:操作"
	Name        string `json:"name" binding:"required,max=100"`  // 权限名称，必填，中文显示名
	Description string `json:"description" binding:"max=255"`    // 权限描述
	Module      string `json:"module" binding:"required,max=50"` // 所属模块，必填，如 user、poem
}

// PermissionUpdateRequest 更新权限请求，IsActive 使用指针以区分未提供和 false
type PermissionUpdateRequest struct {
	Name        string `json:"name" binding:"max=100"`        // 权限名称
	Description string `json:"description" binding:"max=255"` // 权限描述
	Module      string `json:"module" binding:"max=50"`       // 所属模块
	IsActive    *bool  `json:"isActive"`                      // 是否启用，指针类型以区分未提供和 false
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uint64 `json:"id"`          // 权限ID
	Code        string `json:"code"`        // 权限编码
	Name        string `json:"name"`        // 权限名称
	Description string `json:"description"` // 权限描述
	Module      string `json:"module"`      // 所属模块
	IsActive    bool   `json:"isActive"`    // 是否启用
}

// PermissionListRequest 权限列表请求，支持分页和按模块筛选
type PermissionListRequest struct {
	Page     int    `form:"page,default=1"`      // 页码，默认第1页
	PageSize int    `form:"pageSize,default=50"` // 每页数量，默认50条
	Module   string `form:"module"`              // 按模块筛选，如 user、poem
}
