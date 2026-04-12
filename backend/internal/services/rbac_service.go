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
	rolePermissionRepo *repository.RolePermissionRepository
	userRoleRepo       *repository.UserRoleRepository
}

func NewRBACService(
	roleRepo *repository.RoleRepository,
	rolePermissionRepo *repository.RolePermissionRepository,
	userRoleRepo *repository.UserRoleRepository,
) *RBACService {
	return &RBACService{
		roleRepo:           roleRepo,
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
	// 获取角色权限 (字符串数组)
	permissions, err := s.rolePermissionRepo.GetCodesByRoleID(id)
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

// ==================== Role-Permission Service ====================

// AssignPermissionsToRole 为角色分配权限 (权限编码字符串数组)
func (s *RBACService) AssignPermissionsToRole(roleID uint64, permissionCodes []string) error {
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
	for _, code := range permissionCodes {
		rp := &models.RolePermission{
			RoleID: roleID,
			Code:   code,
		}
		if err := s.rolePermissionRepo.Create(rp); err != nil {
			return err
		}
	}

	return nil
}

// GetRolePermissions 获取角色的权限列表(仅仅是编码)
func (s *RBACService) GetRolePermissions(roleID uint64) ([]string, error) {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	codes, err := s.rolePermissionRepo.GetCodesByRoleID(roleID)
	if err != nil {
		return nil, err
	}

	return codes, nil
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

// ==================== Initialization ====================

// InitDefaultRolesAndPermissions 初始化默认角色和权限
func (s *RBACService) InitDefaultRolesAndPermissions() error {
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

		// 分配权限给角色 (通过字符串编码)
		if permCodes, ok := models.RolePermissionMap[role.Name]; ok {
			s.AssignPermissionsToRole(roleID, permCodes)
		}
	}

	return nil
}
