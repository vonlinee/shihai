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
	poemService           *services.PoemService
	poemAnnotationService *services.PoemAnnotationService
}

func NewPoemHandler(poemService *services.PoemService, annotationService ...*services.PoemAnnotationService) *PoemHandler {
	handler := &PoemHandler{poemService: poemService}
	if len(annotationService) > 0 {
		handler.poemAnnotationService = annotationService[0]
	}
	return handler
}

// GetPoemList 获取诗词列表
// @Summary 查询诗词列表
// @Tags 诗词
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "关键词"
// @Param dynasty query string false "朝代"
// @Param author query string false "作者"
// @Param genre query string false "体裁"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/poems [get]
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
// @Summary 获取诗词详情
// @Tags 诗词
// @Param id path string true "诗词 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/poems/{id} [get]
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
// @Summary 创建诗词
// @Tags 后台诗词
// @Security BearerAuth
// @Param request body dto.PoemCreateRequest true "诗词信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems [post]
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
// @Summary 更新诗词
// @Tags 后台诗词
// @Security BearerAuth
// @Param id path string true "诗词 ID"
// @Param request body dto.PoemUpdateRequest true "诗词信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/{id} [put]
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
// @Summary 删除诗词
// @Tags 后台诗词
// @Security BearerAuth
// @Param id path string true "诗词 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/{id} [delete]
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

// BatchDeletePoems 批量删除诗词
// @Summary 批量删除诗词
// @Tags 后台诗词
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "诗词 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems [delete]
func (h *PoemHandler) BatchDeletePoems(c *gin.Context) {
	var req dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.poemService.BatchDeletePoems(req.IDs); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "poems deleted successfully", nil)
}

// LikePoem 点赞诗词
// @Summary 点赞诗词
// @Tags 诗词
// @Param id path string true "诗词 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/poems/{id}/like [post]
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
// @Summary 随机获取诗词
// @Tags 诗词
// @Param limit query int false "返回数量" default(5)
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/poems/random [get]
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
// @Summary 查询朝代列表
// @Tags 诗词基础数据
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/dynasties [get]
func (h *PoemHandler) GetDynastyList(c *gin.Context) {
	dynasties, err := h.poemService.GetDynastyList()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, dynasties)
}

// GetPoetList 获取诗人列表
// @Summary 查询诗人列表
// @Tags 诗词基础数据
// @Param keyword query string false "关键词"
// @Param dynastyId query string false "朝代 ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/poets [get]
func (h *PoemHandler) GetPoetList(c *gin.Context) {
	keyword := c.Query("keyword")
	dynastyIDParam := c.Query("dynastyId")
	pageParam := c.Query("page")
	pageSizeParam := c.Query("pageSize")
	var dynastyID uint64
	if dynastyIDParam != "" {
		parsedDynastyID, err := strconv.ParseUint(dynastyIDParam, 10, 64)
		if err != nil {
			utils.BadRequest(c, "invalid dynasty id")
			return
		}
		dynastyID = parsedDynastyID
	}
	page, _ := strconv.Atoi(pageParam)
	pageSize, _ := strconv.Atoi(pageSizeParam)
	isPaginated := pageParam != "" || pageSizeParam != ""
	if isPaginated {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
	}

	poets, total, err := h.poemService.GetPoetList(keyword, dynastyID, page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	if isPaginated {
		utils.PageSuccess(c, poets, total, page, pageSize)
		return
	}
	utils.Success(c, poets)
}

// GetGenreList 获取体裁列表
// @Summary 查询体裁列表
// @Tags 诗词基础数据
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/genres [get]
func (h *PoemHandler) GetGenreList(c *gin.Context) {
	genres, err := h.poemService.GetGenreList()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, genres)
}

// GetGenreCategories 获取分层体裁列表
// @Summary 查询分层体裁列表
// @Tags 诗词基础数据
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/genre-categories [get]
func (h *PoemHandler) GetGenreCategories(c *gin.Context) {
	categories, err := h.poemService.GetGenreCategories()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, categories)
}

// GetPoemTypes 查询体裁管理列表
// @Summary 查询体裁管理列表
// @Tags 后台基础数据
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/poem-types [get]
func (h *PoemHandler) GetPoemTypes(c *gin.Context) {
	poemTypes, err := h.poemService.GetPoemTypes()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, poemTypes)
}

