package middleware

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// DemoLoginGuard restricts demo-login to development or allowlisted client IPs.
func DemoLoginGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("APP_ENV") == "development" {
			c.Next()
			return
		}

		clientIP := net.ParseIP(c.ClientIP())
		if clientIP == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "demo login not allowed"})
			return
		}

		if clientIP.IsLoopback() || clientIP.IsPrivate() {
			c.Next()
			return
		}

		for _, raw := range strings.Split(os.Getenv("DEMO_LOGIN_ALLOWLIST"), ",") {
			entry := strings.TrimSpace(raw)
			if entry == "" {
				continue
			}
			if strings.Contains(entry, "/") {
				_, network, err := net.ParseCIDR(entry)
				if err == nil && network.Contains(clientIP) {
					c.Next()
					return
				}
				continue
			}
			if ip := net.ParseIP(entry); ip != nil && ip.Equal(clientIP) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "demo login not allowed from this network"})
	}
}
