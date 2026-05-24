package handler

import (
	"net/http"
	"strconv"

	"ai-forum/admin-service/internal/client"
	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// BoardAdminHandler handles board management endpoints
type BoardAdminHandler struct {
	ForumClient *client.ForumClient
}

// NewBoardAdminHandler creates a new BoardAdminHandler
func NewBoardAdminHandler(forumClient *client.ForumClient) *BoardAdminHandler {
	return &BoardAdminHandler{ForumClient: forumClient}
}

// CreateBoard creates a new board
func (h *BoardAdminHandler) CreateBoard(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.ForumClient.CreateBoard(map[string]interface{}{
		"name":        req.Name,
		"slug":        req.Slug,
		"description": req.Description,
		"sort_order":  req.SortOrder,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "board created"})
}

// UpdateBoard updates a board
func (h *BoardAdminHandler) UpdateBoard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Enabled     *bool  `json:"enabled"`
		SortOrder   *int   `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	payload := map[string]interface{}{}
	if req.Name != "" {
		payload["name"] = req.Name
	}
	if req.Slug != "" {
		payload["slug"] = req.Slug
	}
	payload["description"] = req.Description
	if req.Enabled != nil {
		payload["enabled"] = *req.Enabled
	}
	if req.SortOrder != nil {
		payload["sort_order"] = *req.SortOrder
	}

	if err := h.ForumClient.UpdateBoard(uint(id), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "board updated"})
}

// DeleteBoard soft-deletes a board
func (h *BoardAdminHandler) DeleteBoard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	if err := h.ForumClient.DeleteBoard(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "board disabled"})
}

// ListBoards returns all boards including disabled
func (h *BoardAdminHandler) ListBoards(c *gin.Context) {
	boards, err := h.ForumClient.ListAllBoards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"boards": boards})
}

// PostAdminHandler handles admin post management
type PostAdminHandler struct {
	PostAdminService *service.PostAdminService
	ForumClient      *client.ForumClient
}

// NewPostAdminHandler creates a new PostAdminHandler
func NewPostAdminHandler(svc *service.PostAdminService, forumClient *client.ForumClient) *PostAdminHandler {
	return &PostAdminHandler{PostAdminService: svc, ForumClient: forumClient}
}

// ListPosts returns all posts for admin management
func (h *PostAdminHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	posts, total, err := h.ForumClient.ListAllPosts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// DeletePost deletes any post (admin action)
func (h *PostAdminHandler) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.PostAdminService.DeletePost(c.Request.Context(), uint(postID), operatorID, operatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

// SetPostFeatured sets/unsets featured flag
func (h *PostAdminHandler) SetPostFeatured(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req struct {
		Featured bool `json:"featured" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.PostAdminService.SetPostFeatured(c.Request.Context(), uint(postID), req.Featured, operatorID, operatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := "post set as featured"
	if !req.Featured {
		msg = "post unset as featured"
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// SetPostPinned sets/unsets pinned flag
func (h *PostAdminHandler) SetPostPinned(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req struct {
		Pinned bool `json:"pinned" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.PostAdminService.SetPostPinned(c.Request.Context(), uint(postID), req.Pinned, operatorID, operatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := "post set as pinned"
	if !req.Pinned {
		msg = "post unset as pinned"
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
