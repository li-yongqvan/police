package role

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

type fakeFetcher struct {
	name  Name
	err   error
	calls int
}

func (f *fakeFetcher) Fetch(context.Context, uint) (Name, error) {
	f.calls++
	return f.name, f.err
}

type fakeCache map[uint]Name

func (c fakeCache) Get(_ context.Context, userID uint) (Name, bool) {
	n, ok := c[userID]
	return n, ok
}

func (c fakeCache) Set(_ context.Context, userID uint, name Name, _ time.Duration) error {
	c[userID] = name
	return nil
}

func newTestResolver(f *fakeFetcher, c fakeCache) *CachedResolver {
	return NewCachedResolver(f, c, time.Minute, slog.Default())
}

func TestResolveServesCacheWithoutCallingAuthority(t *testing.T) {
	cache := fakeCache{7: Admin}
	f := &fakeFetcher{name: Admin}
	r := newTestResolver(f, cache)

	if got := r.Resolve(context.Background(), 7); got != string(Admin) {
		t.Fatalf("got %q, want admin", got)
	}
	if f.calls != 0 {
		t.Fatalf("authority called %d times, want 0 on cache hit", f.calls)
	}
}

func TestResolveFetchesAndCachesOnMiss(t *testing.T) {
	cache := fakeCache{}
	f := &fakeFetcher{name: PlatformAdmin}
	r := newTestResolver(f, cache)

	if got := r.Resolve(context.Background(), 9); got != string(PlatformAdmin) {
		t.Fatalf("got %q, want platform_admin", got)
	}
	if cached, ok := cache[9]; !ok || cached != PlatformAdmin {
		t.Fatalf("cache = %q present=%v, want platform_admin cached", cached, ok)
	}
}

func TestResolveDegradesToStudentWhenAuthorityFails(t *testing.T) {
	cache := fakeCache{}
	f := &fakeFetcher{err: errors.New("dial timeout")}
	r := newTestResolver(f, cache)

	if got := r.Resolve(context.Background(), 9); got != string(Student) {
		t.Fatalf("got %q, want student", got)
	}
	if _, ok := cache[9]; ok {
		t.Fatalf("degraded value must not be cached")
	}
}

func TestResolveTreatsOutOfDomainRoleAsFailure(t *testing.T) {
	cache := fakeCache{}
	f := &fakeFetcher{name: "moderator"}
	r := newTestResolver(f, cache)

	if got := r.Resolve(context.Background(), 9); got != string(Student) {
		t.Fatalf("got %q, want student", got)
	}
	if _, ok := cache[9]; ok {
		t.Fatalf("out-of-domain value must not be cached")
	}
}

func TestResolveTreatsMalformedCacheAsMiss(t *testing.T) {
	cache := fakeCache{9: "moderator"}
	f := &fakeFetcher{name: Admin}
	r := newTestResolver(f, cache)

	if got := r.Resolve(context.Background(), 9); got != string(Admin) {
		t.Fatalf("got %q, want admin", got)
	}
}
