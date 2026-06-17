package handler

import (
	"net/http"
	"strconv"

	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// StatsHandler handles admin stats endpoints
type StatsHandler struct {
	StatsService *service.StatsService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{StatsService: svc}
}

// GetOverview returns admin dashboard overview stats
func (h *StatsHandler) GetOverview(c *gin.Context) {
	overview, err := h.StatsService.GetOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": overview})
}

// GetDailyStats returns daily statistics
func (h *StatsHandler) GetDailyStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days > 30 {
		days = 30
	}
	stats, err := h.StatsService.GetDailyStats(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
