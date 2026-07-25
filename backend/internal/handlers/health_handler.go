package handlers

import "github.com/gin-gonic/gin"

// HealthHandler adapts service health APIs to HTTP.
type HealthHandler struct{}

// NewHealthHandler creates a HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check reports whether the service is running.
// @Summary 服务健康检查
// @Tags 系统
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
