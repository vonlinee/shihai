package handlers

import (
	"net/http"
	"shihai/internal/dto"
	"shihai/internal/services"
	"shihai/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 认证
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, user)
}

// AdminCreateUser 管理员创建用户
// @Summary 创建用户
// @Tags 后台用户
// @Security BearerAuth
// @Param request body dto.AdminCreateUserRequest true "用户信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/users [post]
func (h *UserHandler) AdminCreateUser(c *gin.Context) {
	var req dto.AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.AdminCreateUser(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, user)
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户资料
// @Tags 用户
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "")
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint64))
	if err != nil {
		utils.NotFound(c, "")
		return
	}

	utils.Success(c, user)
}

// UpdateProfile 更新当前用户信息
// @Summary 更新当前用户资料
// @Tags 用户
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "用户资料"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(userID.(uint64), &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, user)
}

// ChangePassword 修改密码
// @Summary 修改当前用户密码
// @Tags 用户
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "密码信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.userService.ChangePassword(userID.(uint64), &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "password changed successfully", nil)
}

// GetUserList 获取用户列表（管理员）
// @Summary 查询用户列表
// @Tags 后台用户
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "关键词"
// @Param role query string false "角色"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/users [get]
func (h *UserHandler) GetUserList(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	users, total, err := h.userService.GetUserList(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageSuccess(c, users, total, req.Page, req.PageSize)
}

// GetUserByID 根据ID获取用户
// @Summary 获取用户详情
// @Tags 后台用户
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		utils.NotFound(c, "")
		return
	}

	utils.Success(c, user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 后台用户
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "user deleted successfully", nil)
}
