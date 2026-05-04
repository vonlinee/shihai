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
	permRepo           *repository.PermissionRepository
}

func NewRBACService(
	roleRepo *repository.RoleRepository,
	rolePermissionRepo *repository.RolePermissionRepository,
	userRoleRepo *repository.UserRoleRepository,
	permRepo *repository.PermissionRepository,
) *RBACService {
	return &RBACService{
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
		userRoleRepo:       userRoleRepo,
		permRepo:           permRepo,
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
//
// 调用时机：服务器启动时调用一次（cmd/server/main.go）。
//
// 执行流程：
//   1) 从 models.AllPermissions （permission_codes.go 中的全局唯一来源）逐条 UPSERT
//      到 permission 表，确保所有接口权限点都存在。
//   2) 从数据库读取全量权限编码，用于 admin 角色的自动全量同步。
//   3) 逐个处理 DefaultRoles：
//      - admin：直接获取全量权限，新增的接口权限会自动同步。
//      - editor / reviewer / user：从 RolePermissionMap 中读取定义。
//
// 幂等性：本函数可以重复调用，不会产生重复数据。
func (s *RBACService) InitDefaultRolesAndPermissions() error {
	// 1. 初始化权限表：以 AllPermissions 为唯一来源，逐条 UPSERT
	for _, perm := range models.AllPermissions {
		existing, err := s.permRepo.GetByCode(perm.Code)
		if err != nil {
			// 权限不存在，创建
			newPerm := &models.Permission{
				Code:        perm.Code,
				Name:        perm.Name,
				Description: perm.Description,
				Module:      perm.Module,
				IsActive:    true,
			}
			_ = s.permRepo.Create(newPerm)
		} else {
			// 已存在则更新展示元数据（Name/Module/Description）
			existing.Name = perm.Name
			existing.Description = perm.Description
			existing.Module = perm.Module
			_ = s.permRepo.Update(existing)
		}
	}

	// 2. 读取权限表全量编码，用于 admin 角色的自动全量同步
	// 备注：这里从数据库读取而非直接使用 AllPermissions，以适配运营人员
	// 通过后台手动新增的权限点（不在代码中的）也能同步给 admin。
	allPerms, _ := s.permRepo.ListAll()
	allPermCodes := make([]string, 0, len(allPerms))
	for _, p := range allPerms {
		allPermCodes = append(allPermCodes, p.Code)
	}

	// 3. 创建默认角色并分配权限
	for _, role := range models.DefaultRoles {
		existing, err := s.roleRepo.GetByName(role.Name)
		var roleID uint64
		if err != nil {
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

		// admin 角色自动获取全部权限，未来新增的权限点会自动同步给 admin
		if role.Name == models.RoleAdmin {
			_ = s.AssignPermissionsToRole(roleID, allPermCodes)
			continue
		}

		// 非 admin 角色从 RolePermissionMap 中读取预设权限列表
		if permCodes, ok := models.RolePermissionMap[role.Name]; ok {
			_ = s.AssignPermissionsToRole(roleID, permCodes)
		}
	}

	return nil
}

// ==================== Permission Service ====================

// GetPermissionList 获取权限列表
func (s *RBACService) GetPermissionList(req *dto.PermissionListRequest) ([]dto.PermissionResponse, int64, error) {
	perms, total, err := s.permRepo.List(req.Page, req.PageSize, req.Module)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.PermissionResponse
	for _, p := range perms {
		responses = append(responses, s.toPermissionResponse(&p))
	}
	return responses, total, nil
}

// GetAllPermissions 获取所有权限
func (s *RBACService) GetAllPermissions() ([]dto.PermissionResponse, error) {
	perms, err := s.permRepo.ListAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.PermissionResponse
	for _, p := range perms {
		responses = append(responses, s.toPermissionResponse(&p))
	}
	return responses, nil
}

// GetPermissionByID 根据ID获取权限
func (s *RBACService) GetPermissionByID(id uint64) (*dto.PermissionResponse, error) {
	perm, err := s.permRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("permission not found")
	}
	resp := s.toPermissionResponse(perm)
	return &resp, nil
}

// CreatePermission 创建权限
func (s *RBACService) CreatePermission(req *dto.PermissionCreateRequest) (*dto.PermissionResponse, error) {
	// 检查编码是否已存在
	_, err := s.permRepo.GetByCode(req.Code)
	if err == nil {
		return nil, errors.New("permission code already exists")
	}

	perm := &models.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Module:      req.Module,
		IsActive:    true,
	}
	if err := s.permRepo.Create(perm); err != nil {
		return nil, err
	}

	resp := s.toPermissionResponse(perm)
	return &resp, nil
}

// UpdatePermission 更新权限
func (s *RBACService) UpdatePermission(id uint64, req *dto.PermissionUpdateRequest) (*dto.PermissionResponse, error) {
	perm, err := s.permRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("permission not found")
	}

	if req.Name != "" {
		perm.Name = req.Name
	}
	if req.Description != "" {
		perm.Description = req.Description
	}
	if req.Module != "" {
		perm.Module = req.Module
	}
	if req.IsActive != nil {
		perm.IsActive = *req.IsActive
	}

	if err := s.permRepo.Update(perm); err != nil {
		return nil, err
	}

	resp := s.toPermissionResponse(perm)
	return &resp, nil
}

// DeletePermission 删除权限
func (s *RBACService) DeletePermission(id uint64) error {
	return s.permRepo.Delete(id)
}

// toPermissionResponse 转换为权限响应
func (s *RBACService) toPermissionResponse(p *models.Permission) dto.PermissionResponse {
	return dto.PermissionResponse{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		Module:      p.Module,
		IsActive:    p.IsActive,
	}
}
