package handler

import (
	"net/http"
	"strconv"

	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	Service *service.ExtrasService
}

func NewNotificationHandler(s *service.ExtrasService) *NotificationHandler {
	return &NotificationHandler{Service: s}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, total, err := h.Service.ListNotifications(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取消息失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notifications": items, "total": total, "page": page, "limit": limit})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.Service.MarkNotificationRead(c.Request.Context(), userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
