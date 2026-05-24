package handler

import (
	"net/http"
	"strconv"

	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// ModerationHandler handles sensitive word operations
type ModerationHandler struct {
	AdminService *service.AdminService
}

// NewModerationHandler creates a new ModerationHandler
func NewModerationHandler(svc *service.AdminService) *ModerationHandler {
	return &ModerationHandler{AdminService: svc}
}

// AddSensitiveWord adds a new sensitive word to the filter list
func (h *ModerationHandler) AddSensitiveWord(c *gin.Context) {
	var req struct {
		Word     string `json:"word" binding:"required"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	if req.Category == "" {
		req.Category = "general"
	}

	if err := h.AdminService.AddSensitiveWord(c.Request.Context(), req.Word, req.Category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "敏感词添加成功"})
}

// CheckSensitiveWords checks text for sensitive words (internal endpoint)
func (h *ModerationHandler) CheckSensitiveWords(c *gin.Context) {
	var req struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	clean, matched := h.AdminService.CheckSensitiveWords(req.Text)
	if !clean {
		c.JSON(http.StatusOK, gin.H{"clean": false, "matched_words": matched})
		return
	}
	c.JSON(http.StatusOK, gin.H{"clean": true, "matched_words": []string{}})
}

// ListSensitiveWords lists all sensitive words
func (h *ModerationHandler) ListSensitiveWords(c *gin.Context) {
	words, err := h.AdminService.ListSensitiveWords(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, words)
}

// DeleteSensitiveWord deletes a sensitive word by ID
func (h *ModerationHandler) DeleteSensitiveWord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.AdminService.DeleteSensitiveWord(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "敏感词已删除"})
}
