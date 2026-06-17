package service

import (
	"context"
	"encoding/json"
	"fmt"

	"ai-forum/admin-service/internal/client"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostAdminService handles admin post management operations
type PostAdminService struct {
	ForumClient *client.ForumClient
	DB          *pgxpool.Pool
}

// NewPostAdminService creates a new PostAdminService
func NewPostAdminService(forumClient *client.ForumClient, db *pgxpool.Pool) *PostAdminService {
	return &PostAdminService{ForumClient: forumClient, DB: db}
}

// DeletePost deletes any post and logs the operation
func (s *PostAdminService) DeletePost(ctx context.Context, postID uint, operatorID uint, operatorName string) error {
	if err := s.ForumClient.DeletePost(postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	detail, _ := json.Marshal(map[string]interface{}{
		"operator_id":   operatorID,
		"operator_name": operatorName,
	})
	_, err := s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, 'delete_post', 'post', $3, $4)`,
		int64(operatorID), operatorName, int64(postID), string(detail),
	)
	return err
}

// SetPostFeatured sets/unsets featured flag and logs the operation
func (s *PostAdminService) SetPostFeatured(ctx context.Context, postID uint, featured bool, operatorID uint, operatorName string) error {
	if err := s.ForumClient.SetPostFeatured(postID, featured); err != nil {
		return fmt.Errorf("failed to set featured: %w", err)
	}

	action := "set_featured"
	if !featured {
		action = "unset_featured"
	}
	detail, _ := json.Marshal(map[string]interface{}{
		"operator_id":   operatorID,
		"operator_name": operatorName,
		"featured":      featured,
	})
	_, err := s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, $3, 'post', $4, $5)`,
		int64(operatorID), operatorName, action, int64(postID), string(detail),
	)
	return err
}

// SetPostPinned sets/unsets pinned flag and logs the operation
func (s *PostAdminService) SetPostPinned(ctx context.Context, postID uint, pinned bool, operatorID uint, operatorName string) error {
	if err := s.ForumClient.SetPostPinned(postID, pinned); err != nil {
		return fmt.Errorf("failed to set pinned: %w", err)
	}

	action := "set_pinned"
	if !pinned {
		action = "unset_pinned"
	}
	detail, _ := json.Marshal(map[string]interface{}{
		"operator_id":   operatorID,
		"operator_name": operatorName,
		"pinned":        pinned,
	})
	_, err := s.DB.Exec(ctx,
		`INSERT INTO schema_admin.operation_logs (operator_id, operator_username, action, target_type, target_id, detail)
		 VALUES ($1, $2, $3, 'post', $4, $5)`,
		int64(operatorID), operatorName, action, int64(postID), string(detail),
	)
	return err
}
