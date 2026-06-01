package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// StatsService handles statistics queries for forum data
type StatsService struct {
	DB *pgxpool.Pool
}

// NewStatsService creates a new StatsService
func NewStatsService(db *pgxpool.Pool) *StatsService {
	return &StatsService{DB: db}
}

// ForumOverview contains aggregate forum statistics
type ForumOverview struct {
	TotalPosts    int64 `json:"total_posts"`
	TotalComments int64 `json:"total_comments"`
	TotalLikes    int64 `json:"total_likes"`
	PostsToday    int64 `json:"posts_today"`
	CommentsToday int64 `json:"comments_today"`
}

// DailyCount represents a date-count pair
type DailyCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// BoardActivity represents activity stats for a single board
type BoardActivity struct {
	BoardID      uint   `json:"board_id"`
	BoardName    string `json:"board_name"`
	PostCount    int64  `json:"post_count"`
	CommentCount int64  `json:"comment_count"`
}

// GetOverview returns aggregate forum statistics
func (s *StatsService) GetOverview(ctx context.Context) (*ForumOverview, error) {
	overview := &ForumOverview{}
	today := time.Now().Truncate(24 * time.Hour)

	// Total posts (published)
	err := s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_forum.posts WHERE status = 'published'",
	).Scan(&overview.TotalPosts)
	if err != nil {
		return nil, fmt.Errorf("failed to count total posts: %w", err)
	}

	// Total comments (comments table has no status column)
	err = s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_forum.comments",
	).Scan(&overview.TotalComments)
	if err != nil {
		return nil, fmt.Errorf("failed to count total comments: %w", err)
	}

	// Total likes
	err = s.DB.QueryRow(ctx,
		"SELECT COALESCE(SUM(like_count), 0) FROM schema_forum.posts WHERE status = 'published'",
	).Scan(&overview.TotalLikes)
	if err != nil {
		return nil, fmt.Errorf("failed to sum likes: %w", err)
	}

	// Posts today
	err = s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_forum.posts WHERE status = 'published' AND created_at >= $1",
		today,
	).Scan(&overview.PostsToday)
	if err != nil {
		return nil, fmt.Errorf("failed to count today's posts: %w", err)
	}

	// Comments today
	err = s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_forum.comments WHERE created_at >= $1",
		today,
	).Scan(&overview.CommentsToday)
	if err != nil {
		return nil, fmt.Errorf("failed to count today's comments: %w", err)
	}

	return overview, nil
}

// GetDailyPosts returns post counts for the last N days
func (s *StatsService) GetDailyPosts(ctx context.Context, days int) ([]DailyCount, error) {
	if days < 1 {
		days = 7
	}
	cutoff := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -days+1)

	query := `
		SELECT DATE(created_at)::text as day, COUNT(*) as cnt
		FROM schema_forum.posts
		WHERE status = 'published' AND created_at >= $1
		GROUP BY DATE(created_at)
		ORDER BY day
	`
	rows, err := s.DB.Query(ctx, query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily posts: %w", err)
	}
	defer rows.Close()

	var result []DailyCount
	for rows.Next() {
		var dc DailyCount
		if err := rows.Scan(&dc.Date, &dc.Count); err != nil {
			continue
		}
		result = append(result, dc)
	}
	return result, nil
}

// GetDailyComments returns comment counts for the last N days
func (s *StatsService) GetDailyComments(ctx context.Context, days int) ([]DailyCount, error) {
	if days < 1 {
		days = 7
	}
	cutoff := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -days+1)

	query := `
		SELECT DATE(created_at)::text as day, COUNT(*) as cnt
		FROM schema_forum.comments
		WHERE created_at >= $1
		GROUP BY DATE(created_at)
		ORDER BY day
	`
	rows, err := s.DB.Query(ctx, query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily comments: %w", err)
	}
	defer rows.Close()

	var result []DailyCount
	for rows.Next() {
		var dc DailyCount
		if err := rows.Scan(&dc.Date, &dc.Count); err != nil {
			continue
		}
		result = append(result, dc)
	}
	return result, nil
}

// GetBoardActivity returns activity stats per board
func (s *StatsService) GetBoardActivity(ctx context.Context) ([]BoardActivity, error) {
	query := `
		SELECT b.id, b.name,
		       COUNT(DISTINCT p.id) FILTER (WHERE p.status = 'published') as post_count,
		       COUNT(DISTINCT c.id) as comment_count
		FROM schema_forum.boards b
		LEFT JOIN schema_forum.posts p ON p.board_id = b.id
		LEFT JOIN schema_forum.comments c ON c.post_id = p.id
		WHERE b.enabled = true
		GROUP BY b.id, b.name
		ORDER BY post_count DESC
	`
	rows, err := s.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query board activity: %w", err)
	}
	defer rows.Close()

	var result []BoardActivity
	for rows.Next() {
		var ba BoardActivity
		if err := rows.Scan(&ba.BoardID, &ba.BoardName, &ba.PostCount, &ba.CommentCount); err != nil {
			continue
		}
		result = append(result, ba)
	}
	return result, nil
}
