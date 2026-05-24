package handler

import (
	"net/http"
	"strconv"

	"ai-forum/forum-service/internal/model"
	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

// CommentHandler handles comment operations
type CommentHandler struct {
	Service *service.ForumService
}

// NewCommentHandler creates a new CommentHandler
func NewCommentHandler(svc *service.ForumService) *CommentHandler {
	return &CommentHandler{Service: svc}
}

// ListComments returns comments for a post
func (h *CommentHandler) ListComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	comments, total, err := h.Service.ListComments(c.Request.Context(), uint(postID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// CreateComment adds a comment to a post
func (h *CommentHandler) CreateComment(c *gin.Context) {
	authorID := c.GetUint("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	comment, err := h.Service.CreateComment(c.Request.Context(), authorID, uint(postID), req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}
