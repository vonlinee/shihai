package dto

import "time"

// RegisterRequest 注册请求，用于普通用户注册
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`  // 用户名，3-50字符
	Password string `json:"password" binding:"required,min=6,max=100"` // 密码，6-100字符
	Name     string `json:"name" binding:"max=50"`                     // 显示名称，可选，最多50字符
}

// AdminCreateUserRequest 管理员创建用户请求，支持指定角色
type AdminCreateUserRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`  // 用户名，3-50字符
	Password string   `json:"password" binding:"required,min=6,max=100"` // 密码，6-100字符
	Name     string   `json:"name" binding:"max=50"`                     // 显示名称，可选
	RoleIds  []uint64 `json:"roleIds"`                                   // 角色ID列表，为空则分配默认 user 角色
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// UpdateUserRequest 更新用户信息请求，用户只能修改自己的基本信息
type UpdateUserRequest struct {
	Name   string `json:"name" binding:"max=50"`    // 显示名称
	Avatar string `json:"avatar" binding:"max=255"` // 头像URL
	Gender string `json:"gender" binding:"max=10"`  // 性别：male/female/other
	Age    int    `json:"age"`                      // 年龄
	Phone  string `json:"phone" binding:"max=20"`   // 手机号码
}

// ChangePasswordRequest 修改密码请求，需要验证旧密码
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`               // 当前密码
	NewPassword string `json:"newPassword" binding:"required,min=6,max=100"` // 新密码，6-100字符
}

// UserResponse 用户响应，包含用户基本信息和角色信息
// Role 为主角色（优先级：admin > editor > reviewer > user），Roles 为所有角色列表
type UserResponse struct {
	ID        uint64    `json:"id"`        // 用户ID
	Username  string    `json:"username"`  // 用户名
	Name      string    `json:"name"`      // 显示名称
	Avatar    string    `json:"avatar"`    // 头像URL
	Gender    string    `json:"gender"`    // 性别
	Age       int       `json:"age"`       // 年龄
	Phone     string    `json:"phone"`     // 手机号
	Role      string    `json:"role"`      // 主角色，用于前端简单判断
	Roles     []string  `json:"roles"`     // 所有角色名称列表
	IsActive  bool      `json:"isActive"`  // 是否启用
	CreatedAt time.Time `json:"createdAt"` // 注册时间
}

// LoginResponse 登录响应，返回 JWT Token 和用户信息
type LoginResponse struct {
	Token string       `json:"token"` // JWT 认证令牌
	User  UserResponse `json:"user"`  // 登录用户信息
}

// UserListRequest 用户列表请求，支持分页、关键词搜索和角色筛选
type UserListRequest struct {
	Page     int    `form:"page,default=1"`      // 页码，默认第1页
	PageSize int    `form:"pageSize,default=10"` // 每页数量，默认10条
	Keyword  string `form:"keyword"`             // 搜索关键词，匹配用户名或姓名
	Role     string `form:"role"`                // 按角色筛选，如 admin、editor
}
