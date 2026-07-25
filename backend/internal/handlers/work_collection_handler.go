package handlers

import (
	"net/http"
	"shihai/internal/dto"
	"shihai/internal/services"
	"shihai/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkCollectionHandler struct {
	workCollectionService *services.WorkCollectionService
}

func NewWorkCollectionHandler(workCollectionService *services.WorkCollectionService) *WorkCollectionHandler {
	return &WorkCollectionHandler{workCollectionService: workCollectionService}
}

// GetWorkCollections 查询作品集列表。
// @Summary 查询作品集列表
// @Tags 后台作品集
// @Security BearerAuth
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键词"
// @Param published query bool false "发布状态"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections [get]
func (h *WorkCollectionHandler) GetWorkCollections(c *gin.Context) {
	var req dto.WorkCollectionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	collections, total, err := h.workCollectionService.GetCollections(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	utils.PageSuccess(c, collections, total, req.Page, req.PageSize)
}

// GetWorkCollectionByID 获取作品集详情。
// @Summary 获取作品集详情
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/admin/work-collections/{id} [get]
func (h *WorkCollectionHandler) GetWorkCollectionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid work collection id")
		return
	}

	collection, err := h.workCollectionService.GetCollectionByID(id)
	if err != nil {
		utils.NotFound(c, "")
		return
	}

	utils.Success(c, collection)
}

// CreateWorkCollection 创建作品集。
// @Summary 创建作品集
// @Tags 后台作品集
// @Security BearerAuth
// @Param request body dto.WorkCollectionCreateRequest true "作品集信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections [post]
func (h *WorkCollectionHandler) CreateWorkCollection(c *gin.Context) {
	var req dto.WorkCollectionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	collection, err := h.workCollectionService.CreateCollection(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, collection)
}

// UpdateWorkCollection 更新作品集。
// @Summary 更新作品集
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Param request body dto.WorkCollectionUpdateRequest true "作品集信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id} [put]
func (h *WorkCollectionHandler) UpdateWorkCollection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid work collection id")
		return
	}

	var req dto.WorkCollectionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	collection, err := h.workCollectionService.UpdateCollection(id, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, collection)
}

// DeleteWorkCollection 删除作品集。
// @Summary 删除作品集
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id} [delete]
func (h *WorkCollectionHandler) DeleteWorkCollection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid work collection id")
		return
	}

	if err := h.workCollectionService.DeleteCollection(id); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "work collection deleted successfully", nil)
}

// AddWorkCollectionItem 添加作品集条目。
// @Summary 添加作品集条目
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Param request body dto.WorkCollectionItemCreateRequest true "条目信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id}/items [post]
func (h *WorkCollectionHandler) AddWorkCollectionItem(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid work collection id")
		return
	}

	var req dto.WorkCollectionItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	item, err := h.workCollectionService.AddItem(collectionID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, item)
}

// UpdateWorkCollectionItem 更新作品集条目。
// @Summary 更新作品集条目
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Param itemId path string true "条目 ID"
// @Param request body dto.WorkCollectionItemUpdateRequest true "条目信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id}/items/{itemId} [put]
func (h *WorkCollectionHandler) UpdateWorkCollectionItem(c *gin.Context) {
	collectionID, itemID, ok := h.parseCollectionItemParams(c)
	if !ok {
		return
	}

	var req dto.WorkCollectionItemUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	item, err := h.workCollectionService.UpdateItem(collectionID, itemID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, item)
}

// DeleteWorkCollectionItem 删除作品集条目。
// @Summary 删除作品集条目
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Param itemId path string true "条目 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id}/items/{itemId} [delete]
func (h *WorkCollectionHandler) DeleteWorkCollectionItem(c *gin.Context) {
	collectionID, itemID, ok := h.parseCollectionItemParams(c)
	if !ok {
		return
	}

	if err := h.workCollectionService.DeleteItem(collectionID, itemID); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "work collection item deleted successfully", nil)
}

func (h *WorkCollectionHandler) parseCollectionItemParams(c *gin.Context) (uint64, uint64, bool) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid work collection id")
		return 0, 0, false
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid work collection item id")
		return 0, 0, false
	}

	return collectionID, itemID, true
}
