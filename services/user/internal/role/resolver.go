package role

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Resolver is the port BC-Identity depends on to obtain the Authoritative
// Role Name. It is total: caching and the degradation ladder are resolved
// behind the seam, so callers never handle errors.
type Resolver interface {
	Resolve(ctx context.Context, userID uint) string
}

// Fetcher is the transport seam to the Role Authority. The HTTP adapter
// and the in-memory test adapter both implement it.
type Fetcher interface {
	Fetch(ctx context.Context, userID uint) (Name, error)
}

// Cache is the consumer-side role cache seam over Redis.
type Cache interface {
	Get(ctx context.Context, userID uint) (Name, bool)
	Set(ctx context.Context, userID uint, name Name, ttl time.Duration) error
}

// CachedResolver implements Resolver with the read-through degradation
// ladder: cache hit -> return; miss -> fetch -> validate -> cache -> return;
// fetch failure -> student + log. Degraded values are never written to the
// cache, so an outage cannot poison it.
type CachedResolver struct {
	upstream Fetcher
	cache    Cache
	ttl      time.Duration
	log      *slog.Logger
}

func NewCachedResolver(upstream Fetcher, cache Cache, ttl time.Duration, log *slog.Logger) *CachedResolver {
	if log == nil {
		log = slog.Default()
	}
	return &CachedResolver{upstream: upstream, cache: cache, ttl: ttl, log: log}
}

func (r *CachedResolver) Resolve(ctx context.Context, userID uint) string {
	if cached, ok := r.cache.Get(ctx, userID); ok {
		if _, valid := ValidName(string(cached)); valid {
			return string(cached)
		}
	}

	name, err := r.upstream.Fetch(ctx, userID)
	if err == nil {
		canonical, valid := ValidName(string(name))
		if valid {
			if serr := r.cache.Set(ctx, userID, canonical, r.ttl); serr != nil {
				r.log.WarnContext(ctx, "role cache write failed", "user_id", userID, "err", serr)
			}
			return string(canonical)
		}
		err = fmt.Errorf("role authority returned out-of-domain role %q", string(name))
	}

	r.log.WarnContext(ctx, "role authority unavailable, degrading to student", "user_id", userID, "err", err)
	return string(Student)
}
