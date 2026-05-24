package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns service health status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "user-service",
	})
}
