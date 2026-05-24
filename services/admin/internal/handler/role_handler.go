package handler

import (
	"net/http"
	"strconv"

	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// RoleHandler handles role management endpoints
type RoleHandler struct {
	RoleService *service.RoleService
}

// NewRoleHandler creates a new RoleHandler
func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{RoleService: svc}
}

// ListRoles returns all available roles with permissions
func (h *RoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.RoleService.ListRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

// AssignRole assigns a role to a user
func (h *RoleHandler) AssignRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	operatorID := uint(0)
	if oid, ok := c.Get("user_id"); ok {
		if o, ok := oid.(uint); ok {
			operatorID = o
		}
	}

	if err := h.RoleService.AssignRole(c.Request.Context(), uint(userID), req.RoleID, operatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role assigned"})
}

// RemoveRole removes a role from a user
func (h *RoleHandler) RemoveRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	operatorID := uint(0)
	if oid, ok := c.Get("user_id"); ok {
		if o, ok := oid.(uint); ok {
			operatorID = o
		}
	}

	if err := h.RoleService.RemoveRole(c.Request.Context(), uint(userID), uint(roleID), operatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role removed"})
}

// GetUserRoles returns the roles assigned to a user
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	roles, err := h.RoleService.GetUserRoles(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"roles": roles})
}
