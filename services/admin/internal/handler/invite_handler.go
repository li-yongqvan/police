package handler

import (
	"net/http"

	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// InviteHandler handles invite code generation
type InviteHandler struct {
	AdminService *service.AdminService
	UserClient   *service.UserClient
}

// NewInviteHandler creates a new InviteHandler
func NewInviteHandler(adminSvc *service.AdminService, userClient *service.UserClient) *InviteHandler {
	return &InviteHandler{AdminService: adminSvc, UserClient: userClient}
}

// GenerateInviteCode generates a single invite code
func (h *InviteHandler) GenerateInviteCode(c *gin.Context) {
	code, err := h.AdminService.GenerateInviteCode(c.Request.Context(), h.UserClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": code})
}

// GenerateInviteCodesBatch generates multiple invite codes
func (h *InviteHandler) GenerateInviteCodesBatch(c *gin.Context) {
	var req struct {
		Count int `json:"count" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	if req.Count <= 0 || req.Count > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数量必须在1到1000之间"})
		return
	}

	codes, err := h.AdminService.GenerateInviteCodesBatch(c.Request.Context(), req.Count, h.UserClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"codes": codes, "count": len(codes)})
}
