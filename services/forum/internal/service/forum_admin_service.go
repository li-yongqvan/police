package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PendingPostRow represents a post pending review (from DB query)
type PendingPostRow struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	AuthorID     uint      `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	BoardID      uint      `json:"board_id"`
	BoardName    string    `json:"board_name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	MatchedWords []string  `json:"matched_words"`
}

// ForumAdminService handles admin operations on forum-service
type ForumAdminService struct {
	DB *pgxpool.Pool
}

// NewForumAdminService creates a new ForumAdminService
func NewForumAdminService(db *pgxpool.Pool) *ForumAdminService {
	return &ForumAdminService{DB: db}
}

// ListPendingPosts returns posts with status='pending_review', paginated
func (s *ForumAdminService) ListPendingPosts(ctx context.Context, page, limit int) ([]PendingPostRow, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int
	err := s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_forum.posts WHERE status = 'pending_review'",
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pending posts: %w", err)
	}

	query := `
		SELECT p.id, p.title, p.content, p.author_id, u.username, p.board_id, b.name, p.status, p.matched_words, p.created_at
		FROM schema_forum.posts p
		JOIN schema_auth.users u ON p.author_id = u.id
		JOIN schema_forum.boards b ON p.board_id = b.id
		WHERE p.status = 'pending_review'
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := s.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pending posts: %w", err)
	}
	defer rows.Close()

	var posts []PendingPostRow
	for rows.Next() {
		var p PendingPostRow
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.BoardID, &p.BoardName, &p.Status, &p.MatchedWords, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}
	return posts, total, nil
}

// ApprovePost sets post status to 'published'
func (s *ForumAdminService) ApprovePost(ctx context.Context, postID uint) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.posts SET status = 'published', updated_at = NOW() WHERE id = $1",
		postID,
	)
	if err != nil {
		return fmt.Errorf("failed to approve post: %w", err)
	}
	return nil
}

// RejectPost sets post status to 'rejected'
func (s *ForumAdminService) RejectPost(ctx context.Context, postID uint) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.posts SET status = 'rejected', updated_at = NOW() WHERE id = $1",
		postID,
	)
	if err != nil {
		return fmt.Errorf("failed to reject post: %w", err)
	}
	return nil
}

// AdminDeletePost sets post status to 'deleted' (admin action, no author check)
func (s *ForumAdminService) AdminDeletePost(ctx context.Context, postID uint) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.posts SET status = 'deleted', updated_at = NOW() WHERE id = $1",
		postID,
	)
	if err != nil {
		return fmt.Errorf("failed to admin-delete post: %w", err)
	}
	return nil
}

// SetPostFeatured sets the is_featured flag on a post
func (s *ForumAdminService) SetPostFeatured(ctx context.Context, postID uint, featured bool) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.posts SET is_featured = $1, updated_at = NOW() WHERE id = $2",
		featured, postID,
	)
	if err != nil {
		return fmt.Errorf("failed to set post featured: %w", err)
	}
	return nil
}

// SetPostPinned sets the is_pinned flag on a post
func (s *ForumAdminService) SetPostPinned(ctx context.Context, postID uint, pinned bool) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.posts SET is_pinned = $1, updated_at = NOW() WHERE id = $2",
		pinned, postID,
	)
	if err != nil {
		return fmt.Errorf("failed to set post pinned: %w", err)
	}
	return nil
}

// ChangePostStatus sets the status of a post to the given value
func (s *ForumAdminService) ChangePostStatus(ctx context.Context, postID uint, status string) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.posts SET status = $1, updated_at = NOW() WHERE id = $2",
		status, postID,
	)
	if err != nil {
		return fmt.Errorf("failed to change post status: %w", err)
	}
	return nil
}

// BatchDeletePosts sets status='deleted' for multiple posts
func (s *ForumAdminService) BatchDeletePosts(ctx context.Context, postIDs []uint) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.posts SET status = 'deleted', updated_at = NOW() WHERE id = ANY($1)",
		postIDs,
	)
	if err != nil {
		return fmt.Errorf("failed to batch delete posts: %w", err)
	}
	return nil
}

// CreateBoard inserts a new board into schema_forum.boards
func (s *ForumAdminService) CreateBoard(ctx context.Context, name, slug, description string, sortOrder int) error {
	_, err := s.DB.Exec(ctx,
		"INSERT INTO schema_forum.boards (name, slug, description, enabled, sort_order) VALUES ($1, $2, $3, true, $4)",
		name, slug, description, sortOrder,
	)
	if err != nil {
		return fmt.Errorf("failed to create board: %w", err)
	}
	return nil
}

// UpdateBoard updates a board's attributes
func (s *ForumAdminService) UpdateBoard(ctx context.Context, id uint, name, slug, description string, enabled bool, sortOrder int) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.boards SET name = $1, slug = $2, description = $3, enabled = $4, sort_order = $5 WHERE id = $6",
		name, slug, description, enabled, sortOrder, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update board: %w", err)
	}
	return nil
}

// DeleteBoard soft-deletes a board (sets enabled=false)
func (s *ForumAdminService) DeleteBoard(ctx context.Context, id uint) error {
	_, err := s.DB.Exec(ctx,
		"UPDATE schema_forum.boards SET enabled = false WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}
	return nil
}

// ListAllBoards returns all boards including disabled ones
func (s *ForumAdminService) ListAllBoards(ctx context.Context) ([]BoardInfo, error) {
	rows, err := s.DB.Query(ctx,
		"SELECT id, name, slug, description, enabled, sort_order, created_at FROM schema_forum.boards ORDER BY sort_order, id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query boards: %w", err)
	}
	defer rows.Close()

	var boards []BoardInfo
	for rows.Next() {
		var b BoardInfo
		if err := rows.Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.Enabled, &b.SortOrder, &b.CreatedAt); err != nil {
			continue
		}
		boards = append(boards, b)
	}
	return boards, nil
}

// AdminPostRow represents a post returned by the admin list endpoint
type AdminPostRow struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorID   uint      `json:"author_id"`
	AuthorName string    `json:"author_name"`
	BoardID    uint      `json:"board_id"`
	BoardName  string    `json:"board_name"`
	Status     string    `json:"status"`
	IsFeatured bool      `json:"is_featured"`
	IsPinned   bool      `json:"is_pinned"`
	CreatedAt  time.Time `json:"created_at"`
}

// ListAllPosts returns all posts (paginated) for admin management
func (s *ForumAdminService) ListAllPosts(ctx context.Context, page, limit int) ([]AdminPostRow, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int
	err := s.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM schema_forum.posts WHERE status != 'deleted'",
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	query := `
		SELECT p.id, p.title, p.content, p.author_id, u.username, p.board_id, b.name,
		       p.status, p.is_featured, p.is_pinned, p.created_at
		FROM schema_forum.posts p
		JOIN schema_auth.users u ON p.author_id = u.id
		JOIN schema_forum.boards b ON p.board_id = b.id
		WHERE p.status != 'deleted'
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := s.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []AdminPostRow
	for rows.Next() {
		var p AdminPostRow
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.BoardID, &p.BoardName, &p.Status, &p.IsFeatured, &p.IsPinned, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}
	return posts, total, nil
}

// BoardInfo represents a board (mirrors model.Board)
type BoardInfo struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}
