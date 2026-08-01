package handlers

import (
	"strconv"

	"shihai/internal/dto"
	"shihai/pkg/utils"

	"github.com/gin-gonic/gin"
)

type correctionService interface {
	CreateCorrection(userID uint64, req dto.CorrectionCreateRequest) (*dto.CorrectionResponse, error)
	GetCorrection(id uint64) (*dto.CorrectionResponse, error)
	ListCorrections(req dto.CorrectionListRequest) ([]dto.CorrectionResponse, int64, error)
	ListUserCorrections(userID uint64, req dto.CorrectionListRequest) ([]dto.CorrectionResponse, int64, error)
	UpdateCorrectionStatus(id uint64, req dto.CorrectionStatusUpdateRequest) (*dto.CorrectionResponse, error)
}

// CorrectionHandler adapts correction management APIs to HTTP.
type CorrectionHandler struct {
	correctionService correctionService
}

// NewCorrectionHandler creates a CorrectionHandler.
func NewCorrectionHandler(correctionService correctionService) *CorrectionHandler {
	return &CorrectionHandler{correctionService: correctionService}
}

// CreateCorrection handles POST /api/corrections.
// @Summary 提交诗词纠错申请
// @Tags 纠错
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CorrectionCreateRequest true "纠错申请"
// @Success 200 {object} utils.Response{data=dto.CorrectionResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/corrections [post]
func (h *CorrectionHandler) CreateCorrection(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "请先登录")
		return
	}
	userID, ok := userIDValue.(uint64)
	if !ok || userID == 0 {
		utils.Unauthorized(c, "登录状态无效")
		return
	}

	var req dto.CorrectionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	correction, err := h.correctionService.CreateCorrection(userID, req)
	if err != nil {
		utils.InternalServerError(c, "提交纠错申请失败")
		return
	}
	utils.Success(c, correction)
}

// ListCorrections handles GET /api/admin/corrections.
// @Summary 查询纠错列表
// @Tags 后台纠错
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "关键词"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/corrections [get]
func (h *CorrectionHandler) ListCorrections(c *gin.Context) {
	var req dto.CorrectionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	corrections, total, err := h.correctionService.ListCorrections(req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	utils.PageSuccess(c, corrections, total, page, pageSize)
}

// ListMyCorrections handles GET /api/corrections/my.
// @Summary 查询我的纠错申请
// @Tags 纠错
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/corrections/my [get]
func (h *CorrectionHandler) ListMyCorrections(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "请先登录")
		return
	}
	userID, ok := userIDValue.(uint64)
	if !ok || userID == 0 {
		utils.Unauthorized(c, "登录状态无效")
		return
	}

	var req dto.CorrectionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	corrections, total, err := h.correctionService.ListUserCorrections(userID, req)
	if err != nil {
		utils.InternalServerError(c, "查询我的纠错申请失败")
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	utils.PageSuccess(c, corrections, total, page, pageSize)
}

// GetCorrection handles GET /api/admin/corrections/{id}.
// @Summary 查询纠错详情
// @Tags 后台纠错
// @Security BearerAuth
// @Param id path string true "纠错申请 ID"
// @Success 200 {object} utils.Response{data=dto.CorrectionResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/corrections/{id} [get]
func (h *CorrectionHandler) GetCorrection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid correction id")
		return
	}

	correction, err := h.correctionService.GetCorrection(id)
	if err != nil {
		utils.InternalServerError(c, "查询纠错详情失败")
		return
	}
	utils.Success(c, correction)
}

// UpdateCorrectionStatus handles PUT /api/admin/corrections/{id}/status.
// @Summary 更新纠错状态
// @Tags 后台纠错
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "纠错申请 ID"
// @Param request body dto.CorrectionStatusUpdateRequest true "纠错状态"
// @Success 200 {object} utils.Response{data=dto.CorrectionResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/corrections/{id}/status [put]
func (h *CorrectionHandler) UpdateCorrectionStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid correction id")
		return
	}

	var req dto.CorrectionStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	correction, err := h.correctionService.UpdateCorrectionStatus(id, req)
	if err != nil {
		utils.InternalServerError(c, "更新纠错状态失败")
		return
	}
	utils.Success(c, correction)
}
