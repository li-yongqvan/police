package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"ai-forum/user-service/pkg/database"
	"ai-forum/user-service/pkg/jwt"

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

		// Set user_id in context
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}

		// Set username in context
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}

		// Set level in context (may need to refresh from DB for accuracy)
		if level, ok := claims["level"].(float64); ok {
			c.Set("level", int(level))
		} else if userID, ok := claims["user_id"].(float64); ok {
			// Try to fetch level from DB
			userIDUint := uint(userID)
			var userLevel int
			db, err := database.GetPool()
			if err == nil {
				row := db.QueryRow(context.Background(),
					"SELECT level FROM schema_auth.users WHERE id = $1", userIDUint)
				if err := row.Scan(&userLevel); err == nil {
					c.Set("level", userLevel)
				}
			}
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
