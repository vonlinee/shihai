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
// @Summary 创建角色
// @Tags RBAC 角色
// @Security BearerAuth
// @Param request body dto.RoleCreateRequest true "角色信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles [post]
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
// @Summary 查询角色列表
// @Tags RBAC 角色
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles [get]
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
// @Summary 获取角色详情
// @Tags RBAC 角色
// @Security BearerAuth
// @Param id path string true "角色 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/rbac/roles/{id} [get]
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
// @Summary 更新角色
// @Tags RBAC 角色
// @Security BearerAuth
// @Param id path string true "角色 ID"
// @Param request body dto.RoleUpdateRequest true "角色信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles/{id} [put]
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
// @Summary 删除角色
// @Tags RBAC 角色
// @Security BearerAuth
// @Param id path string true "角色 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles/{id} [delete]
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

// ==================== Role-Permission Handlers ====================

// AssignPermissionsToRole 为角色分配权限
// @Summary 分配角色权限
// @Tags RBAC 角色
// @Security BearerAuth
// @Param id path string true "角色 ID"
// @Param request body dto.AssignPermissionRequest true "权限编码列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles/{id}/permissions [put]
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

	if err := h.rbacService.AssignPermissionsToRole(roleID, req.Permissions); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "permissions assigned successfully", nil)
}

// GetRolePermissions 获取角色的权限列表
// @Summary 查询角色权限
// @Tags RBAC 角色
// @Security BearerAuth
// @Param id path string true "角色 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles/{id}/permissions [get]
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
// @Summary 分配用户角色
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Param request body dto.AssignRoleRequest true "角色 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/users/{id}/roles [put]
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

	if err := h.rbacService.AssignRolesToUser(userID, []uint64(req.RoleIDs)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "roles assigned successfully", nil)
}

// GetUserRoles 获取用户的角色列表
// @Summary 查询用户角色
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/users/{id}/roles [get]
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
// @Summary 查询用户权限
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/users/{id}/permissions [get]
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
// @Summary 检查当前用户权限
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Param request body dto.CheckPermissionRequest true "权限编码"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/rbac/check [post]
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
// @Summary 查询当前用户权限
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/rbac/my/permissions [get]
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

// ==================== Permission Handlers ====================

// GetPermissionList 获取权限列表
// @Summary 查询权限列表
// @Tags RBAC 权限
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(50)
// @Param module query string false "模块"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions [get]
func (h *RBACHandler) GetPermissionList(c *gin.Context) {
	var req dto.PermissionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	perms, total, err := h.rbacService.GetPermissionList(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageSuccess(c, perms, total, req.Page, req.PageSize)
}

// GetAllPermissions 获取所有权限
// @Summary 查询全部权限
// @Tags RBAC 权限
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions/all [get]
func (h *RBACHandler) GetAllPermissions(c *gin.Context) {
	perms, err := h.rbacService.GetAllPermissions()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, perms)
}

// GetPermissionByID 根据ID获取权限
// @Summary 获取权限详情
// @Tags RBAC 权限
// @Security BearerAuth
// @Param id path string true "权限 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/rbac/permissions/{id} [get]
func (h *RBACHandler) GetPermissionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid permission id")
		return
	}

	perm, err := h.rbacService.GetPermissionByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, perm)
}

// CreatePermission 创建权限
// @Summary 创建权限
// @Tags RBAC 权限
// @Security BearerAuth
// @Param request body dto.PermissionCreateRequest true "权限信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions [post]
func (h *RBACHandler) CreatePermission(c *gin.Context) {
	var req dto.PermissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	perm, err := h.rbacService.CreatePermission(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, perm)
}

// UpdatePermission 更新权限
// @Summary 更新权限
// @Tags RBAC 权限
// @Security BearerAuth
// @Param id path string true "权限 ID"
// @Param request body dto.PermissionUpdateRequest true "权限信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions/{id} [put]
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

	perm, err := h.rbacService.UpdatePermission(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, perm)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Tags RBAC 权限
// @Security BearerAuth
// @Param id path string true "权限 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions/{id} [delete]
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
