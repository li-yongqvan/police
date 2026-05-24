package handler

import (
	"net/http"
	"strconv"

	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

// InteractionHandler handles likes and collections
type InteractionHandler struct {
	Service *service.ForumService
}

// NewInteractionHandler creates a new InteractionHandler
func NewInteractionHandler(svc *service.ForumService) *InteractionHandler {
	return &InteractionHandler{Service: svc}
}

// LikePost handles liking a post
func (h *InteractionHandler) LikePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}

	resp, err := h.Service.LikePost(c.Request.Context(), userID, uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞操作失败"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CollectPost handles collecting/bookmarking a post
func (h *InteractionHandler) CollectPost(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}

	resp, err := h.Service.CollectPost(c.Request.Context(), userID, uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "收藏操作失败"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
