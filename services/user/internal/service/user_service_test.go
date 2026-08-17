package service

import (
	"context"
	"testing"
)

type staticResolver struct {
	role string
}

func (s *staticResolver) Resolve(context.Context, uint) string {
	return s.role
}

func TestResolveAppRoleDelegatesToRoleAuthorityPort(t *testing.T) {
	svc := &UserService{roles: &staticResolver{role: "platform_admin"}}
	if got := svc.ResolveAppRole(context.Background(), 42); got != "platform_admin" {
		t.Fatalf("ResolveAppRole = %q, want platform_admin", got)
	}
}
