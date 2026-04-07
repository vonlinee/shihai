package dto

import "time"

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Name     string `json:"name" binding:"max=50"`
}

// AdminCreateUserRequest 管理员创建用户请求
type AdminCreateUserRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Password string   `json:"password" binding:"required,min=6,max=100"`
	Name     string   `json:"name" binding:"max=50"`
	RoleIds  []uint64 `json:"roleIds"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Name   string `json:"name" binding:"max=50"`
	Avatar string `json:"avatar" binding:"max=255"`
	Gender string `json:"gender" binding:"max=10"`
	Age    int    `json:"age"`
	Phone  string `json:"phone" binding:"max=20"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=100"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Gender    string    `json:"gender"`
	Age       int       `json:"age"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Roles     []string  `json:"roles"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
	Role     string `form:"role"`
}
