package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"ai-forum/forum-service/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens on protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "default-secret-change-in-production"
		}
		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		if level, ok := claims["level"].(float64); ok {
			c.Set("level", int(level))
		}

		c.Next()
	}
}

// OptionalAuthMiddleware sets user context when a valid Bearer token is present.
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "default-secret-change-in-production"
		}
		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.Next()
			return
		}
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		if level, ok := claims["level"].(float64); ok {
			c.Set("level", int(level))
		}
		c.Next()
	}
}

// RequireLevel checks if the user has sufficient level
func RequireLevel(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		level, exists := c.Get("level")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}
		if levelInt, ok := level.(int); ok && levelInt < minLevel {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足，需要等级" + strconv.Itoa(minLevel)})
			return
		}
		c.Next()
	}
}
