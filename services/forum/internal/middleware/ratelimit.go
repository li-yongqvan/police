package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func allow(ctx context.Context, rdb *redis.Client, key string, max int, window time.Duration) bool {
	if rdb == nil {
		return true
	}
	n, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return true
	}
	if n == 1 {
		_ = rdb.Expire(ctx, key, window).Err()
	}
	return n <= int64(max)
}

// RateLimitByUser limits authenticated actions per user.
func RateLimitByUser(rdb *redis.Client, prefix string, max int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("user_id")
		if !ok {
			c.Next()
			return
		}
		key := fmt.Sprintf("rl:%s:user:%v", prefix, uid)
		if !allow(c.Request.Context(), rdb, key, max, window) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "发帖过于频繁，请稍后再试",
			})
			return
		}
		c.Next()
	}
}
