package services

import (
	"errors"
	"shihai/internal/dto"
	"shihai/internal/models"
	"shihai/internal/repository"
)

// RBACService RBAC服务
type RBACService struct {
	roleRepo           *repository.RoleRepository
	permissionRepo     *repository.PermissionRepository
	rolePermissionRepo *repository.RolePermissionRepository
	userRoleRepo       *repository.UserRoleRepository
}

func NewRBACService(
	roleRepo *repository.RoleRepository,
	permissionRepo *repository.PermissionRepository,
	rolePermissionRepo *repository.RolePermissionRepository,
	userRoleRepo *repository.UserRoleRepository,
) *RBACService {
	return &RBACService{
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		rolePermissionRepo: rolePermissionRepo,
		userRoleRepo:       userRoleRepo,
	}
}

// ==================== Role Service ====================

// CreateRole 创建角色
func (s *RBACService) CreateRole(req *dto.RoleCreateRequest) (*dto.RoleResponse, error) {
	// 检查角色名是否已存在
	_, err := s.roleRepo.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("role name already exists")
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return s.toRoleResponse(role), nil
}

// GetRoleByID 根据ID获取角色
func (s *RBACService) GetRoleByID(id uint64) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	resp := s.toRoleResponse(role)
	// 获取角色权限
	permissions, err := s.rolePermissionRepo.GetPermissionsByRoleID(id)
	if err == nil {
		resp.Permissions = permissions
	}

	return resp, nil
}

// GetRoleByName 根据名称获取角色
func (s *RBACService) GetRoleByName(name string) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByName(name)
	if err != nil {
		return nil, errors.New("role not found")
	}
	return s.toRoleResponse(role), nil
}

// GetRoleList 获取角色列表
func (s *RBACService) GetRoleList(req *dto.RoleListRequest) ([]dto.RoleResponse, int64, error) {
	roles, total, err := s.roleRepo.List(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.RoleResponse
	for _, role := range roles {
		responses = append(responses, *s.toRoleResponse(&role))
	}

	return responses, total, nil
}

// UpdateRole 更新角色
func (s *RBACService) UpdateRole(id uint64, req *dto.RoleUpdateRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if req.Name != "" {
		// 检查新名称是否已被其他角色使用
		existing, err := s.roleRepo.GetByName(req.Name)
		if err == nil && existing.ID != id {
			return nil, errors.New("role name already exists")
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	return s.toRoleResponse(role), nil
}

// DeleteRole 删除角色
func (s *RBACService) DeleteRole(id uint64) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	// 删除角色的所有权限关联
	if err := s.rolePermissionRepo.DeleteByRoleID(id); err != nil {
		return err
	}

	return s.roleRepo.Delete(id)
}

// ==================== Permission Service ====================

// CreatePermission 创建权限
func (s *RBACService) CreatePermission(req *dto.PermissionCreateRequest) (*dto.PermissionResponse, error) {
	// 检查权限编码是否已存在
	_, err := s.permissionRepo.GetByCode(req.Code)
	if err == nil {
		return nil, errors.New("permission code already exists")
	}

	permission := &models.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Module:      req.Module,
		IsActive:    true,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return s.toPermissionResponse(permission), nil
}

// GetPermissionByID 根据ID获取权限
func (s *RBACService) GetPermissionByID(id uint64) (*dto.PermissionResponse, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("permission not found")
	}
	return s.toPermissionResponse(permission), nil
}

// GetPermissionList 获取权限列表
func (s *RBACService) GetPermissionList(req *dto.PermissionListRequest) ([]dto.PermissionResponse, int64, error) {
	permissions, total, err := s.permissionRepo.List(req.Page, req.PageSize, req.Module)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.PermissionResponse
	for _, permission := range permissions {
		responses = append(responses, *s.toPermissionResponse(&permission))
	}

	return responses, total, nil
}

// GetAllPermissions 获取所有权限（按模块分组）
func (s *RBACService) GetAllPermissions() ([]dto.ModulePermissions, error) {
	permissions, err := s.permissionRepo.ListAll()
	if err != nil {
		return nil, err
	}

	// 按模块分组
	moduleMap := make(map[string][]dto.PermissionResponse)
	for _, p := range permissions {
		resp := s.toPermissionResponse(&p)
		moduleMap[p.Module] = append(moduleMap[p.Module], *resp)
	}

	var result []dto.ModulePermissions
	for module, perms := range moduleMap {
		result = append(result, dto.ModulePermissions{
			Module:      module,
			Permissions: perms,
		})
	}

	return result, nil
}

// UpdatePermission 更新权限
func (s *RBACService) UpdatePermission(id uint64, req *dto.PermissionUpdateRequest) (*dto.PermissionResponse, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("permission not found")
	}

	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Description != "" {
		permission.Description = req.Description
	}
	if req.Module != "" {
		permission.Module = req.Module
	}
	if req.IsActive != nil {
		permission.IsActive = *req.IsActive
	}

	if err := s.permissionRepo.Update(permission); err != nil {
		return nil, err
	}

	return s.toPermissionResponse(permission), nil
}

// DeletePermission 删除权限
func (s *RBACService) DeletePermission(id uint64) error {
	_, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return errors.New("permission not found")
	}
	return s.permissionRepo.Delete(id)
}

