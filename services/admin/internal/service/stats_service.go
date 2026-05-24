package service

import (
	"context"
	"fmt"
	"time"

	"ai-forum/admin-service/internal/client"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserOverview contains aggregate user statistics (mirrors user-service)
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

// DailyCount represents a date-count pair for forum data
type DailyCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// BoardActivity represents activity stats per board
type BoardActivity struct {
	BoardID      uint   `json:"board_id"`
	BoardName    string `json:"board_name"`
	PostCount    int64  `json:"post_count"`
	CommentCount int64  `json:"comment_count"`
}

// DailyStatsRow represents a row from statistics_daily
type DailyStatsRow struct {
	StatDate      string                   `json:"stat_date"`
	NewUsers      int                      `json:"new_users"`
	NewPosts      int                      `json:"new_posts"`
	NewComments   int                      `json:"new_comments"`
	ActiveUsers   int                      `json:"active_users"`
	TotalUsers    int                      `json:"total_users"`
	TotalPosts    int                      `json:"total_posts"`
	TotalComments int                      `json:"total_comments"`
	BoardActivity []map[string]interface{} `json:"board_activity"`
}

// OverviewData combines user and forum stats for admin dashboard
type OverviewData struct {
	TotalUsers    int64 `json:"total_users"`
	UsersToday    int64 `json:"users_today"`
	TotalPosts    int64 `json:"total_posts"`
	PostsToday    int64 `json:"posts_today"`
	TotalComments int64 `json:"total_comments"`
	BannedUsers   int64 `json:"banned_users"`
}

// StatsService handles aggregated statistics for admin dashboard
type StatsService struct {
	UserClient  *UserClient
	ForumClient *client.ForumClient
	DB          *pgxpool.Pool
}

// NewStatsService creates a new StatsService
func NewStatsService(userClient *UserClient, forumClient *client.ForumClient, db *pgxpool.Pool) *StatsService {
	return &StatsService{
		UserClient:  userClient,
		ForumClient: forumClient,
		DB:          db,
	}
}

// GetOverview aggregates stats from user-service and forum-service
func (s *StatsService) GetOverview(ctx context.Context) (*OverviewData, error) {
	userOverview, err := s.UserClient.GetUserOverview()
	if err != nil {
		return nil, fmt.Errorf("failed to get user overview: %w", err)
	}

	forumOverview, err := s.ForumClient.GetForumOverview()
	if err != nil {
		return nil, fmt.Errorf("failed to get forum overview: %w", err)
	}

	return &OverviewData{
		TotalUsers:    userOverview.TotalUsers,
		UsersToday:    userOverview.UsersToday,
		TotalPosts:    forumOverview.TotalPosts,
		PostsToday:    forumOverview.PostsToday,
		TotalComments: forumOverview.TotalComments,
		BannedUsers:   userOverview.BannedUsers,
	}, nil
}

// GetDailyStats returns historical daily stats from DB, or aggregates live from services if not available
func (s *StatsService) GetDailyStats(ctx context.Context, days int) ([]DailyStatsRow, error) {
	if days < 1 || days > 30 {
		days = 7
	}

	// Try to get from DB first
	rows, err := s.DB.Query(ctx, `
		SELECT stat_date::text, new_users, new_posts, new_comments, active_users,
		       total_users, total_posts, total_comments, board_activity
		FROM schema_admin.statistics_daily
		WHERE stat_date >= CURRENT_DATE - ($1 || ' days')::interval
		ORDER BY stat_date
	`, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily stats: %w", err)
	}
	defer rows.Close()

	var result []DailyStatsRow
	for rows.Next() {
		var r DailyStatsRow
		if err := rows.Scan(&r.StatDate, &r.NewUsers, &r.NewPosts, &r.NewComments, &r.ActiveUsers, &r.TotalUsers, &r.TotalPosts, &r.TotalComments, &r.BoardActivity); err != nil {
			continue
		}
		result = append(result, r)
	}

	if len(result) > 0 {
		return result, nil
	}

	// If no data in DB, aggregate live from services
	return s.aggregateLiveStats(ctx, days)
}

func (s *StatsService) aggregateLiveStats(ctx context.Context, days int) ([]DailyStatsRow, error) {
	userOverview, uErr := s.UserClient.GetUserOverview()
	forumOverview, fErr := s.ForumClient.GetForumOverview()
	dailyUsers, _ := s.UserClient.GetDailyUsers(days)
	dailyPosts, _ := s.ForumClient.GetDailyPosts(days)
	boardActivity, _ := s.ForumClient.GetBoardActivity()

	resultMap := make(map[string]*DailyStatsRow)
	for _, du := range dailyUsers {
		resultMap[du.Date] = &DailyStatsRow{StatDate: du.Date, NewUsers: int(du.Count)}
	}
	for _, dp := range dailyPosts {
		if r, ok := resultMap[dp.Date]; ok {
			r.NewPosts = int(dp.Count)
		} else {
			resultMap[dp.Date] = &DailyStatsRow{StatDate: dp.Date, NewPosts: int(dp.Count)}
		}
	}

	var result []DailyStatsRow
	for _, r := range resultMap {
		result = append(result, *r)
	}

	// Add current overview as latest day
	if uErr == nil && fErr == nil {
		var boardJSON []map[string]interface{}
		for _, ba := range boardActivity {
			boardJSON = append(boardJSON, map[string]interface{}{
				"board_id":      ba.BoardID,
				"board_name":    ba.BoardName,
				"post_count":    ba.PostCount,
				"comment_count": ba.CommentCount,
			})
		}
		result = append(result, DailyStatsRow{
			StatDate:      time.Now().Format("2006-01-02"),
			NewUsers:      int(userOverview.UsersToday),
			TotalUsers:    int(userOverview.TotalUsers),
			NewPosts:      int(forumOverview.PostsToday),
			TotalPosts:    int(forumOverview.TotalPosts),
			NewComments:   int(forumOverview.CommentsToday),
			TotalComments: int(forumOverview.TotalComments),
			BoardActivity: boardJSON,
		})
	}

	return result, nil
}

// ComputeDailyStats computes stats for a given date and stores them in the DB
func (s *StatsService) ComputeDailyStats(ctx context.Context, targetDate time.Time) error {
	dateStr := targetDate.Format("2006-01-02")
	nextDay := targetDate.AddDate(0, 0, 1)

	// Get counts from forum-service schema
	var newPosts, newComments int
	err := s.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM schema_forum.posts
		WHERE status = 'published' AND created_at >= $1 AND created_at < $2
	`, targetDate, nextDay).Scan(&newPosts)
	if err != nil {
		return fmt.Errorf("failed to count new posts: %w", err)
	}

	err = s.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM schema_forum.comments
		WHERE status = 'published' AND created_at >= $1 AND created_at < $2
	`, targetDate, nextDay).Scan(&newComments)
	if err != nil {
		return fmt.Errorf("failed to count new comments: %w", err)
	}

	// Get counts from user-service schema
	var newUsers int
	err = s.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM schema_auth.users
		WHERE created_at >= $1 AND created_at < $2
	`, targetDate, nextDay).Scan(&newUsers)
	if err != nil {
		return fmt.Errorf("failed to count new users: %w", err)
	}

	// Get totals
	var totalUsers, totalPosts, totalComments int
	s.DB.QueryRow(ctx, "SELECT COUNT(*) FROM schema_auth.users").Scan(&totalUsers)
	s.DB.QueryRow(ctx, "SELECT COUNT(*) FROM schema_forum.posts WHERE status = 'published'").Scan(&totalPosts)
	s.DB.QueryRow(ctx, "SELECT COUNT(*) FROM schema_forum.comments WHERE status = 'published'").Scan(&totalComments)

	// Insert or update
	_, err = s.DB.Exec(ctx, `
		INSERT INTO schema_admin.statistics_daily
			(stat_date, new_users, new_posts, new_comments, total_users, total_posts, total_comments)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (stat_date) DO UPDATE SET
			new_users = EXCLUDED.new_users,
			new_posts = EXCLUDED.new_posts,
			new_comments = EXCLUDED.new_comments,
			total_users = EXCLUDED.total_users,
			total_posts = EXCLUDED.total_posts,
			total_comments = EXCLUDED.total_comments
	`, dateStr, newUsers, newPosts, newComments, totalUsers, totalPosts, totalComments)

	if err != nil {
		return fmt.Errorf("failed to upsert daily stats: %w", err)
	}

	return nil
}
