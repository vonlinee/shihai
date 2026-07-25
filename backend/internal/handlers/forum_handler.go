package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"shihai/internal/dto"
	"shihai/internal/services"
	"shihai/pkg/utils"

	"github.com/gin-gonic/gin"
)

type forumService interface {
	ListPosts(page, pageSize int, keyword string, includeDeleted bool) ([]dto.ForumPostResponse, int64, error)
	GetPostByID(id uint64, incrementViews bool) (*dto.ForumPostResponse, error)
	CreatePost(userID uint64, req *dto.ForumPostCreateRequest) (*dto.ForumPostResponse, error)
	UpdatePost(id uint64, userID uint64, canModerate bool, req *dto.ForumPostUpdateRequest) (*dto.ForumPostResponse, error)
	DeletePost(id uint64, userID uint64, canModerate bool) error
	SetPostPinned(id uint64, isPinned bool) error
	ListReplies(postID uint64, page, pageSize int) ([]dto.ForumReplyResponse, int64, error)
	CreateReply(postID uint64, userID uint64, req *dto.ForumReplyCreateRequest) (*dto.ForumReplyResponse, error)
	DeleteReply(id uint64, userID uint64, canModerate bool) error
}

type ForumHandler struct {
	forumService forumService
}

func NewForumHandler(forumService forumService) *ForumHandler {
	return &ForumHandler{forumService: forumService}
}

// ListPosts 查询论坛帖子列表。
// @Summary 查询论坛帖子列表
// @Tags 论坛
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param keyword query string false "关键词"
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts [get]
func (h *ForumHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	posts, total, err := h.forumService.ListPosts(page, pageSize, keyword, false)
	if err != nil {
		utils.InternalServerError(c, "failed to list forum posts")
		return
	}

	utils.PageSuccess(c, posts, total, page, pageSize)
}

// ListAllPosts 查询后台论坛帖子列表。
// @Summary 查询后台论坛帖子列表
// @Tags 后台论坛
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param keyword query string false "关键词"
// @Param includeDeleted query bool false "是否包含已删除帖子"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/forum/posts [get]
func (h *ForumHandler) ListAllPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	includeDeleted := c.Query("includeDeleted") == "true"

	posts, total, err := h.forumService.ListPosts(page, pageSize, keyword, includeDeleted)
	if err != nil {
		utils.InternalServerError(c, "failed to list forum posts")
		return
	}

	utils.PageSuccess(c, posts, total, page, pageSize)
}

// GetPostByID 获取论坛帖子详情。
// @Summary 获取论坛帖子详情
// @Tags 论坛
// @Param id path string true "帖子 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/posts/{id} [get]
func (h *ForumHandler) GetPostByID(c *gin.Context) {
	id, ok := parseUintParam(c, "id", "invalid forum post id")
	if !ok {
		return
	}

	post, err := h.forumService.GetPostByID(id, true)
	if err != nil {
		respondForumError(c, err)
		return
	}

	utils.Success(c, post)
}

// CreatePost 创建论坛帖子。
// @Summary 创建论坛帖子
// @Tags 论坛
// @Security BearerAuth
// @Param request body dto.ForumPostCreateRequest true "帖子信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /api/forum/posts [post]
func (h *ForumHandler) CreatePost(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.Unauthorized(c, "")
		return
	}

	var req dto.ForumPostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	post, err := h.forumService.CreatePost(userID, &req)
	if err != nil {
		respondForumError(c, err)
		return
	}

	utils.Success(c, post)
}

// UpdatePost 更新自己的论坛帖子。
// @Summary 更新自己的论坛帖子
// @Tags 论坛
// @Security BearerAuth
// @Param id path string true "帖子 ID"
// @Param request body dto.ForumPostUpdateRequest true "帖子信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/posts/{id} [put]
func (h *ForumHandler) UpdatePost(c *gin.Context) {
	id, ok := parseUintParam(c, "id", "invalid forum post id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		utils.Unauthorized(c, "")
		return
	}

	var req dto.ForumPostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	post, err := h.forumService.UpdatePost(id, userID, false, &req)
	if err != nil {
		respondForumError(c, err)
		return
	}

	utils.Success(c, post)
}

// DeletePost 删除自己的论坛帖子。
// @Summary 删除自己的论坛帖子
// @Tags 论坛
// @Security BearerAuth
// @Param id path string true "帖子 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/posts/{id} [delete]
func (h *ForumHandler) DeletePost(c *gin.Context) {
	h.deletePost(c, false)
}

