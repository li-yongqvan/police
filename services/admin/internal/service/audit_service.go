package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ai-forum/admin-service/internal/client"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PendingPost represents a post pending audit (response DTO)
type PendingPost struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	AuthorName   string    `json:"author_name"`
	BoardName    string    `json:"board_name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	Content      string    `json:"content"`
	AuthorID     uint      `json:"author_id"`
	BoardID      uint      `json:"board_id"`
	MatchedWords []string  `json:"matched_words,omitempty"`
}

// AuditService handles audit workflow operations
type AuditService struct {
	DB          *pgxpool.Pool
	ForumClient *client.ForumClient
}

// NewAuditService creates a new AuditService
func NewAuditService(db *pgxpool.Pool, forumClient *client.ForumClient) *AuditService {
	return &AuditService{DB: db, ForumClient: forumClient}
}

// ListPendingAudit returns posts pending review
func (s *AuditService) ListPendingAudit(ctx context.Context, page, limit int) ([]PendingPost, int, error) {
	clientPosts, total, err := s.ForumClient.ListPendingPosts(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list pending posts: %w", err)
	}

	var posts []PendingPost
	for _, p := range clientPosts {
		posts = append(posts, PendingPost{
			ID:           p.ID,
			Title:        p.Title,
			AuthorName:   p.AuthorName,
			BoardName:    p.BoardName,
			Status:       "pending_review",
			CreatedAt:    p.CreatedAt,
			Content:      p.Content,
			AuthorID:     p.AuthorID,
			BoardID:      p.BoardID,
			MatchedWords: p.MatchedWords,
		})
	}
	return posts, total, nil
}

// ApprovePost approves a pending post and logs the operation
func (s *AuditService) ApprovePost(ctx context.Context, postID uint, operatorID uint, operatorName string) error {
	if err := s.ForumClient.ChangePostStatus(postID, "approved"); err != nil {
		return fmt.Errorf("failed to approve post: %w", err)
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"operator_id":   operatorID,
		"operator_name": operatorName,
	})
	_, err := s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'approve_post', 'post', $3, $4)`,
		int64(operatorID), operatorName, int64(postID), string(detail),
	)
	if err != nil {
		return fmt.Errorf("failed to log operation: %w", err)
	}
	return nil
}

// RejectPost rejects a pending post and logs the operation
func (s *AuditService) RejectPost(ctx context.Context, postID uint, operatorID uint, operatorName string, reason string) error {
	if err := s.ForumClient.ChangePostStatus(postID, "rejected"); err != nil {
		return fmt.Errorf("failed to reject post: %w", err)
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"operator_id":   operatorID,
		"operator_name": operatorName,
		"reason":        reason,
	})
	_, err := s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'reject_post', 'post', $3, $4)`,
		int64(operatorID), operatorName, int64(postID), string(detail),
	)
	if err != nil {
		return fmt.Errorf("failed to log operation: %w", err)
	}
	return nil
}

// BatchDeletePosts batch deletes posts and logs each operation
func (s *AuditService) BatchDeletePosts(ctx context.Context, postIDs []uint, operatorID uint, operatorName string, reason string) error {
	if err := s.ForumClient.BatchDeletePosts(postIDs); err != nil {
		return fmt.Errorf("failed to batch delete posts: %w", err)
	}

	for _, postID := range postIDs {
		detail, _ := json.Marshal(map[string]interface{}{
			"operator_id":   operatorID,
			"operator_name": operatorName,
			"reason":        reason,
		})
		_, err := s.DB.Exec(ctx,
			`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
			 VALUES ($1, $2, 'batch_delete_post', 'post', $3, $4)`,
			int64(operatorID), operatorName, int64(postID), string(detail),
		)
		if err != nil {
			return fmt.Errorf("failed to log operation: %w", err)
		}
	}
	return nil
}
