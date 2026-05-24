package handler

import (
	"net/http"
	"strconv"

	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// AuditHandler handles audit workflow endpoints
type AuditHandler struct {
	AuditService *service.AuditService
}

// NewAuditHandler creates a new AuditHandler
func NewAuditHandler(svc *service.AuditService) *AuditHandler {
	return &AuditHandler{AuditService: svc}
}

// ListPendingAudit returns posts pending review
func (h *AuditHandler) ListPendingAudit(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	posts, total, err := h.AuditService.ListPendingAudit(c.Request.Context(), page, limit)
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

// ApprovePost approves a pending post
func (h *AuditHandler) ApprovePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.AuditService.ApprovePost(c.Request.Context(), uint(postID), operatorID, operatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post approved"})
}

// RejectPost rejects a pending post
func (h *AuditHandler) RejectPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.AuditService.RejectPost(c.Request.Context(), uint(postID), operatorID, operatorName, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post rejected"})
}

// BatchDeletePosts deletes multiple posts at once
func (h *AuditHandler) BatchDeletePosts(c *gin.Context) {
	var req struct {
		PostIDs []uint `json:"post_ids" binding:"required"`
		Reason  string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if len(req.PostIDs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch size limited to 100"})
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.AuditService.BatchDeletePosts(c.Request.Context(), req.PostIDs, operatorID, operatorName, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posts deleted", "count": len(req.PostIDs)})
}

func getOperatorInfo(c *gin.Context) (uint, string) {
	operatorID := uint(0)
	operatorName := "admin_system"
	if oid, ok := c.Get("user_id"); ok {
		if o, ok := oid.(uint); ok {
			operatorID = o
		}
	}
	if on, ok := c.Get("username"); ok {
		if n, ok := on.(string); ok {
			operatorName = n
		}
	}
	return operatorID, operatorName
}
