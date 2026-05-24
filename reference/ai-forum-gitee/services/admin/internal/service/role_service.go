package service

import (
	"context"
	"encoding/json"
	"fmt"

	"ai-forum/admin-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// RoleService handles role and permission assignment
type RoleService struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

// NewRoleService creates a new RoleService
func NewRoleService(db *pgxpool.Pool) *RoleService {
	return &RoleService{DB: db}
}

// SetRedis sets the Redis client (optional, called from main)
func (s *RoleService) SetRedis(rdb *redis.Client) {
	s.RDB = rdb
}

// ListRoles returns all available roles
func (s *RoleService) ListRoles(ctx context.Context) ([]*model.Role, error) {
	rows, err := s.DB.Query(ctx,
		"SELECT id, name, description, permissions, created_at FROM schema_admin.roles ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []*model.Role
	for rows.Next() {
		r := &model.Role{}
		var permJSON []byte
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &permJSON, &r.CreatedAt); err != nil {
			continue
		}
		if permJSON != nil {
			_ = json.Unmarshal(permJSON, &r.Permissions)
		}
		roles = append(roles, r)
	}
	return roles, nil
}

// AssignRole assigns a role to a user
func (s *RoleService) AssignRole(ctx context.Context, userID uint, roleID uint, operatorID uint) error {
	_, err := s.DB.Exec(ctx,
		"INSERT INTO schema_admin.user_roles (user_id, role_id, assigned_by) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		int64(userID), roleID, int64(operatorID),
	)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Get role name for Redis
	var roleName string
	err = s.DB.QueryRow(ctx,
		"SELECT name FROM schema_admin.roles WHERE id = $1",
		roleID,
	).Scan(&roleName)
	if err == nil && s.RDB != nil {
		// Store role in Redis for quick middleware lookup
		s.RDB.Set(ctx, fmt.Sprintf("role:%d", userID), roleName, 0)
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"user_id":     userID,
		"role_id":     roleID,
		"operator_id": operatorID,
	})
	_, err = s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'assign_role', 'role', $3, $4)`,
		int64(operatorID), "", int64(userID), string(detail),
	)
	return err
}

// RemoveRole removes a role from a user
func (s *RoleService) RemoveRole(ctx context.Context, userID uint, roleID uint, operatorID uint) error {
	_, err := s.DB.Exec(ctx,
		"DELETE FROM schema_admin.user_roles WHERE user_id = $1 AND role_id = $2",
		int64(userID), roleID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	if s.RDB != nil {
		s.RDB.Del(ctx, fmt.Sprintf("role:%d", userID))
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"user_id":     userID,
		"role_id":     roleID,
		"operator_id": operatorID,
	})
	_, err = s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'remove_role', 'role', $3, $4)`,
		int64(operatorID), "", int64(userID), string(detail),
	)
	return err
}

// GetUserRoles returns role names for a user
func (s *RoleService) GetUserRoles(ctx context.Context, userID uint) ([]string, error) {
	rows, err := s.DB.Query(ctx,
		`SELECT r.name FROM schema_admin.roles r
		 JOIN schema_admin.user_roles ur ON r.id = ur.role_id
		 WHERE ur.user_id = $1`,
		int64(userID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		roles = append(roles, name)
	}
	return roles, nil
}
