package handler

import (
	"net/http"
	"strconv"

	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// InviteAdminHandler handles invite code management endpoints
type InviteAdminHandler struct {
	UserClient *service.UserClient
}

// NewInviteAdminHandler creates a new InviteAdminHandler
func NewInviteAdminHandler(userClient *service.UserClient) *InviteAdminHandler {
	return &InviteAdminHandler{UserClient: userClient}
}

// ListInviteCodes returns paginated invite codes
func (h *InviteAdminHandler) ListInviteCodes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	codes, total, err := h.UserClient.ListInviteCodes(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"codes": codes,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetInviteCodeStatus returns details about an invite code
func (h *InviteAdminHandler) GetInviteCodeStatus(c *gin.Context) {
	code := c.Param("code")
	result, err := h.UserClient.GetInviteCodeStatus(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// VoidInviteCode marks an unused invite code as voided
func (h *InviteAdminHandler) VoidInviteCode(c *gin.Context) {
	code := c.Param("code")
	if err := h.UserClient.VoidInviteCode(code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "invite code voided"})
}
