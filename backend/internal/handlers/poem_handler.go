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

// GetPoetList 获取诗人列表
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
func (h *PoemHandler) GetGenreList(c *gin.Context) {
	genres, err := h.poemService.GetGenreList()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, genres)
}

// CreateDynasty 创建朝代
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

// CreatePoet 创建诗人
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