// CreatePoemType 创建体裁
// @Summary 创建体裁
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.PoemTypeCreateRequest true "体裁信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poem-types [post]
func (h *PoemHandler) CreatePoemType(c *gin.Context) {
	var req dto.PoemTypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	poemType, err := h.poemService.CreatePoemType(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, poemType)
}

// UpdatePoemType 更新体裁
// @Summary 更新体裁
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "体裁 ID"
// @Param request body dto.PoemTypeUpdateRequest true "体裁信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poem-types/{id} [put]
func (h *PoemHandler) UpdatePoemType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem type id")
		return
	}

	var req dto.PoemTypeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	poemType, err := h.poemService.UpdatePoemType(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, poemType)
}

// DeletePoemType 删除体裁
// @Summary 删除体裁
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "体裁 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poem-types/{id} [delete]
func (h *PoemHandler) DeletePoemType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem type id")
		return
	}

	if err := h.poemService.DeletePoemType(id); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "poem type deleted successfully", nil)
}

// BatchDeletePoemTypes 批量删除体裁
// @Summary 批量删除体裁
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "体裁 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poem-types [delete]
func (h *PoemHandler) BatchDeletePoemTypes(c *gin.Context) {
	var req dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.poemService.BatchDeletePoemTypes(req.IDs); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "poem types deleted successfully", nil)
}

// GetCiTunes 查询词牌管理列表
// @Summary 查询词牌管理列表
// @Tags 后台基础数据
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/ci-tunes [get]
func (h *PoemHandler) GetCiTunes(c *gin.Context) {
	ciTunes, err := h.poemService.GetCiTunes()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, ciTunes)
}

// CreateCiTune 创建词牌
// @Summary 创建词牌
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.CiTuneCreateRequest true "词牌信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/ci-tunes [post]
func (h *PoemHandler) CreateCiTune(c *gin.Context) {
	var req dto.CiTuneCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	ciTune, err := h.poemService.CreateCiTune(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, ciTune)
}

// UpdateCiTune 更新词牌
// @Summary 更新词牌
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "词牌 ID"
// @Param request body dto.CiTuneUpdateRequest true "词牌信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/ci-tunes/{id} [put]
func (h *PoemHandler) UpdateCiTune(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid ci tune id")
		return
	}

	var req dto.CiTuneUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	ciTune, err := h.poemService.UpdateCiTune(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, ciTune)
}

// DeleteCiTune 删除词牌
// @Summary 删除词牌
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "词牌 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/ci-tunes/{id} [delete]
func (h *PoemHandler) DeleteCiTune(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid ci tune id")
		return
	}

	if err := h.poemService.DeleteCiTune(id); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "ci tune deleted successfully", nil)
}

// BatchDeleteCiTunes 批量删除词牌
// @Summary 批量删除词牌
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "词牌 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/ci-tunes [delete]
func (h *PoemHandler) BatchDeleteCiTunes(c *gin.Context) {
	var req dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.poemService.BatchDeleteCiTunes(req.IDs); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "ci tunes deleted successfully", nil)
}

// ParseCiTuneTitle 根据标题解析词牌
// @Summary 根据标题解析词牌
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.CiTuneTitleParseRequest true "标题信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/ci-tunes/parse-title [post]
func (h *PoemHandler) ParseCiTuneTitle(c *gin.Context) {
	var req dto.CiTuneTitleParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.poemService.ParseCiTuneTitle(req.Title)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// CreateDynasty 创建朝代
// @Summary 创建朝代
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.DynastyCreateRequest true "朝代信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/dynasties [post]
func (h *PoemHandler) CreateDynasty(c *gin.Context) {
	var req dto.DynastyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	dynasty, err := h.poemService.CreateDynasty(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, dynasty)
}

// UpdateDynasty 更新朝代
// @Summary 更新朝代
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "朝代 ID"
// @Param request body dto.DynastyUpdateRequest true "朝代信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/dynasties/{id} [put]
func (h *PoemHandler) UpdateDynasty(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid dynasty id")
		return
	}

	var req dto.DynastyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	dynasty, err := h.poemService.UpdateDynasty(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, dynasty)
}

