package service

import (
	"context"
	"fmt"
)

// Role names in the authorization domain. The Role Authority never emits a
// name outside this set; the three values are the locked Session Token
// Contract role domain (see CONTEXT.md, ADR-0003).
const (
	RoleStudent       = "student"
	RoleAdmin         = "admin"
	RolePlatformAdmin = "platform_admin"
)

// rolePriority is the Role Resolution priority order: the first name in
// this list that appears among a user's assignments wins.
var rolePriority = []string{RolePlatformAdmin, RoleAdmin}

// ResolveRoleName turns raw role assignments into the single Authoritative
// Role Name: priority plus fallback. Unknown names (including the legacy
// "user" role) collapse to student.
func ResolveRoleName(assignments []string) string {
	for _, p := range rolePriority {
		for _, a := range assignments {
			if a == p {
				return p
			}
		}
	}
	return RoleStudent
}

// ResolveAuthoritativeRole returns the authoritative role name for a user
// by resolving their role assignments from schema_admin. An empty
// assignment list is not an error: it resolves to student, so consumers
// never branch on "user not found".
func (s *RoleService) ResolveAuthoritativeRole(ctx context.Context, userID uint) (string, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT r.name FROM schema_admin.user_roles ur
		JOIN schema_admin.roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`, int64(userID))
	if err != nil {
		return "", fmt.Errorf("failed to query role assignments: %w", err)
	}
	defer rows.Close()

	var assignments []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return "", fmt.Errorf("failed to scan role assignment: %w", err)
		}
		assignments = append(assignments, name)
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("failed to iterate role assignments: %w", err)
	}
	return ResolveRoleName(assignments), nil
}