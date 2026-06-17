package handler

import (
	"net/http"
	"strconv"

	"ai-forum/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

// UserAdminHandler handles admin operations on user endpoints
type UserAdminHandler struct {
	AdminService *service.UserAdminService
}

// NewUserAdminHandler creates a new UserAdminHandler
func NewUserAdminHandler(svc *service.UserAdminService) *UserAdminHandler {
	return &UserAdminHandler{AdminService: svc}
}

// BanUser bans a user account
func (h *UserAdminHandler) BanUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

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

	if err := h.AdminService.BanUser(c.Request.Context(), uint(userID), req.Reason, operatorID, operatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user banned"})
}

// UnbanUser unbans a user account
func (h *UserAdminHandler) UnbanUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

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

	if err := h.AdminService.UnbanUser(c.Request.Context(), uint(userID), operatorID, operatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unbanned"})
}

// ListUsers returns paginated user list with optional status filter
func (h *UserAdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	statusFilter := c.DefaultQuery("status", "all")

	users, total, err := h.AdminService.ListUsers(c.Request.Context(), page, limit, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateUserLevel updates a user's level
func (h *UserAdminHandler) UpdateUserLevel(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Level int `json:"level" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

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

	if err := h.AdminService.UpdateUserLevel(c.Request.Context(), uint(userID), req.Level, operatorID, operatorName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user level updated"})
}

// GetUserLogs returns operation logs for a user
func (h *UserAdminHandler) GetUserLogs(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	logs, total, err := h.AdminService.GetUserLogs(c.Request.Context(), uint(userID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetInviteCodeStatus returns details about an invite code
func (h *UserAdminHandler) GetInviteCodeStatus(c *gin.Context) {
	code := c.Param("code")

	detail, err := h.AdminService.GetInviteCodeStatus(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}

// ListInviteCodes returns paginated invite codes
func (h *UserAdminHandler) ListInviteCodes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	codes, total, err := h.AdminService.ListInviteCodes(c.Request.Context(), page, limit)
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

// VoidInviteCode marks an unused invite code as voided
func (h *UserAdminHandler) VoidInviteCode(c *gin.Context) {
	code := c.Param("code")

	if err := h.AdminService.VoidInviteCode(c.Request.Context(), code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "invite code voided"})
}
