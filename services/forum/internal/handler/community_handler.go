package handler

import (
	"net/http"

	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CommunityHandler struct {
	Service *service.ExtrasService
}

func NewCommunityHandler(s *service.ExtrasService) *CommunityHandler {
	return &CommunityHandler{Service: s}
}

func (h *CommunityHandler) Stats(c *gin.Context) {
	stats, err := h.Service.GetCommunityStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
