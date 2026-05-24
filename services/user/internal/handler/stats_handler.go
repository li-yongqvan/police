package handler

import (
	"net/http"
	"strconv"

	"ai-forum/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

// UserStatsHandler handles user statistics endpoints
type UserStatsHandler struct {
	StatsService *service.UserStatsService
}

// NewUserStatsHandler creates a new UserStatsHandler
func NewUserStatsHandler(svc *service.UserStatsService) *UserStatsHandler {
	return &UserStatsHandler{StatsService: svc}
}

// GetStatsOverview returns user overview statistics
func (h *UserStatsHandler) GetStatsOverview(c *gin.Context) {
	overview, err := h.StatsService.GetOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": overview})
}

// GetDailyUsers returns daily user registration counts
func (h *UserStatsHandler) GetDailyUsers(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days > 30 {
		days = 30
	}
	users, err := h.StatsService.GetDailyUsers(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetLevelDistribution returns user distribution by level
func (h *UserStatsHandler) GetLevelDistribution(c *gin.Context) {
	dist, err := h.StatsService.GetLevelDistribution(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": dist})
}
