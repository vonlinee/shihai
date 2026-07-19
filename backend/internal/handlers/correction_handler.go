package handlers

import (
	"shihai/internal/dto"
	"shihai/pkg/utils"

	"github.com/gin-gonic/gin"
)

type correctionService interface {
	ListCorrections(req dto.CorrectionListRequest) ([]dto.CorrectionResponse, int64, error)
}

// CorrectionHandler adapts correction management APIs to HTTP.
type CorrectionHandler struct {
	correctionService correctionService
}

// NewCorrectionHandler creates a CorrectionHandler.
func NewCorrectionHandler(correctionService correctionService) *CorrectionHandler {
	return &CorrectionHandler{correctionService: correctionService}
}

// ListCorrections handles GET /api/admin/corrections.
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