// ==================== Role-Permission Service ====================

// AssignPermissionsToRole 为角色分配权限
func (s *RBACService) AssignPermissionsToRole(roleID uint64, permissionIDs []uint64) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	// 清除原有权限
	if err := s.rolePermissionRepo.DeleteByRoleID(roleID); err != nil {
		return err
	}

	// 添加新权限
	for _, permID := range permissionIDs {
		// 检查权限是否存在
		_, err := s.permissionRepo.GetByID(permID)
		if err != nil {
			continue // 跳过不存在的权限
		}

		rp := &models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
		if err := s.rolePermissionRepo.Create(rp); err != nil {
			return err
		}
	}

	return nil
}

// GetRolePermissions 获取角色的权限列表
func (s *RBACService) GetRolePermissions(roleID uint64) ([]dto.PermissionResponse, error) {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	rps, err := s.rolePermissionRepo.ListByRoleID(roleID)
	if err != nil {
		return nil, err
	}

	var responses []dto.PermissionResponse
	for _, rp := range rps {
		if rp.Permission.ID > 0 {
			responses = append(responses, *s.toPermissionResponse(&rp.Permission))
		}
	}

	return responses, nil
}

// ==================== User-Role Service ====================

// AssignRolesToUser 为用户分配角色
func (s *RBACService) AssignRolesToUser(userID uint64, roleIDs []uint64) error {
	// 清除原有角色
	if err := s.userRoleRepo.DeleteByUserID(userID); err != nil {
		return err
	}

	// 添加新角色
	for _, roleID := range roleIDs {
		// 检查角色是否存在
		_, err := s.roleRepo.GetByID(roleID)
		if err != nil {
			continue // 跳过不存在的角色
		}

		ur := &models.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := s.userRoleRepo.Create(ur); err != nil {
			return err
		}
	}

	return nil
}

// GetUserRoles 获取用户的角色列表
func (s *RBACService) GetUserRoles(userID uint64) ([]dto.RoleResponse, error) {
	urs, err := s.userRoleRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.RoleResponse
	for _, ur := range urs {
		if ur.Role.ID > 0 {
			responses = append(responses, *s.toRoleResponse(&ur.Role))
		}
	}

	return responses, nil
}

// GetUserPermissions 获取用户的所有权限
func (s *RBACService) GetUserPermissions(userID uint64) ([]string, error) {
	return s.userRoleRepo.GetUserPermissions(userID)
}

// CheckUserPermission 检查用户是否有指定权限
func (s *RBACService) CheckUserPermission(userID uint64, permissionCode string) (bool, error) {
	permissions, err := s.userRoleRepo.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p == permissionCode {
			return true, nil
		}
	}

	return false, nil
}

// CheckUserRole 检查用户是否有指定角色
func (s *RBACService) CheckUserRole(userID uint64, roleName string) (bool, error) {
	roles, err := s.userRoleRepo.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		if r == roleName {
			return true, nil
		}
	}

	return false, nil
}

// ==================== Helper Methods ====================

func (s *RBACService) toRoleResponse(role *models.Role) *dto.RoleResponse {
	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsActive:    role.IsActive,
	}
}

func (s *RBACService) toPermissionResponse(permission *models.Permission) *dto.PermissionResponse {
	return &dto.PermissionResponse{
		ID:          permission.ID,
		Code:        permission.Code,
		Name:        permission.Name,
		Description: permission.Description,
		Module:      permission.Module,
		IsActive:    permission.IsActive,
	}
}

// ==================== Initialization ====================

// InitDefaultRolesAndPermissions 初始化默认角色和权限
func (s *RBACService) InitDefaultRolesAndPermissions() error {
	// 创建默认权限
	permissionMap := make(map[string]uint64)
	for _, perm := range models.DefaultPermissions {
		existing, err := s.permissionRepo.GetByCode(perm.Code)
		if err != nil {
			// 权限不存在，创建
			newPerm := &models.Permission{
				Code:        perm.Code,
				Name:        perm.Name,
				Description: perm.Description,
				Module:      perm.Module,
				IsActive:    true,
			}
			if err := s.permissionRepo.Create(newPerm); err != nil {
				continue
			}
			permissionMap[perm.Code] = newPerm.ID
		} else {
			permissionMap[perm.Code] = existing.ID
		}
	}

	// 创建默认角色并分配权限
	for _, role := range models.DefaultRoles {
		existing, err := s.roleRepo.GetByName(role.Name)
		var roleID uint64
		if err != nil {
			// 角色不存在，创建
			newRole := &models.Role{
				Name:        role.Name,
				Description: role.Description,
				IsActive:    true,
			}
			if err := s.roleRepo.Create(newRole); err != nil {
				continue
			}
			roleID = newRole.ID
		} else {
			roleID = existing.ID
		}

		// 分配权限给角色
		if permCodes, ok := models.RolePermissionMap[role.Name]; ok {
			var permIDs []uint64
			for _, code := range permCodes {
				if id, ok := permissionMap[code]; ok {
					permIDs = append(permIDs, id)
				}
			}
			s.AssignPermissionsToRole(roleID, permIDs)
		}
	}

	return nil
}
