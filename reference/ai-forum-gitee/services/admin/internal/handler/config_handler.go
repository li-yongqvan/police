package handler

import (
	"net/http"

	"ai-forum/admin-service/internal/service"

	"github.com/gin-gonic/gin"
)

// ConfigHandler handles system config endpoints
type ConfigHandler struct {
	ConfigService *service.ConfigService
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(svc *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{ConfigService: svc}
}

// GetConfig returns all system configuration
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	configs, err := h.ConfigService.GetAllConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make(map[string]string)
	for _, c := range configs {
		result[c.Key] = c.Value
	}
	c.JSON(http.StatusOK, gin.H{"configs": result, "list": configs})
}

// GetConfigByKey returns a single config entry
func (h *ConfigHandler) GetConfigByKey(c *gin.Context) {
	key := c.Param("key")
	config, err := h.ConfigService.GetConfig(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

// UpdateConfig updates a single config entry
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	key := c.Param("key")

	var req struct {
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	operatorID := int64(0)
	if oid, ok := c.Get("user_id"); ok {
		if o, ok := oid.(uint); ok {
			operatorID = int64(o)
		}
	}

	if err := h.ConfigService.UpdateConfig(c.Request.Context(), key, req.Value, operatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "config updated"})
}
