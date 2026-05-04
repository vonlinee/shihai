package middleware

import (
	"shihai/internal/services"
	"shihai/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware RBAC中间件
type RBACMiddleware struct {
	rbacService *services.RBACService
}

func NewRBACMiddleware(rbacService *services.RBACService) *RBACMiddleware {
	return &RBACMiddleware{rbacService: rbacService}
}

// RequirePermission 检查权限中间件
func (m *RBACMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		hasPermission, err := m.rbacService.CheckUserPermission(userID.(uint64), permission)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			c.Abort()
			return
		}

		if !hasPermission {
			utils.Forbidden(c, "没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 检查任意权限中间件（满足其中一个即可）
func (m *RBACMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		userPermissions, err := m.rbacService.GetUserPermissions(userID.(uint64))
		if err != nil {
			utils.InternalServerError(c, err.Error())
			c.Abort()
			return
		}

		// 检查是否拥有任意一个权限
		for _, required := range permissions {
			for _, userPerm := range userPermissions {
				if userPerm == required {
					c.Next()
					return
				}
			}
		}

		utils.Forbidden(c, "没有权限执行此操作")
		c.Abort()
	}
}

// RequireAllPermissions 检查所有权限中间件（必须满足所有）
func (m *RBACMiddleware) RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		userPermissions, err := m.rbacService.GetUserPermissions(userID.(uint64))
		if err != nil {
			utils.InternalServerError(c, err.Error())
			c.Abort()
			return
		}

		// 创建权限集合
		permSet := make(map[string]bool)
		for _, p := range userPermissions {
			permSet[p] = true
		}

		// 检查是否拥有所有权限
		for _, required := range permissions {
			if !permSet[required] {
				utils.Forbidden(c, "没有权限执行此操作")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireRole 检查角色中间件
func (m *RBACMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		hasRole, err := m.rbacService.CheckUserRole(userID.(uint64), role)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			c.Abort()
			return
		}

		if !hasRole {
			utils.Forbidden(c, "没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole 检查任意角色中间件（满足其中一个即可）
func (m *RBACMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		userRoles, err := m.rbacService.GetUserRoles(userID.(uint64))
		if err != nil {
			utils.InternalServerError(c, err.Error())
			c.Abort()
			return
		}

		// 检查是否拥有任意一个角色
		for _, required := range roles {
			for _, userRole := range userRoles {
				if userRole.Name == required {
					c.Next()
					return
				}
			}
		}

		utils.Forbidden(c, "没有权限执行此操作")
		c.Abort()
	}
}

// AdminOnly 仅管理员访问（保留原有功能，使用RBAC实现）
func (m *RBACMiddleware) AdminOnly() gin.HandlerFunc {
	return m.RequireRole("admin")
}

// PermissionChecker 权限检查器，用于在handler中检查权限
type PermissionChecker struct {
	c           *gin.Context
	rbacService *services.RBACService
}

// NewPermissionChecker 创建权限检查器
func (m *RBACMiddleware) NewPermissionChecker(c *gin.Context) *PermissionChecker {
	return &PermissionChecker{
		c:           c,
		rbacService: m.rbacService,
	}
}

// HasPermission 检查是否有指定权限
func (pc *PermissionChecker) HasPermission(permission string) bool {
	userID, exists := pc.c.Get("userID")
	if !exists {
		return false
	}

	has, _ := pc.rbacService.CheckUserPermission(userID.(uint64), permission)
	return has
}

// HasRole 检查是否有指定角色
func (pc *PermissionChecker) HasRole(role string) bool {
	userID, exists := pc.c.Get("userID")
	if !exists {
		return false
	}

	has, _ := pc.rbacService.CheckUserRole(userID.(uint64), role)
	return has
}

// GetUserID 获取用户ID
func (pc *PermissionChecker) GetUserID() (uint64, bool) {
	userID, exists := pc.c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(uint64), true
}

// Middleware 创建gin中间件函数
func (m *RBACMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将RBAC服务存入上下文，供后续使用
		c.Set("rbacService", m.rbacService)
		c.Next()
	}
}

// GetRBACService 从上下文获取RBAC服务
func GetRBACService(c *gin.Context) *services.RBACService {
	service, exists := c.Get("rbacService")
	if !exists {
		return nil
	}
	return service.(*services.RBACService)
}
