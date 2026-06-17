package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// UserListItem represents a user in admin listing
type UserListItem struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Level     int       `json:"level"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Roles     []string  `json:"roles"`
}

// InviteCodeDetail represents invite code details
type InviteCodeDetail struct {
	Code            string     `json:"code"`
	Status          string     `json:"status"`
	CreatedBy       int64      `json:"created_by"`
	CreatedByName   string     `json:"created_by_name"`
	UsedBy          *int64     `json:"used_by"`
	UsedByName      *string    `json:"used_by_name"`
	UsedAt          *time.Time `json:"used_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// InviteCodeItem represents an invite code in admin listing
type InviteCodeItem struct {
	Code          string     `json:"code"`
	CreatedByName string     `json:"created_by_name"`
	Status        string     `json:"status"`
	UsedByName    *string    `json:"used_by_name"`
	UsedAt        *time.Time `json:"used_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

// LogEntry represents an operation log entry
type LogEntry struct {
	ID               uint                 `json:"id"`
	OperatorID       int64                `json:"operator_id"`
	OperatorUsername string               `json:"operator_username"`
	Action           string               `json:"action"`
	TargetType       string               `json:"target_type"`
	TargetID         int64                `json:"target_id"`
	Detail           map[string]interface{} `json:"detail"`
	CreatedAt        time.Time            `json:"created_at"`
}

// UserAdminService handles admin operations on user-service
type UserAdminService struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

// NewUserAdminService creates a new UserAdminService
func NewUserAdminService(db *pgxpool.Pool, rdb *redis.Client) *UserAdminService {
	return &UserAdminService{DB: db, RDB: rdb}
}

// BanUser bans a user: updates status, deletes refresh token, logs operation
func (s *UserAdminService) BanUser(ctx context.Context, userID uint, reason string, operatorID uint, operatorName string) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update user status to banned
	_, err = tx.Exec(ctx,
		"UPDATE schema_auth.users SET status = 'banned', updated_at = NOW() WHERE id = $1",
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to ban user: %w", err)
	}

	// Log operation
	_, err = tx.Exec(ctx,
		`INSERT INTO schema_auth.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'ban_user', 'user', $3, $4)`,
		operatorID, operatorName, int64(userID), fmt.Sprintf(`{"reason": "%s"}`, reason),
	)
	if err != nil {
		return fmt.Errorf("failed to log operation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Delete refresh token for instant ban effect (per D-04, D-05)
	s.RDB.Del(ctx, fmt.Sprintf("refresh:%d", userID))

	return nil
}

// UnbanUser sets user status back to active
func (s *UserAdminService) UnbanUser(ctx context.Context, userID uint, operatorID uint, operatorName string) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_auth.users SET status = 'active', updated_at = NOW() WHERE id = $1",
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to unban user: %w", err)
	}

	_, err = s.DB.Exec(ctx,
		`INSERT INTO schema_auth.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'unban_user', 'user', $3, '{}')`,
		operatorID, operatorName, int64(userID),
	)
	return err
}

// ListUsers returns paginated users with optional status filter
func (s *UserAdminService) ListUsers(ctx context.Context, page, limit int, statusFilter string) ([]UserListItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	countQuery := "SELECT COUNT(*) FROM schema_auth.users"
	var total int
	if statusFilter != "" && statusFilter != "all" {
		countQuery += " WHERE status = $1"
		err := s.DB.QueryRow(ctx, countQuery, statusFilter).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count users: %w", err)
		}
	} else {
		err := s.DB.QueryRow(ctx, countQuery).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count users: %w", err)
		}
	}

	query := `
		SELECT u.id, u.username, u.nickname, u.level, u.status, u.created_at
		FROM schema_auth.users u`
	args := []interface{}{}
	if statusFilter != "" && statusFilter != "all" {
		query += " WHERE u.status = $1"
		args = append(args, statusFilter)
	}
	query += " ORDER BY u.id DESC"
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []UserListItem
	for rows.Next() {
		var u UserListItem
		if err := rows.Scan(&u.ID, &u.Username, &u.Nickname, &u.Level, &u.Status, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, total, nil
}

// UpdateUserLevel updates a user's level and logs the operation
func (s *UserAdminService) UpdateUserLevel(ctx context.Context, userID uint, newLevel int, operatorID uint, operatorName string) error {
	if newLevel < 0 || newLevel > 5 {
		return fmt.Errorf("level must be between 0 and 5")
	}

	_, err := s.DB.Exec(ctx,
		"UPDATE schema_auth.users SET level = $1, updated_at = NOW() WHERE id = $2",
		newLevel, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user level: %w", err)
	}

	_, err = s.DB.Exec(ctx,
		`INSERT INTO schema_auth.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'update_user_level', 'user', $3, $4)`,
		operatorID, operatorName, int64(userID), fmt.Sprintf(`{"new_level": %d}`, newLevel),
	)
	return err
}

// GetUserLogs returns operation logs for a specific user (target or operator)
func (s *UserAdminService) GetUserLogs(ctx context.Context, userID uint, page, limit int) ([]LogEntry, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var total int
	err := s.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM schema_auth.operation_logs
		 WHERE target_id = $1 OR operator_id = $1`,
		int64(userID),
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count logs: %w", err)
	}

	rows, err := s.DB.Query(ctx,
		`SELECT id, operator_id, operator_username, action, target_type, target_id, detail, created_at
		 FROM schema_auth.operation_logs
		 WHERE target_id = $1 OR operator_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		int64(userID), limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.ID, &l.OperatorID, &l.OperatorUsername, &l.Action, &l.TargetType, &l.TargetID, &l.Detail, &l.CreatedAt); err != nil {
			continue
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}

// GetInviteCodeStatus returns details about an invite code
func (s *UserAdminService) GetInviteCodeStatus(ctx context.Context, code string) (*InviteCodeDetail, error) {
	detail := &InviteCodeDetail{}
	err := s.DB.QueryRow(ctx,
		`SELECT ic.code, ic.status, ic.created_by, COALESCE(cb.username, ''), ic.used_by, ub.username, ic.used_at, ic.created_at
		 FROM schema_auth.invite_codes ic
		 LEFT JOIN schema_auth.users cb ON ic.created_by = cb.id
		 LEFT JOIN schema_auth.users ub ON ic.used_by = ub.id
		 WHERE ic.code = $1`,
		code,
	).Scan(&detail.Code, &detail.Status, &detail.CreatedBy, &detail.CreatedByName, &detail.UsedBy, &detail.UsedByName, &detail.UsedAt, &detail.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invite code not found")
	}
	return detail, nil
}

// ListInviteCodes returns paginated invite codes
func (s *UserAdminService) ListInviteCodes(ctx context.Context, page, limit int) ([]InviteCodeItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var total int
	err := s.DB.QueryRow(ctx, "SELECT COUNT(*) FROM schema_auth.invite_codes").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count invite codes: %w", err)
	}

	rows, err := s.DB.Query(ctx,
		`SELECT ic.code, COALESCE(cb.username, ''), ic.status, ub.username, ic.used_at, ic.created_at
		 FROM schema_auth.invite_codes ic
		 LEFT JOIN schema_auth.users cb ON ic.created_by = cb.id
		 LEFT JOIN schema_auth.users ub ON ic.used_by = ub.id
		 ORDER BY ic.created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query invite codes: %w", err)
	}
	defer rows.Close()

	var items []InviteCodeItem
	for rows.Next() {
		var item InviteCodeItem
		if err := rows.Scan(&item.Code, &item.CreatedByName, &item.Status, &item.UsedByName, &item.UsedAt, &item.CreatedAt); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, total, nil
}

// VoidInviteCode marks an unused invite code as voided
func (s *UserAdminService) VoidInviteCode(ctx context.Context, code string) error {
	result, err := s.DB.Exec(ctx,
		"UPDATE schema_auth.invite_codes SET status = 'voided' WHERE code = $1 AND status = 'unused'",
		code,
	)
	if err != nil {
		return fmt.Errorf("failed to void invite code: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("invite code not found or already used/voided")
	}
	return nil
}
