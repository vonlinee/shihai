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
