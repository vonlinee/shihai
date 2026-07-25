package handlers

import (
	"net/http"

	"shihai/internal/dto"
	"shihai/pkg/utils"

	"github.com/gin-gonic/gin"
)

type textConversionService interface {
	ConvertTexts(req *dto.TextConversionRequest) (*dto.TextConversionResponse, error)
}

// TextConversionHandler 处理中文简繁转换请求。
type TextConversionHandler struct {
	textConversionService textConversionService
}

// NewTextConversionHandler 创建中文简繁转换处理器。
func NewTextConversionHandler(textConversionService textConversionService) *TextConversionHandler {
	return &TextConversionHandler{textConversionService: textConversionService}
}

// ConvertTexts 批量转换文本，并按输入顺序返回转换结果。
// @Summary 批量简繁转换
// @Tags 后台文本转换
// @Security BearerAuth
// @Param request body dto.TextConversionRequest true "转换信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/text-conversion [post]
func (h *TextConversionHandler) ConvertTexts(c *gin.Context) {
	var req dto.TextConversionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.textConversionService.ConvertTexts(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, resp)
}
