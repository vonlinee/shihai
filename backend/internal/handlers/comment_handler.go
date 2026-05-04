package handlers

import (
	"net/http"
	"shihai/internal/dto"
	"shihai/internal/services"
	"shihai/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// GetComments 获取诗词的评论列表
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

// CreateComment 创建评论
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
