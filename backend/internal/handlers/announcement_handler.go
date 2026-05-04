package handlers

import (
	"net/http"
	"shihai/internal/services"
	"shihai/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	announcementService *services.AnnouncementService
}

func NewAnnouncementHandler(announcementService *services.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{announcementService: announcementService}
}

// GetAnnouncements 获取公告列表
func (h *AnnouncementHandler) GetAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	onlyPinned := c.Query("pinned") == "true"

	announcements, total, err := h.announcementService.GetAnnouncements(page, pageSize, onlyPinned)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageSuccess(c, announcements, total, page, pageSize)
}

// GetAnnouncementByID 根据ID获取公告
func (h *AnnouncementHandler) GetAnnouncementByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid announcement id")
		return
	}

	announcement, err := h.announcementService.GetAnnouncementByID(id)
	if err != nil {
		utils.NotFound(c, "")
		return
	}

	utils.Success(c, announcement)
}

// CreateAnnouncement 创建公告
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		IsPinned bool   `json:"isPinned"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	announcement, err := h.announcementService.CreateAnnouncement(req.Title, req.Content, req.IsPinned)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, announcement)
}

// UpdateAnnouncement 更新公告
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid announcement id")
		return
	}

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		IsPinned bool   `json:"isPinned"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	announcement, err := h.announcementService.UpdateAnnouncement(id, req.Title, req.Content, req.IsPinned)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, announcement)
}

// DeleteAnnouncement 删除公告
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid announcement id")
		return
	}

	err = h.announcementService.DeleteAnnouncement(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "announcement deleted successfully", nil)
}