// DeleteDynasty 删除朝代
// @Summary 删除朝代
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "朝代 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/dynasties/{id} [delete]
func (h *PoemHandler) DeleteDynasty(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid dynasty id")
		return
	}

	err = h.poemService.DeleteDynasty(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "dynasty deleted successfully", nil)
}

// BatchDeleteDynasties 批量删除朝代
// @Summary 批量删除朝代
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "朝代 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/dynasties [delete]
func (h *PoemHandler) BatchDeleteDynasties(c *gin.Context) {
	var req dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.poemService.BatchDeleteDynasties(req.IDs); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "dynasties deleted successfully", nil)
}

// CreatePoet 创建诗人
// @Summary 创建诗人
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.PoetCreateRequest true "诗人信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poets [post]
func (h *PoemHandler) CreatePoet(c *gin.Context) {
	var req dto.PoetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	poet, err := h.poemService.CreatePoet(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, poet)
}

// UpdatePoet 更新诗人
// @Summary 更新诗人
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "诗人 ID"
// @Param request body dto.PoetUpdateRequest true "诗人信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poets/{id} [put]
func (h *PoemHandler) UpdatePoet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poet id")
		return
	}

	var req dto.PoetUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	poet, err := h.poemService.UpdatePoet(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, poet)
}

// DeletePoet 删除诗人
// @Summary 删除诗人
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "诗人 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poets/{id} [delete]
func (h *PoemHandler) DeletePoet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poet id")
		return
	}

	err = h.poemService.DeletePoet(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "poet deleted successfully", nil)
}

// BatchDeletePoets 批量删除诗人
// @Summary 批量删除诗人
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "诗人 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poets [delete]
func (h *PoemHandler) BatchDeletePoets(c *gin.Context) {
	var req dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.poemService.BatchDeletePoets(req.IDs); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "poets deleted successfully", nil)
}

// GetPoemAnnotations 查询诗词标注。
// @Summary 查询诗词标注
// @Tags 后台诗词标注
// @Security BearerAuth
// @Param id path string true "诗词 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/{id}/annotations [get]
func (h *PoemHandler) GetPoemAnnotations(c *gin.Context) {
	if h.poemAnnotationService == nil {
		utils.InternalServerError(c, "poem annotation service not configured")
		return
	}
	poemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem id")
		return
	}

	annotations, err := h.poemAnnotationService.ListAnnotations(poemID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.Success(c, annotations)
}

// CreatePoemAnnotation 创建诗词标注。
// @Summary 创建诗词标注
// @Tags 后台诗词标注
// @Security BearerAuth
// @Param id path string true "诗词 ID"
// @Param request body dto.PoemAnnotationCreateRequest true "标注信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/{id}/annotations [post]
func (h *PoemHandler) CreatePoemAnnotation(c *gin.Context) {
	if h.poemAnnotationService == nil {
		utils.InternalServerError(c, "poem annotation service not configured")
		return
	}
	poemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid poem id")
		return
	}

	var req dto.PoemAnnotationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	annotation, err := h.poemAnnotationService.CreateAnnotation(poemID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, annotation)
}

// UpdatePoemAnnotation 更新诗词标注。
// @Summary 更新诗词标注
// @Tags 后台诗词标注
// @Security BearerAuth
// @Param id path string true "标注 ID"
// @Param request body dto.PoemAnnotationUpdateRequest true "标注信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poem-annotations/{id} [put]
func (h *PoemHandler) UpdatePoemAnnotation(c *gin.Context) {
	if h.poemAnnotationService == nil {
		utils.InternalServerError(c, "poem annotation service not configured")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid annotation id")
		return
	}

	var req dto.PoemAnnotationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	annotation, err := h.poemAnnotationService.UpdateAnnotation(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, annotation)
}

// DeletePoemAnnotation 删除诗词标注。
// @Summary 删除诗词标注
// @Tags 后台诗词标注
// @Security BearerAuth
// @Param id path string true "标注 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poem-annotations/{id} [delete]
func (h *PoemHandler) DeletePoemAnnotation(c *gin.Context) {
	if h.poemAnnotationService == nil {
		utils.InternalServerError(c, "poem annotation service not configured")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid annotation id")
		return
	}

	if err := h.poemAnnotationService.DeleteAnnotation(id); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "poem annotation deleted successfully", nil)
}
