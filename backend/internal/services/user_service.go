package services

import (
	"errors"
	"shihai/internal/dto"
	"shihai/internal/models"
	"shihai/internal/repository"
	"shihai/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户，如果未提供姓名，则使用用户名作为默认姓名
	name := req.Name
	if name == "" {
		name = req.Username
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     name,
		Role:     "user",
		IsActive: true,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// Login 用户登录
func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 检查用户状态
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *s.toUserResponse(user),
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint64) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return s.toUserResponse(user), nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint64, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 更新字段
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Age > 0 {
		user.Age = req.Age
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint64, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(id, string(hashedPassword))
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(req *dto.UserListRequest) ([]dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(req.Page, req.PageSize, req.Keyword, req.Role)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, *s.toUserResponse(&user))
	}

	return responses, total, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint64) error {
	return s.userRepo.Delete(id)
}

// toUserResponse 转换为响应格式
func (s *UserService) toUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Age:       user.Age,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
}
