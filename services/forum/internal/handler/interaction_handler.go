package handler

import (
	"net/http"
	"strconv"

	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

// InteractionHandler handles likes and collections.
type InteractionHandler struct {
	Service *service.ForumService
}

// NewInteractionHandler creates a new InteractionHandler.
func NewInteractionHandler(svc *service.ForumService) *InteractionHandler {
	return &InteractionHandler{Service: svc}
}

// LikePost handles liking a post.
func (h *InteractionHandler) LikePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	resp, err := h.Service.LikePost(c.Request.Context(), userID, uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like post"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DislikePost handles disliking a post.
func (h *InteractionHandler) DislikePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	resp, err := h.Service.DislikePost(c.Request.Context(), userID, uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to dislike post"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CollectPost handles collecting/bookmarking a post.
func (h *InteractionHandler) CollectPost(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	resp, err := h.Service.CollectPost(c.Request.Context(), userID, uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect post"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// LikeComment handles liking a comment.
func (h *InteractionHandler) LikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	resp, err := h.Service.LikeComment(c.Request.Context(), userID, uint(commentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like comment"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DislikeComment handles disliking a comment.
func (h *InteractionHandler) DislikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	resp, err := h.Service.DislikeComment(c.Request.Context(), userID, uint(commentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to dislike comment"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ListMyCollections returns the current user's collected posts.
func (h *InteractionHandler) ListMyCollections(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, total, err := h.Service.ListUserCollections(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list collections"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
