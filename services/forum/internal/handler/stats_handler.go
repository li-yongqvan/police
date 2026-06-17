package handler

import (
	"net/http"
	"strconv"

	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

// StatsHandler handles statistics endpoints
type StatsHandler struct {
	StatsService *service.StatsService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{StatsService: svc}
}

// GetStatsOverview returns forum overview statistics
func (h *StatsHandler) GetStatsOverview(c *gin.Context) {
	overview, err := h.StatsService.GetOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": overview})
}

// GetDailyPosts returns daily post counts
func (h *StatsHandler) GetDailyPosts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days > 30 {
		days = 30
	}
	posts, err := h.StatsService.GetDailyPosts(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// GetDailyComments returns daily comment counts
func (h *StatsHandler) GetDailyComments(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days > 30 {
		days = 30
	}
	comments, err := h.StatsService.GetDailyComments(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// GetBoardActivity returns board activity statistics
func (h *StatsHandler) GetBoardActivity(c *gin.Context) {
	activity, err := h.StatsService.GetBoardActivity(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": activity})
}
