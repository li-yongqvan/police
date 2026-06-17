package handler

import (
	"net/http"
	"strconv"

	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// UserAdminHandler handles admin user management endpoints (proxy to user-service)
type UserAdminHandler struct {
	UserAdminService *service.UserAdminService
}

// NewUserAdminHandler creates a new UserAdminHandler
func NewUserAdminHandler(svc *service.UserAdminService) *UserAdminHandler {
	return &UserAdminHandler{UserAdminService: svc}
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

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.UserAdminService.BanUser(c.Request.Context(), uint(userID), req.Reason, operatorID, operatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user banned"})
}

// UnbanUser restores a banned user account
func (h *UserAdminHandler) UnbanUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.UserAdminService.UnbanUser(c.Request.Context(), uint(userID), operatorID, operatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unbanned"})
}

// ListUsers returns paginated user list
func (h *UserAdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	statusFilter := c.DefaultQuery("status", "all")

	users, total, err := h.UserAdminService.ListUsers(c.Request.Context(), page, limit, statusFilter)
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

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.UserAdminService.UpdateUserLevel(c.Request.Context(), uint(userID), req.Level, operatorID, operatorName); err != nil {
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

	logs, total, err := h.UserAdminService.GetUserLogs(c.Request.Context(), uint(userID), page, limit)
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
