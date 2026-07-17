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