// AdminDeletePost 后台删除论坛帖子。
// @Summary 后台删除论坛帖子
// @Tags 后台论坛
// @Security BearerAuth
// @Param id path string true "帖子 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/admin/forum/posts/{id} [delete]
func (h *ForumHandler) AdminDeletePost(c *gin.Context) {
	h.deletePost(c, true)
}

// SetPostPinned 设置论坛帖子置顶状态。
// @Summary 设置论坛帖子置顶状态
// @Tags 后台论坛
// @Security BearerAuth
// @Param id path string true "帖子 ID"
// @Param request body dto.ForumPostPinRequest true "置顶状态"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/admin/forum/posts/{id}/pin [put]
func (h *ForumHandler) SetPostPinned(c *gin.Context) {
	id, ok := parseUintParam(c, "id", "invalid forum post id")
	if !ok {
		return
	}

	var req dto.ForumPostPinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.forumService.SetPostPinned(id, req.IsPinned); err != nil {
		respondForumError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "forum post pin state updated", nil)
}

// ListReplies 查询论坛帖子回复。
// @Summary 查询论坛帖子回复
// @Tags 论坛
// @Param id path string true "帖子 ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(50)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/posts/{id}/replies [get]
func (h *ForumHandler) ListReplies(c *gin.Context) {
	postID, ok := parseUintParam(c, "id", "invalid forum post id")
	if !ok {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	replies, total, err := h.forumService.ListReplies(postID, page, pageSize)
	if err != nil {
		respondForumError(c, err)
		return
	}

	utils.PageSuccess(c, replies, total, page, pageSize)
}

// CreateReply 创建论坛回复。
// @Summary 创建论坛回复
// @Tags 论坛
// @Security BearerAuth
// @Param id path string true "帖子 ID"
// @Param request body dto.ForumReplyCreateRequest true "回复信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /api/forum/posts/{id}/replies [post]
func (h *ForumHandler) CreateReply(c *gin.Context) {
	postID, ok := parseUintParam(c, "id", "invalid forum post id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		utils.Unauthorized(c, "")
		return
	}

	var req dto.ForumReplyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	reply, err := h.forumService.CreateReply(postID, userID, &req)
	if err != nil {
		respondForumError(c, err)
		return
	}

	utils.Success(c, reply)
}

// DeleteReply 删除自己的论坛回复。
// @Summary 删除自己的论坛回复
// @Tags 论坛
// @Security BearerAuth
// @Param id path string true "回复 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/replies/{id} [delete]
func (h *ForumHandler) DeleteReply(c *gin.Context) {
	h.deleteReply(c, false)
}

// AdminDeleteReply 后台删除论坛回复。
// @Summary 后台删除论坛回复
// @Tags 后台论坛
// @Security BearerAuth
// @Param id path string true "回复 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/admin/forum/replies/{id} [delete]
func (h *ForumHandler) AdminDeleteReply(c *gin.Context) {
	h.deleteReply(c, true)
}

func (h *ForumHandler) deletePost(c *gin.Context, canModerate bool) {
	id, ok := parseUintParam(c, "id", "invalid forum post id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		utils.Unauthorized(c, "")
		return
	}

	if err := h.forumService.DeletePost(id, userID, canModerate); err != nil {
		respondForumError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "forum post deleted successfully", nil)
}

func (h *ForumHandler) deleteReply(c *gin.Context, canModerate bool) {
	id, ok := parseUintParam(c, "id", "invalid forum reply id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		utils.Unauthorized(c, "")
		return
	}

	if err := h.forumService.DeleteReply(id, userID, canModerate); err != nil {
		respondForumError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "forum reply deleted successfully", nil)
}

func currentUserID(c *gin.Context) (uint64, bool) {
	value, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	userID, ok := value.(uint64)
	return userID, ok
}

func parseUintParam(c *gin.Context, name string, message string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		utils.BadRequest(c, message)
		return 0, false
	}
	return id, true
}

func respondForumError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrForumNotFound):
		utils.NotFound(c, "forum resource not found")
	case errors.Is(err, services.ErrForumPermissionDenied):
		utils.Forbidden(c, "permission denied")
	case errors.Is(err, services.ErrForumInvalidRequest):
		utils.Error(c, http.StatusUnprocessableEntity, "invalid forum request")
	default:
		utils.InternalServerError(c, "forum operation failed")
	}
}
