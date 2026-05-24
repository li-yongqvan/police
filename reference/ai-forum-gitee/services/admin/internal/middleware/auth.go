package middleware

import (
	"net/http"
	"os"
	"strings"

	"ai-forum/admin-service/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and checks admin role
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
		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Check admin role
		if role, ok := claims["role"].(string); !ok || (role != "admin" && role != "platform_admin") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}

		c.Next()
	}
}
