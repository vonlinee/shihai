package apidocs

import (
	"shihai/internal/dto"
	"shihai/pkg/utils"
)

var (
	_ dto.RoleCreateRequest
	_ utils.Response
)

// swaggerListRoles documents GET /api/rbac/roles.
// @Summary 查询角色列表
// @Tags RBAC 角色
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles [get]
func swaggerListRoles() {}

// swaggerGetRole documents GET /api/rbac/roles/{id}.
// @Summary 获取角色详情
// @Tags RBAC 角色
// @Security BearerAuth
// @Param id path string true "角色 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/rbac/roles/{id} [get]
func swaggerGetRole() {}

// swaggerCreateRole documents POST /api/rbac/roles.
// @Summary 创建角色
// @Tags RBAC 角色
// @Security BearerAuth
// @Param request body dto.RoleCreateRequest true "角色信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles [post]
func swaggerCreateRole() {}

// swaggerUpdateRole documents PUT /api/rbac/roles/{id}.
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
func swaggerUpdateRole() {}

// swaggerDeleteRole documents DELETE /api/rbac/roles/{id}.
// @Summary 删除角色
// @Tags RBAC 角色
// @Security BearerAuth
// @Param id path string true "角色 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles/{id} [delete]
func swaggerDeleteRole() {}

// swaggerGetRolePermissions documents GET /api/rbac/roles/{id}/permissions.
// @Summary 查询角色权限
// @Tags RBAC 角色
// @Security BearerAuth
// @Param id path string true "角色 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles/{id}/permissions [get]
func swaggerGetRolePermissions() {}

// swaggerAssignRolePermissions documents PUT /api/rbac/roles/{id}/permissions.
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
func swaggerAssignRolePermissions() {}

// swaggerGetUserRoles documents GET /api/rbac/users/{id}/roles.
// @Summary 查询用户角色
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/users/{id}/roles [get]
func swaggerGetUserRoles() {}

// swaggerAssignUserRoles documents PUT /api/rbac/users/{id}/roles.
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
func swaggerAssignUserRoles() {}

// swaggerGetUserPermissions documents GET /api/rbac/users/{id}/permissions.
// @Summary 查询用户权限
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/users/{id}/permissions [get]
func swaggerGetUserPermissions() {}

// swaggerGetMyPermissions documents GET /api/rbac/my/permissions.
// @Summary 查询当前用户权限
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/rbac/my/permissions [get]
func swaggerGetMyPermissions() {}

// swaggerCheckPermission documents POST /api/rbac/check.
// @Summary 检查当前用户权限
// @Tags RBAC 用户授权
// @Security BearerAuth
// @Param request body dto.CheckPermissionRequest true "权限编码"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/rbac/check [post]
func swaggerCheckPermission() {}

// swaggerListPermissions documents GET /api/rbac/permissions.
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
func swaggerListPermissions() {}

// swaggerGetAllPermissions documents GET /api/rbac/permissions/all.
// @Summary 查询全部权限
// @Tags RBAC 权限
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions/all [get]
func swaggerGetAllPermissions() {}

// swaggerGetPermission documents GET /api/rbac/permissions/{id}.
// @Summary 获取权限详情
// @Tags RBAC 权限
// @Security BearerAuth
// @Param id path string true "权限 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/rbac/permissions/{id} [get]
func swaggerGetPermission() {}

// swaggerCreatePermission documents POST /api/rbac/permissions.
// @Summary 创建权限
// @Tags RBAC 权限
// @Security BearerAuth
// @Param request body dto.PermissionCreateRequest true "权限信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions [post]
func swaggerCreatePermission() {}

// swaggerUpdatePermission documents PUT /api/rbac/permissions/{id}.
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
func swaggerUpdatePermission() {}

// swaggerDeletePermission documents DELETE /api/rbac/permissions/{id}.
// @Summary 删除权限
// @Tags RBAC 权限
// @Security BearerAuth
// @Param id path string true "权限 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions/{id} [delete]
func swaggerDeletePermission() {}
