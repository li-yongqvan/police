package service

import (
	"context"
	"encoding/json"
	"fmt"

	"ai-forum/admin-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConfigService handles system configuration operations
type ConfigService struct {
	DB *pgxpool.Pool
}

// NewConfigService creates a new ConfigService
func NewConfigService(db *pgxpool.Pool) *ConfigService {
	return &ConfigService{DB: db}
}

// GetAllConfig returns all system configuration entries
func (s *ConfigService) GetAllConfig(ctx context.Context) ([]*model.SystemConfig, error) {
	rows, err := s.DB.Query(ctx,
		"SELECT key, value, description, updated_at, updated_by FROM schema_admin.system_config ORDER BY key",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query config: %w", err)
	}
	defer rows.Close()

	var configs []*model.SystemConfig
	for rows.Next() {
		c := &model.SystemConfig{}
		if err := rows.Scan(&c.Key, &c.Value, &c.Description, &c.UpdatedAt, &c.UpdatedBy); err != nil {
			continue
		}
		configs = append(configs, c)
	}
	return configs, nil
}

// GetConfig returns a single config entry by key
func (s *ConfigService) GetConfig(ctx context.Context, key string) (*model.SystemConfig, error) {
	c := &model.SystemConfig{}
	err := s.DB.QueryRow(ctx,
		"SELECT key, value, description, updated_at, updated_by FROM schema_admin.system_config WHERE key = $1",
		key,
	).Scan(&c.Key, &c.Value, &c.Description, &c.UpdatedAt, &c.UpdatedBy)
	if err != nil {
		return nil, fmt.Errorf("config key not found: %s", key)
	}
	return c, nil
}

// UpdateConfig updates a config entry and logs the operation
func (s *ConfigService) UpdateConfig(ctx context.Context, key, value string, operatorID int64) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_admin.system_config SET value = $1, updated_at = NOW(), updated_by = $2 WHERE key = $3",
		value, operatorID, key,
	)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"new_value":   value,
		"operator_id": operatorID,
	})
	_, err = s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'update_config', 'system_config', 0, $3)`,
		operatorID, "", string(detail),
	)
	if err != nil {
		return fmt.Errorf("failed to log operation: %w", err)
	}
	return nil
}
