package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

// RoleRepository 角色仓库
type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create 创建角色
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// GetByID 根据ID获取角色
func (r *RoleRepository) GetByID(id uint64) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName 根据名称获取角色
func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *RoleRepository) List(page, pageSize int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Update 更新角色
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// PermissionRepository 权限仓库
type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// Create 创建权限
func (r *PermissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// GetByID 根据ID获取权限
func (r *PermissionRepository) GetByID(id uint64) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByCode 根据编码获取权限
func (r *PermissionRepository) GetByCode(code string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// List 获取权限列表
func (r *PermissionRepository) List(page, pageSize int, module string) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	query := r.db.Model(&models.Permission{})
	if module != "" {
		query = query.Where("module = ?", module)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// ListByModule 按模块获取权限
func (r *PermissionRepository) ListByModule(module string) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.Where("module = ?", module).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// ListAll 获取所有权限
func (r *PermissionRepository) ListAll() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// Update 更新权限
func (r *PermissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

// Delete 删除权限
func (r *PermissionRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

// RolePermissionRepository 角色权限关联仓库
type RolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

// Create 创建角色权限关联
func (r *RolePermissionRepository) Create(rp *models.RolePermission) error {
	return r.db.Create(rp).Error
}

// Delete 删除角色权限关联
func (r *RolePermissionRepository) Delete(roleID, permissionID uint64) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&models.RolePermission{}).Error
}

// DeleteByRoleID 删除角色的所有权限
func (r *RolePermissionRepository) DeleteByRoleID(roleID uint64) error {
	return r.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error
}

// ListByRoleID 获取角色的权限列表
func (r *RolePermissionRepository) ListByRoleID(roleID uint64) ([]models.RolePermission, error) {
	var rps []models.RolePermission
	if err := r.db.Where("role_id = ?", roleID).Preload("Permission").Find(&rps).Error; err != nil {
		return nil, err
	}
	return rps, nil
}

// GetPermissionsByRoleID 获取角色的权限编码列表
func (r *RolePermissionRepository) GetPermissionsByRoleID(roleID uint64) ([]string, error) {
	var codes []string
	if err := r.db.Model(&models.RolePermission{}).
		Joins("JOIN permission ON permission.id = role_permission.permission_id").
		Where("role_permission.role_id = ?", roleID).
		Pluck("permission.code", &codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

// UserRoleRepository 用户角色关联仓库
type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

// Create 创建用户角色关联
func (r *UserRoleRepository) Create(ur *models.UserRole) error {
	return r.db.Create(ur).Error
}

// Delete 删除用户角色关联
func (r *UserRoleRepository) Delete(userID, roleID uint64) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

// DeleteByUserID 删除用户的所有角色
func (r *UserRoleRepository) DeleteByUserID(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error
}

// ListByUserID 获取用户的角色列表
func (r *UserRoleRepository) ListByUserID(userID uint64) ([]models.UserRole, error) {
	var urs []models.UserRole
	if err := r.db.Where("user_id = ?", userID).Preload("Role").Find(&urs).Error; err != nil {
		return nil, err
	}
	return urs, nil
}

// GetRolesByUserID 获取用户的角色名称列表
func (r *UserRoleRepository) GetRolesByUserID(userID uint64) ([]string, error) {
	var names []string
	if err := r.db.Model(&models.UserRole{}).
		Joins("JOIN role ON role.id = user_role.role_id").
		Where("user_role.user_id = ?", userID).
		Pluck("role.name", &names).Error; err != nil {
		return nil, err
	}
	return names, nil
}

// GetRoleIDsByUserID 获取用户的角色ID列表
func (r *UserRoleRepository) GetRoleIDsByUserID(userID uint64) ([]uint64, error) {
	var ids []uint64
	if err := r.db.Model(&models.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// CheckUserHasRole 检查用户是否有指定角色
func (r *UserRoleRepository) CheckUserHasRole(userID, roleID uint64) (bool, error) {
	var count int64
	if err := r.db.Model(&models.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserPermissions 获取用户的所有权限编码
func (r *UserRoleRepository) GetUserPermissions(userID uint64) ([]string, error) {
	var codes []string
	if err := r.db.Model(&models.UserRole{}).
		Joins("JOIN role_permission ON role_permission.role_id = user_role.role_id").
		Joins("JOIN permission ON permission.id = role_permission.permission_id").
		Where("user_role.user_id = ?", userID).
		Distinct().
		Pluck("permission.code", &codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}
