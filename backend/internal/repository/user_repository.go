package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint64) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List 获取用户列表
func (r *UserRepository) List(page, pageSize int, keyword, role string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if role != "" {
		query = query.Joins("JOIN user_role ON user_role.user_id = \"user\".id").
			Joins("JOIN role ON role.id = user_role.role_id").
			Where("role.name = ?", role)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdatePassword 更新密码
func (r *UserRepository) UpdatePassword(id uint64, password string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password", password).Error
}
