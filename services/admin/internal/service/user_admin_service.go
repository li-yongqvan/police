package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserAdminService handles admin user management (delegates to user-service)
type UserAdminService struct {
	DB         *pgxpool.Pool
	UserClient *UserClient
}

// NewUserAdminService creates a new UserAdminService
func NewUserAdminService(db *pgxpool.Pool, userClient *UserClient) *UserAdminService {
	return &UserAdminService{DB: db, UserClient: userClient}
}

// BanUser delegates to user-service and logs operation in admin
func (s *UserAdminService) BanUser(ctx context.Context, userID uint, reason string, operatorID uint, operatorName string) error {
	if err := s.UserClient.BanUser(userID, reason); err != nil {
		return fmt.Errorf("failed to ban user: %w", err)
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"operator_id":   operatorID,
		"operator_name": operatorName,
		"reason":        reason,
	})
	_, err := s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'ban_user', 'user', $3, $4)`,
		int64(operatorID), operatorName, int64(userID), string(detail),
	)
	return err
}

// UnbanUser delegates to user-service and logs operation in admin
func (s *UserAdminService) UnbanUser(ctx context.Context, userID uint, operatorID uint, operatorName string) error {
	if err := s.UserClient.UnbanUser(userID); err != nil {
		return fmt.Errorf("failed to unban user: %w", err)
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"operator_id":   operatorID,
		"operator_name": operatorName,
	})
	_, err := s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'unban_user', 'user', $3, $4)`,
		int64(operatorID), operatorName, int64(userID), string(detail),
	)
	return err
}

// ListUsers delegates to user-service
func (s *UserAdminService) ListUsers(ctx context.Context, page, limit int, status string) ([]map[string]interface{}, int, error) {
	return s.UserClient.ListUsers(page, limit, status)
}

// UpdateUserLevel delegates to user-service and logs operation
func (s *UserAdminService) UpdateUserLevel(ctx context.Context, userID uint, newLevel int, operatorID uint, operatorName string) error {
	if newLevel < 0 || newLevel > 5 {
		return fmt.Errorf("level must be between 0 and 5")
	}
	if err := s.UserClient.UpdateUserLevel(userID, newLevel); err != nil {
		return fmt.Errorf("failed to update user level: %w", err)
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"operator_id":   operatorID,
		"operator_name": operatorName,
		"new_level":     newLevel,
	})
	_, err := s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'update_user_level', 'user', $3, $4)`,
		int64(operatorID), operatorName, int64(userID), string(detail),
	)
	return err
}

// GetUserLogs delegates to user-service
func (s *UserAdminService) GetUserLogs(ctx context.Context, userID uint, page, limit int) ([]map[string]interface{}, int, error) {
	return s.UserClient.GetUserLogs(userID, page, limit)
}
