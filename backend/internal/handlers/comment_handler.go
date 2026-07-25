package handlers

import (
	"net/http"
	"shihai/internal/dto"
	"shihai/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type commentService interface {
	GetCommentsByPoem(poemID uint64, page, pageSize int) ([]dto.CommentResponse, int64, error)
	GetAllComments(page, pageSize int) ([]dto.CommentResponse, int64, error)
	CreateComment(userID *uint64, req *dto.CommentCreateRequest) (*dto.CommentResponse, error)
	DeleteComment(id uint64, userID uint64) error
	VoteComment(userID *uint64, visitorID string, req *dto.CommentVoteRequest) error
}

type CommentHandler struct {
	commentService commentService
}

func NewCommentHandler(commentService commentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// GetComments 获取诗词的评论列表
// @Summary 查询诗词评论
// @Tags 评论
// @Param poemId query string true "诗词 ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	poemID, err := strconv.ParseUint(c.Query("poemId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	comments, total, err := h.commentService.GetCommentsByPoem(poemID, page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageSuccess(c, comments, total, page, pageSize)
}

// GetAllComments 查询全部评论。
// @Summary 查询全部评论
// @Tags 后台评论
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/comments/all [get]
func (h *CommentHandler) GetAllComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	comments, total, err := h.commentService.GetAllComments(page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageSuccess(c, comments, total, page, pageSize)
}

// CreateComment 创建评论
// @Summary 创建评论
// @Tags 评论
// @Security BearerAuth
// @Param request body dto.CommentCreateRequest true "评论信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req dto.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 获取当前用户ID（如果已登录）
	var userID *uint64
	if id, exists := c.Get("userID"); exists {
		uid := id.(uint64)
		userID = &uid
	}

	comment, err := h.commentService.CreateComment(userID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, comment)
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Tags 评论
// @Security BearerAuth
// @Param id path string true "评论 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid comment id")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "")
		return
	}

	err = h.commentService.DeleteComment(id, userID.(uint64))
	if err != nil {
		utils.Error(c, http.StatusForbidden, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "comment deleted successfully", nil)
}

// VoteComment 评论投票
// @Summary 评论投票
// @Tags 评论
// @Param X-Visitor-ID header string false "游客标识"
// @Param request body dto.CommentVoteRequest true "投票信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/comments/vote [post]
func (h *CommentHandler) VoteComment(c *gin.Context) {
	var req dto.CommentVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 获取用户ID或游客ID
	var userID *uint64
	var visitorID string
	if id, exists := c.Get("userID"); exists {
		uid := id.(uint64)
		userID = &uid
	} else {
		// 获取游客ID（从请求头或生成）
		visitorID = c.GetHeader("X-Visitor-ID")
		if visitorID == "" {
			visitorID = c.ClientIP()
		}
	}

	err := h.commentService.VoteComment(userID, visitorID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "voted successfully", nil)
}
