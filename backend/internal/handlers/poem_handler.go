package handlers

import (
	"net/http"
	"shihai/internal/dto"
	"shihai/internal/services"
	"shihai/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PoemHandler struct {
	poemService *services.PoemService
}

func NewPoemHandler(poemService *services.PoemService) *PoemHandler {
	return &PoemHandler{poemService: poemService}
}

// GetPoemList 获取诗词列表
func (h *PoemHandler) GetPoemList(c *gin.Context) {
	var req dto.PoemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	poems, total, err := h.poemService.GetPoemList(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageSuccess(c, poems, total, req.Page, req.PageSize)
}

// GetPoemByID 根据ID获取诗词
func (h *PoemHandler) GetPoemByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem id")
		return
	}

	poem, err := h.poemService.GetPoemByID(id)
	if err != nil {
		utils.NotFound(c, "")
		return
	}

	utils.Success(c, poem)
}

// CreatePoem 创建诗词
func (h *PoemHandler) CreatePoem(c *gin.Context) {
	var req dto.PoemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	poem, err := h.poemService.CreatePoem(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, poem)
}

// UpdatePoem 更新诗词
func (h *PoemHandler) UpdatePoem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem id")
		return
	}

	var req dto.PoemUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	poem, err := h.poemService.UpdatePoem(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, poem)
}

// DeletePoem 删除诗词
func (h *PoemHandler) DeletePoem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem id")
		return
	}

	err = h.poemService.DeletePoem(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "poem deleted successfully", nil)
}

// LikePoem 点赞诗词
func (h *PoemHandler) LikePoem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem id")
		return
	}

	err = h.poemService.LikePoem(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "liked successfully", nil)
}

// GetRandomPoems 随机获取诗词
func (h *PoemHandler) GetRandomPoems(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	poems, err := h.poemService.GetRandomPoems(limit)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, poems)
}

// GetDynastyList 获取朝代列表
func (h *PoemHandler) GetDynastyList(c *gin.Context) {
	dynasties, err := h.poemService.GetDynastyList()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, dynasties)
}
