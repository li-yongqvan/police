package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ExtrasService handles notifications, reports, and public community stats.
type ExtrasService struct {
	DB *pgxpool.Pool
}

func NewExtrasService(db *pgxpool.Pool) *ExtrasService {
	return &ExtrasService{DB: db}
}

type NotificationItem struct {
	ID            uint   `json:"id"`
	Type          string `json:"type"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	RelatedPostID *uint  `json:"related_post_id,omitempty"`
	IsRead        bool   `json:"is_read"`
	CreatedAt     string `json:"created_at"`
}

func (s *ExtrasService) ListNotifications(ctx context.Context, userID uint, page, limit int) ([]NotificationItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit
	var total int
	if err := s.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM schema_forum.notifications WHERE user_id = $1`, userID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.DB.Query(ctx, `
		SELECT id, type, title, body, related_post_id, is_read, created_at::text
		FROM schema_forum.notifications WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []NotificationItem
	for rows.Next() {
		var n NotificationItem
		var postID *uint
		if err := rows.Scan(&n.ID, &n.Type, &n.Title, &n.Body, &postID, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, 0, err
		}
		n.RelatedPostID = postID
		items = append(items, n)
	}
	return items, total, nil
}

func (s *ExtrasService) CountUnreadNotifications(ctx context.Context, userID uint) (int, error) {
	var total int
	err := s.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM schema_forum.notifications WHERE user_id = $1 AND is_read = false`,
		userID,
	).Scan(&total)
	return total, err
}

func (s *ExtrasService) MarkNotificationRead(ctx context.Context, userID, id uint) error {
	_, err := s.DB.Exec(ctx,
		`UPDATE schema_forum.notifications SET is_read = true WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	return err
}

func (s *ExtrasService) CreateNotification(ctx context.Context, userID uint, nType, title, body string, postID *uint) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO schema_forum.notifications (user_id, type, title, body, related_post_id)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, nType, title, body, postID)
	return err
}

func (s *ExtrasService) ReportPost(ctx context.Context, reporterID, postID uint, reason string) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO schema_forum.post_reports (post_id, reporter_id, reason) VALUES ($1, $2, $3)
	`, postID, reporterID, reason)
	return err
}

type ReportRow struct {
	ID           uint   `json:"id"`
	PostID       uint   `json:"post_id"`
	PostTitle    string `json:"post_title"`
	ReporterID   uint   `json:"reporter_id"`
	ReporterName string `json:"reporter_name"`
	Reason       string `json:"reason"`
	Status       string `json:"status"`
	AdminNote    string `json:"admin_note,omitempty"`
	CreatedAt    string `json:"created_at"`
}

func (s *ExtrasService) ListReports(ctx context.Context, page, limit int, status string) ([]ReportRow, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if status == "" {
		status = "pending"
	}
	offset := (page - 1) * limit
	var total int
	if err := s.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM schema_forum.post_reports WHERE status = $1`, status,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.DB.Query(ctx, `
		SELECT r.id, r.post_id, COALESCE(p.title, ''), r.reporter_id,
		       COALESCE(u.nickname, u.username, ''), r.reason, r.status,
		       COALESCE(r.admin_note, ''), r.created_at::text
		FROM schema_forum.post_reports r
		LEFT JOIN schema_forum.posts p ON p.id = r.post_id
		LEFT JOIN schema_auth.users u ON u.id = r.reporter_id
		WHERE r.status = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []ReportRow
	for rows.Next() {
		var row ReportRow
		if err := rows.Scan(&row.ID, &row.PostID, &row.PostTitle, &row.ReporterID,
			&row.ReporterName, &row.Reason, &row.Status, &row.AdminNote, &row.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, row)
	}
	return items, total, nil
}

func (s *ExtrasService) ResolveReport(ctx context.Context, reportID uint, action string, deletePost bool, adminNote string) error {
	var postID uint
	var curStatus string
	err := s.DB.QueryRow(ctx,
		`SELECT post_id, status FROM schema_forum.post_reports WHERE id = $1`, reportID,
	).Scan(&postID, &curStatus)
	if err != nil {
		return fmt.Errorf("举报记录不存在")
	}
	if curStatus != "pending" {
		return fmt.Errorf("该举报已处理")
	}
	newStatus := "dismissed"
	if action == "resolved" {
		newStatus = "resolved"
	} else if action != "dismissed" {
		return fmt.Errorf("无效的处理动作")
	}
	_, err = s.DB.Exec(ctx, `
		UPDATE schema_forum.post_reports
		SET status = $1, admin_note = $2, resolved_at = NOW()
		WHERE id = $3
	`, newStatus, adminNote, reportID)
	if err != nil {
		return err
	}
	if deletePost && postID > 0 {
		_, _ = s.DB.Exec(ctx, `UPDATE schema_forum.posts SET status = 'deleted' WHERE id = $1`, postID)
	}
	return nil
}

type CommunityStats struct {
	TotalUsers  int `json:"total_users"`
	TotalPosts  int `json:"total_posts"`
	OnlineUsers int `json:"online_users"`
	PostsToday  int `json:"posts_today"`
}

func (s *ExtrasService) GetCommunityStats(ctx context.Context) (*CommunityStats, error) {
	stats := &CommunityStats{}
	_ = s.DB.QueryRow(ctx, `SELECT COUNT(*) FROM schema_auth.users WHERE status = 'active'`).Scan(&stats.TotalUsers)
	_ = s.DB.QueryRow(ctx, `SELECT COUNT(*) FROM schema_forum.posts WHERE status = 'published'`).Scan(&stats.TotalPosts)
	_ = s.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM schema_forum.posts
		WHERE status = 'published' AND created_at >= CURRENT_DATE
	`).Scan(&stats.PostsToday)
	stats.OnlineUsers = stats.TotalUsers
	if stats.OnlineUsers < 1 {
		stats.OnlineUsers = 0
	}
	return stats, nil
}

func (s *ExtrasService) NotifyPostAuthorOnComment(ctx context.Context, postID, authorID, commenterID uint, title string) {
	if authorID == commenterID {
		return
	}
	body := fmt.Sprintf("你的帖子《%s》收到新回复", title)
	_ = s.CreateNotification(ctx, authorID, "reply", "收到新回复", body, &postID)
}

func (s *ExtrasService) NotifyPostAuthorOnLike(ctx context.Context, postID, authorID, likerID uint, title string) {
	if authorID == likerID {
		return
	}
	body := fmt.Sprintf("你的帖子《%s》收到了新的点赞", title)
	_ = s.CreateNotification(ctx, authorID, "like", "收到新的点赞", body, &postID)
}

func (s *ExtrasService) NotifyCommentAuthorOnLike(ctx context.Context, postID, authorID, likerID uint, postTitle string) {
	if authorID == likerID {
		return
	}
	body := fmt.Sprintf("你在帖子《%s》下的评论收到了新的点赞", postTitle)
	_ = s.CreateNotification(ctx, authorID, "like", "评论收到点赞", body, &postID)
}
