package handlers

import (
	"log"

	"shihai/internal/dto"
	"shihai/pkg/utils"

	"github.com/gin-gonic/gin"
)

type dashboardService interface {
	GetDashboard() (*dto.DashboardResponse, error)
}

// DashboardHandler adapts admin dashboard APIs to HTTP.
type DashboardHandler struct {
	dashboardService dashboardService
}

// NewDashboardHandler creates a DashboardHandler.
func NewDashboardHandler(dashboardService dashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// GetDashboard handles GET /api/admin/dashboard.
// @Summary 获取后台仪表盘数据
// @Description 获取后台仪表盘统计卡片和最近活动
// @Tags 管理后台
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{data=dto.DashboardResponse}
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/dashboard [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	dashboard, err := h.dashboardService.GetDashboard()
	if err != nil {
		log.Printf("get admin dashboard failed: %v", err)
		utils.InternalServerError(c, "获取仪表盘数据失败")
		return
	}
	utils.Success(c, dashboard)
}
