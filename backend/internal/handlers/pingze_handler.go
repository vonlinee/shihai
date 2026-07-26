package handlers

import (
	"net/http"

	"shihai/internal/dto"
	"shihai/pkg/utils"

	"github.com/gin-gonic/gin"
)

type pingzeRecognitionService interface {
	RecognizePingze(req *dto.PingzeRecognitionRequest) (*dto.PingzeRecognitionResponse, error)
}

// PingzeHandler 处理诗词平仄自动识别请求。
type PingzeHandler struct {
	pingzeRecognitionService pingzeRecognitionService
}

// NewPingzeHandler 创建诗词平仄自动识别处理器。
func NewPingzeHandler(pingzeRecognitionService pingzeRecognitionService) *PingzeHandler {
	return &PingzeHandler{pingzeRecognitionService: pingzeRecognitionService}
}

// RecognizePingze 批量识别文本平仄，并按输入顺序返回结果。
// @Summary 批量识别平仄
// @Tags 后台诗词
// @Security BearerAuth
// @Param request body dto.PingzeRecognitionRequest true "平仄识别信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/pingze-recognition [post]
func (h *PingzeHandler) RecognizePingze(c *gin.Context) {
	var req dto.PingzeRecognitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.pingzeRecognitionService.RecognizePingze(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, resp)
}
