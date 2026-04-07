package handlers

import (
	"net/http"
	"shihai/internal/dto"
	"shihai/internal/services"
	"shihai/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RBACHandler struct {
	rbacService *services.RBACService
}

func NewRBACHandler(rbacService *services.RBACService) *RBACHandler {
	return &RBACHandler{rbacService: rbacService}
}

// ==================== Role Handlers ====================

// CreateRole 创建角色
func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req dto.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	role, err := h.rbacService.CreateRole(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, role)
}

// GetRoleList 获取角色列表
func (h *RBACHandler) GetRoleList(c *gin.Context) {
	var req dto.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	roles, total, err := h.rbacService.GetRoleList(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageSuccess(c, roles, total, req.Page, req.PageSize)
}

// GetRoleByID 根据ID获取角色
func (h *RBACHandler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	role, err := h.rbacService.GetRoleByID(id)
	if err != nil {
		utils.NotFound(c, "")
		return
	}

	utils.Success(c, role)
}

// UpdateRole 更新角色
func (h *RBACHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	var req dto.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	role, err := h.rbacService.UpdateRole(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, role)
}

// DeleteRole 删除角色
func (h *RBACHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	if err := h.rbacService.DeleteRole(id); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "role deleted successfully", nil)
}

// ==================== Permission Handlers ====================

// CreatePermission 创建权限
func (h *RBACHandler) CreatePermission(c *gin.Context) {
	var req dto.PermissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	permission, err := h.rbacService.CreatePermission(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, permission)
}

// GetPermissionList 获取权限列表
func (h *RBACHandler) GetPermissionList(c *gin.Context) {
	var req dto.PermissionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	permissions, total, err := h.rbacService.GetPermissionList(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageSuccess(c, permissions, total, req.Page, req.PageSize)
}

// GetAllPermissions 获取所有权限（按模块分组）
func (h *RBACHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.rbacService.GetAllPermissions()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, permissions)
}

// GetPermissionByID 根据ID获取权限
func (h *RBACHandler) GetPermissionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid permission id")
		return
	}

	permission, err := h.rbacService.GetPermissionByID(id)
	if err != nil {
		utils.NotFound(c, "")
		return
	}

	utils.Success(c, permission)
}

// UpdatePermission 更新权限
func (h *RBACHandler) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid permission id")
		return
	}

	var req dto.PermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	permission, err := h.rbacService.UpdatePermission(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, permission)
}

// DeletePermission 删除权限
func (h *RBACHandler) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid permission id")
		return
	}

	if err := h.rbacService.DeletePermission(id); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "permission deleted successfully", nil)
}

// ==================== Role-Permission Handlers ====================

// AssignPermissionsToRole 为角色分配权限
func (h *RBACHandler) AssignPermissionsToRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	var req dto.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.rbacService.AssignPermissionsToRole(roleID, req.PermissionIDs); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "permissions assigned successfully", nil)
}

// GetRolePermissions 获取角色的权限列表
func (h *RBACHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	permissions, err := h.rbacService.GetRolePermissions(roleID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, permissions)
}

// ==================== User-Role Handlers ====================

// AssignRolesToUser 为用户分配角色
func (h *RBACHandler) AssignRolesToUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.rbacService.AssignRolesToUser(userID, req.RoleIDs); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "roles assigned successfully", nil)
}

// GetUserRoles 获取用户的角色列表
func (h *RBACHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	roles, err := h.rbacService.GetUserRoles(userID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, roles)
}

// GetUserPermissions 获取用户的权限列表
func (h *RBACHandler) GetUserPermissions(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	permissions, err := h.rbacService.GetUserPermissions(userID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, permissions)
}

// CheckUserPermission 检查当前用户是否有指定权限
func (h *RBACHandler) CheckUserPermission(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "")
		return
	}

	var req dto.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	hasPermission, err := h.rbacService.CheckUserPermission(userID.(uint64), req.Permission)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"hasPermission": hasPermission})
}

// GetMyPermissions 获取当前用户的权限列表
func (h *RBACHandler) GetMyPermissions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "")
		return
	}

	permissions, err := h.rbacService.GetUserPermissions(userID.(uint64))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	roles, err := h.rbacService.GetUserRoles(userID.(uint64))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"permissions": permissions,
		"roles":       roles,
	})
}
