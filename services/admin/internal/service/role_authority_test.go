package service

import "testing"

func TestResolveRoleName(t *testing.T) {
	tests := []struct {
		name        string
		assignments []string
		want        string
	}{
		{"no assignments fall back to student", nil, RoleStudent},
		{"empty assignments fall back to student", []string{}, RoleStudent},
		{"platform_admin outranks admin", []string{RoleAdmin, RolePlatformAdmin}, RolePlatformAdmin},
		{"platform_admin wins regardless of order", []string{RolePlatformAdmin, RoleAdmin}, RolePlatformAdmin},
		{"admin outranks legacy user", []string{"user", RoleAdmin}, RoleAdmin},
		{"legacy user clamps to student", []string{"user"}, RoleStudent},
		{"unknown name clamps to student", []string{"moderator"}, RoleStudent},
		{"unknown names do not block known ones", []string{"moderator", RoleAdmin}, RoleAdmin},
		{"duplicates keep the priority name", []string{RoleAdmin, RoleAdmin}, RoleAdmin},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveRoleName(tt.assignments); got != tt.want {
				t.Fatalf("ResolveRoleName(%v) = %q, want %q", tt.assignments, got, tt.want)
			}
		})
	}
}