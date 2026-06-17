package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

// ForumAdminHandler handles admin operations on forum endpoints
type ForumAdminHandler struct {
	AdminService *service.ForumAdminService
}

// NewForumAdminHandler creates a new ForumAdminHandler
func NewForumAdminHandler(svc *service.ForumAdminService) *ForumAdminHandler {
	return &ForumAdminHandler{AdminService: svc}
}

// ListPendingPosts returns posts pending review (paginated)
func (h *ForumAdminHandler) ListPendingPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	posts, total, err := h.AdminService.ListPendingPosts(c.Request.Context(), page, limit)
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

// ListAllPosts returns all posts for admin management (paginated)
func (h *ForumAdminHandler) ListAllPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	posts, total, err := h.AdminService.ListAllPosts(c.Request.Context(), page, limit)
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

// ApprovePost approves a pending post (sets status to published)
func (h *ForumAdminHandler) ApprovePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.AdminService.ApprovePost(c.Request.Context(), uint(postID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post approved"})
}

// RejectPost rejects a pending post (sets status to rejected)
func (h *ForumAdminHandler) RejectPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.AdminService.RejectPost(c.Request.Context(), uint(postID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post rejected"})
}

// AdminDeletePost deletes any post (admin action, no author check)
func (h *ForumAdminHandler) AdminDeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.AdminService.AdminDeletePost(c.Request.Context(), uint(postID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

// SetPostFeatured sets/unsets the featured flag on a post
func (h *ForumAdminHandler) SetPostFeatured(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req struct {
		Featured bool `json:"featured"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.AdminService.SetPostFeatured(c.Request.Context(), uint(postID), req.Featured); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := "post set as featured"
	if !req.Featured {
		msg = "post unset as featured"
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// SetPostPinned sets/unsets the pinned flag on a post
func (h *ForumAdminHandler) SetPostPinned(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req struct {
		Pinned bool `json:"pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.AdminService.SetPostPinned(c.Request.Context(), uint(postID), req.Pinned); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := "post set as pinned"
	if !req.Pinned {
		msg = "post unset as pinned"
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// ChangePostStatus changes post status (approved/rejected/published)
func (h *ForumAdminHandler) ChangePostStatus(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.AdminService.ChangePostStatus(c.Request.Context(), uint(postID), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post status changed to " + req.Status})
}

// BatchDeletePosts deletes multiple posts at once
func (h *ForumAdminHandler) BatchDeletePosts(c *gin.Context) {
	var req struct {
		PostIDs []uint `json:"post_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if len(req.PostIDs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch size limited to 100"})
		return
	}

	if err := h.AdminService.BatchDeletePosts(c.Request.Context(), req.PostIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posts deleted", "count": len(req.PostIDs)})
}

// CreateBoard creates a new board (admin only)
func (h *ForumAdminHandler) CreateBoard(c *gin.Context) {
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

	if err := h.AdminService.CreateBoard(c.Request.Context(), req.Name, req.Slug, req.Description, req.SortOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "board created"})
}

// UpdateBoard updates a board's attributes
func (h *ForumAdminHandler) UpdateBoard(c *gin.Context) {
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

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	if err := h.AdminService.UpdateBoard(c.Request.Context(), uint(id), req.Name, req.Slug, req.Description, enabled, sortOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "board updated"})
}

// DeleteBoard soft-deletes a board
func (h *ForumAdminHandler) DeleteBoard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	if err := h.AdminService.DeleteBoard(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "board disabled"})
}

// ListAllBoardsInternal returns all boards including disabled (internal endpoint)
func (h *ForumAdminHandler) ListAllBoardsInternal(c *gin.Context) {
	boards, err := h.AdminService.ListAllBoards(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"boards": boards})
}

// ListBoardsInternal returns all boards including disabled (legacy format for ForumClient compatibility)
func (h *ForumAdminHandler) ListBoardsInternal(c *gin.Context) {
	boards, err := h.AdminService.ListAllBoards(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"boards": boards})
}

// AdminStatus handler for changing post status via dedicated endpoint
func (h *ForumAdminHandler) AdminStatus(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	status := req.Status
	if status == "approved" {
		status = "published"
	}
	if !isValidStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	if err := h.AdminService.ChangePostStatus(c.Request.Context(), uint(postID), status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post status changed to " + status})
}

func isValidStatus(status string) bool {
	valid := []string{"published", "pending_review", "rejected", "deleted"}
	for _, v := range valid {
		if strings.EqualFold(v, status) {
			return true
		}
	}
	return false
}
