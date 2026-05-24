package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserStatsService handles statistics queries for user data
type UserStatsService struct {
	DB *pgxpool.Pool
}

// NewUserStatsService creates a new UserStatsService
func NewUserStatsService(db *pgxpool.Pool) *UserStatsService {
	return &UserStatsService{DB: db}
}

// UserOverview contains aggregate user statistics
type UserOverview struct {
	TotalUsers  int64 `json:"total_users"`
	UsersToday  int64 `json:"users_today"`
	BannedUsers int64 `json:"banned_users"`
}

// LevelDistribution represents user count at a given level
type LevelDistribution struct {
	Level int   `json:"level"`
	Count int64 `json:"count"`
}

// DailyUserCount represents a date-count pair for user registration
type DailyUserCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// GetOverview returns aggregate user statistics
func (s *UserStatsService) GetOverview(ctx context.Context) (*UserOverview, error) {
	overview := &UserOverview{}
	today := time.Now().Truncate(24 * time.Hour)

	err := s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_auth.users",
	).Scan(&overview.TotalUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}

	err = s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_auth.users WHERE created_at >= $1",
		today,
	).Scan(&overview.UsersToday)
	if err != nil {
		return nil, fmt.Errorf("failed to count today's users: %w", err)
	}

	err = s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_auth.users WHERE status = 'banned'",
	).Scan(&overview.BannedUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to count banned users: %w", err)
	}

	return overview, nil
}

// GetLevelDistribution returns user count per level
func (s *UserStatsService) GetLevelDistribution(ctx context.Context) ([]LevelDistribution, error) {
	rows, err := s.DB.Query(ctx,
		"SELECT level, COUNT(*) as cnt FROM schema_auth.users GROUP BY level ORDER BY level",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query level distribution: %w", err)
	}
	defer rows.Close()

	var result []LevelDistribution
	for rows.Next() {
		var ld LevelDistribution
		if err := rows.Scan(&ld.Level, &ld.Count); err != nil {
			continue
		}
		result = append(result, ld)
	}
	return result, nil
}

// GetDailyUsers returns user registration counts for the last N days
func (s *UserStatsService) GetDailyUsers(ctx context.Context, days int) ([]DailyUserCount, error) {
	if days < 1 {
		days = 7
	}
	cutoff := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -days+1)

	query := `
		SELECT DATE(created_at)::text as day, COUNT(*) as cnt
		FROM schema_auth.users
		WHERE created_at >= $1
		GROUP BY DATE(created_at)
		ORDER BY day
	`
	rows, err := s.DB.Query(ctx, query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily users: %w", err)
	}
	defer rows.Close()

	var result []DailyUserCount
	for rows.Next() {
		var dc DailyUserCount
		if err := rows.Scan(&dc.Date, &dc.Count); err != nil {
			continue
		}
		result = append(result, dc)
	}
	return result, nil
}
