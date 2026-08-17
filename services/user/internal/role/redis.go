package role

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache stores resolved roles under role:authority:<id>. The prefix is
// deliberately distinct from the retired admin-service role:<id> keys.
type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) key(userID uint) string {
	return fmt.Sprintf("role:authority:%d", userID)
}

// Get returns the cached role, or a miss on any Redis error: cache failures
// must never fail a login, only route through the authority.
func (c *RedisCache) Get(ctx context.Context, userID uint) (Name, bool) {
	val, err := c.client.Get(ctx, c.key(userID)).Result()
	if err != nil {
		return "", false
	}
	return Name(val), true
}

func (c *RedisCache) Set(ctx context.Context, userID uint, name Name, ttl time.Duration) error {
	return c.client.Set(ctx, c.key(userID), string(name), ttl).Err()
}
