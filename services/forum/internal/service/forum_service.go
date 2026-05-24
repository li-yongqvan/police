package service

import (
	"context"
	"fmt"

	"ai-forum/forum-service/internal/client"
	"ai-forum/forum-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ForumService handles forum business logic
type ForumService struct {
	DB          *pgxpool.Pool
	AdminClient *client.AdminClient
}

// NewForumService creates a new ForumService instance
func NewForumService(db *pgxpool.Pool, adminClient *client.AdminClient) *ForumService {
	return &ForumService{DB: db, AdminClient: adminClient}
}

// ListBoards returns all enabled boards
func (s *ForumService) ListBoards(ctx context.Context) ([]*model.BoardResponse, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT b.id, b.name, b.slug, b.description, b.enabled,
		       COALESCE(p.post_count, 0) as post_count
		FROM schema_forum.boards b
		LEFT JOIN (
			SELECT board_id, COUNT(*) as post_count
			FROM schema_forum.posts
			WHERE status = 'published'
			GROUP BY board_id
		) p ON b.id = p.board_id
		WHERE b.enabled = true
		ORDER BY b.sort_order, b.id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query boards: %w", err)
	}
	defer rows.Close()

	var boards []*model.BoardResponse
	for rows.Next() {
		b := &model.BoardResponse{}
		if err := rows.Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.Enabled, &b.PostCount); err != nil {
			return nil, fmt.Errorf("failed to scan board: %w", err)
		}
		boards = append(boards, b)
	}
	return boards, nil
}

// GetBoard returns a single board by ID
func (s *ForumService) GetBoard(ctx context.Context, id uint) (*model.BoardResponse, error) {
	b := &model.BoardResponse{}
	err := s.DB.QueryRow(ctx, `
		SELECT b.id, b.name, b.slug, b.description, b.enabled,
		       COALESCE(p.post_count, 0) as post_count
		FROM schema_forum.boards b
		LEFT JOIN (
			SELECT board_id, COUNT(*) as post_count
			FROM schema_forum.posts
			WHERE status = 'published' AND board_id = $1
			GROUP BY board_id
		) p ON b.id = p.board_id
		WHERE b.id = $1
	`, id).Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.Enabled, &b.PostCount)
	if err != nil {
		return nil, fmt.Errorf("board not found")
	}
	return b, nil
}

// ListPosts returns paginated posts, optionally filtered by board
func (s *ForumService) ListPosts(ctx context.Context, boardID uint, page, limit int) ([]*model.PostListItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Count total
	countQuery := "SELECT COUNT(*) FROM schema_forum.posts WHERE status = 'published'"
	args := []interface{}{}
	if boardID > 0 {
		countQuery += " AND board_id = $1"
		args = append(args, boardID)
	}
	var total int
	err := s.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Fetch posts
	query := `
		SELECT p.id, p.title, p.author_id, u.username, p.board_id, b.name,
		       p.status, p.is_pinned, p.is_featured, p.like_count, p.comment_count, p.created_at
		FROM schema_forum.posts p
		JOIN schema_auth.users u ON p.author_id = u.id
		JOIN schema_forum.boards b ON p.board_id = b.id
		WHERE p.status = 'published'`
	if boardID > 0 {
		query += " AND p.board_id = $1"
	}
	query += " ORDER BY p.is_pinned DESC, p.created_at DESC"
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*model.PostListItem
	for rows.Next() {
		p := &model.PostListItem{}
		if err := rows.Scan(&p.ID, &p.Title, &p.AuthorID, &p.AuthorName,
			&p.BoardID, &p.BoardName, &p.Status, &p.IsPinned, &p.IsFeatured,
			&p.LikeCount, &p.CommentCount, &p.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, total, nil
}

// GetPost returns full post detail with attachments
func (s *ForumService) GetPost(ctx context.Context, id uint) (*model.PostDetail, error) {
	detail := &model.PostDetail{}
	err := s.DB.QueryRow(ctx, `
		SELECT p.id, p.title, p.content, p.author_id, u.username, p.board_id, b.name,
		       p.status, p.is_pinned, p.is_featured, p.like_count, p.comment_count,
		       p.created_at, p.updated_at
		FROM schema_forum.posts p
		JOIN schema_auth.users u ON p.author_id = u.id
		JOIN schema_forum.boards b ON p.board_id = b.id
		WHERE p.id = $1 AND p.status IN ('published', 'pending_review')
	`, id).Scan(
		&detail.ID, &detail.Title, &detail.Content, &detail.AuthorID, &detail.AuthorName,
		&detail.BoardID, &detail.BoardName, &detail.Status, &detail.IsPinned, &detail.IsFeatured,
		&detail.LikeCount, &detail.CommentCount, &detail.CreatedAt, &detail.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	// Fetch attachments
	rows, err := s.DB.Query(ctx,
		"SELECT id, COALESCE(post_id, 0), COALESCE(comment_id, 0), user_id, filename, file_type, file_path, file_size, created_at FROM schema_forum.attachments WHERE post_id = $1",
		id,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			a := model.Attachment{}
			var pid, cid uint
			if err := rows.Scan(&a.ID, &pid, &cid, &a.UserID, &a.Filename, &a.FileType, &a.FilePath, &a.FileSize, &a.CreatedAt); err == nil {
				detail.Attachments = append(detail.Attachments, a)
			}
		}
	}

	return detail, nil
}

// CreatePost creates a new post with sensitive word check
func (s *ForumService) CreatePost(ctx context.Context, authorID uint, req *model.CreatePostRequest) (*model.Post, error) {
	// Check sensitive words
	clean, err := s.AdminClient.CheckSensitiveWords(req.Title + " " + req.Content)
	if err != nil {
		// If admin service is unavailable, allow post through
		clean = true
	}

	status := "published"
	if !clean {
		status = "pending_review"
	}

	post := &model.Post{}
	err = s.DB.QueryRow(ctx, `
		INSERT INTO schema_forum.posts (title, content, author_id, board_id, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, content, author_id, board_id, status, is_pinned, is_featured, like_count, comment_count, created_at, updated_at
	`, req.Title, req.Content, authorID, req.BoardID, status).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.BoardID,
		&post.Status, &post.IsPinned, &post.IsFeatured, &post.LikeCount,
		&post.CommentCount, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Link attachments if provided
	if len(req.AttachmentIDs) > 0 {
		_, err = s.DB.Exec(ctx,
			"UPDATE schema_forum.attachments SET post_id = $1 WHERE id = ANY($2) AND user_id = $3",
			post.ID, req.AttachmentIDs, authorID,
		)
		if err != nil {
			// Non-fatal, log warning
		}
	}

	return post, nil
}

// UpdatePost updates an existing post (author verification required)
func (s *ForumService) UpdatePost(ctx context.Context, authorID, postID uint, req *model.UpdatePostRequest) (*model.Post, error) {
	// Verify ownership
	var existing model.Post
	err := s.DB.QueryRow(ctx,
		"SELECT id, title, content, author_id, board_id, status, is_pinned, is_featured, like_count, comment_count, created_at, updated_at FROM schema_forum.posts WHERE id = $1",
		postID,
	).Scan(&existing.ID, &existing.Title, &existing.Content, &existing.AuthorID,
		&existing.BoardID, &existing.Status, &existing.IsPinned, &existing.IsFeatured,
		&existing.LikeCount, &existing.CommentCount, &existing.CreatedAt, &existing.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}
	if existing.AuthorID != authorID {
		return nil, fmt.Errorf("无权修改他人帖子")
	}

	// Check sensitive words on updated content
	title := req.Title
	if title == "" {
		title = existing.Title
	}
	content := req.Content
	if content == "" {
		content = existing.Content
	}
	clean, _ := s.AdminClient.CheckSensitiveWords(title + " " + content)
	newStatus := existing.Status
	if !clean {
		newStatus = "pending_review"
	}

	err = s.DB.QueryRow(ctx, `
		UPDATE schema_forum.posts
		SET title = COALESCE(NULLIF($1, ''), title),
		    content = COALESCE(NULLIF($2, ''), content),
		    status = $3,
		    updated_at = NOW()
		WHERE id = $4
		RETURNING id, title, content, author_id, board_id, status, is_pinned, is_featured, like_count, comment_count, created_at, updated_at
	`, title, content, newStatus, postID).Scan(
		&existing.ID, &existing.Title, &existing.Content, &existing.AuthorID,
		&existing.BoardID, &existing.Status, &existing.IsPinned, &existing.IsFeatured,
		&existing.LikeCount, &existing.CommentCount, &existing.CreatedAt, &existing.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}
	return &existing, nil
}

// DeletePost soft-deletes a post (author verification required)
func (s *ForumService) DeletePost(ctx context.Context, authorID, postID uint) error {
	var existingAuthor uint
	err := s.DB.QueryRow(ctx, "SELECT author_id FROM schema_forum.posts WHERE id = $1", postID).Scan(&existingAuthor)
	if err != nil {
		return fmt.Errorf("post not found")
	}
	if existingAuthor != authorID {
		return fmt.Errorf("无权删除他人帖子")
	}

	_, err = s.DB.Exec(ctx,
		"UPDATE schema_forum.posts SET status = 'deleted', updated_at = NOW() WHERE id = $1",
		postID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

// ListComments returns paginated comments for a post
func (s *ForumService) ListComments(ctx context.Context, postID uint, page, limit int) ([]*model.Comment, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	offset := (page - 1) * limit

	var total int
	err := s.DB.QueryRow(ctx, "SELECT COUNT(*) FROM schema_forum.comments WHERE post_id = $1", postID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.DB.Query(ctx,
		"SELECT id, post_id, author_id, content, created_at FROM schema_forum.comments WHERE post_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3",
		postID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Content, &c.CreatedAt); err != nil {
			return nil, 0, err
		}
		comments = append(comments, c)
	}
	return comments, total, nil
}

// CreateComment creates a new comment with sensitive word check
func (s *ForumService) CreateComment(ctx context.Context, authorID, postID uint, content string) (*model.Comment, error) {
	// Check sensitive words
	clean, err := s.AdminClient.CheckSensitiveWords(content)
	if err != nil {
		clean = true
	}
	if !clean {
		return nil, fmt.Errorf("comment contains sensitive words, rejected")
	}

	comment := &model.Comment{}
	err = s.DB.QueryRow(ctx, `
		INSERT INTO schema_forum.comments (post_id, author_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, post_id, author_id, content, created_at
	`, postID, authorID, content).Scan(
		&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Increment comment count
	_, _ = s.DB.Exec(ctx, "UPDATE schema_forum.posts SET comment_count = comment_count + 1 WHERE id = $1", postID)

	return comment, nil
}

// LikePost toggles a like on a post
func (s *ForumService) LikePost(ctx context.Context, userID, postID uint) (*model.LikeResponse, error) {
	// Check if like exists
	var likeID uint
	err := s.DB.QueryRow(ctx, "SELECT id FROM schema_forum.likes WHERE post_id = $1 AND user_id = $2", postID, userID).Scan(&likeID)

	resp := &model.LikeResponse{}
	if err == nil {
		// Like exists, toggle off
		_, err = s.DB.Exec(ctx, "DELETE FROM schema_forum.likes WHERE id = $1", likeID)
		if err != nil {
			return nil, err
		}
		_, _ = s.DB.Exec(ctx, "UPDATE schema_forum.posts SET like_count = like_count - 1 WHERE id = $1", postID)
		resp.Liked = false
	} else {
		// Like does not exist, toggle on
		_, err = s.DB.Exec(ctx, "INSERT INTO schema_forum.likes (post_id, user_id) VALUES ($1, $2)", postID, userID)
		if err != nil {
			return nil, err
		}
		_, _ = s.DB.Exec(ctx, "UPDATE schema_forum.posts SET like_count = like_count + 1 WHERE id = $1", postID)
		resp.Liked = true
	}

	// Get current like count
	_ = s.DB.QueryRow(ctx, "SELECT like_count FROM schema_forum.posts WHERE id = $1", postID).Scan(&resp.LikeCount)
	return resp, nil
}

// CollectPost toggles a collection on a post
func (s *ForumService) CollectPost(ctx context.Context, userID, postID uint) (*model.CollectResponse, error) {
	var collID uint
	err := s.DB.QueryRow(ctx, "SELECT id FROM schema_forum.collections WHERE post_id = $1 AND user_id = $2", postID, userID).Scan(&collID)

	resp := &model.CollectResponse{}
	if err == nil {
		// Collection exists, toggle off
		_, err = s.DB.Exec(ctx, "DELETE FROM schema_forum.collections WHERE id = $1", collID)
		if err != nil {
			return nil, err
		}
		resp.Collected = false
	} else {
		// Collection does not exist, toggle on
		_, err = s.DB.Exec(ctx, "INSERT INTO schema_forum.collections (post_id, user_id) VALUES ($1, $2)", postID, userID)
		if err != nil {
			return nil, err
		}
		resp.Collected = true
	}
	return resp, nil
}

// UploadAttachment saves an uploaded file or link
func (s *ForumService) UploadAttachment(ctx context.Context, userID uint, filename, fileType, filePath string, fileSize int64) (*model.Attachment, error) {
	attachment := &model.Attachment{}
	err := s.DB.QueryRow(ctx, `
		INSERT INTO schema_forum.attachments (user_id, filename, file_type, file_path, file_size)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, COALESCE(post_id, 0), COALESCE(comment_id, 0), user_id, filename, file_type, file_path, file_size, created_at
	`, userID, filename, fileType, filePath, fileSize).Scan(
		&attachment.ID, &attachment.PostID, &attachment.CommentID, &attachment.UserID,
		&attachment.Filename, &attachment.FileType, &attachment.FilePath, &attachment.FileSize,
		&attachment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save attachment: %w", err)
	}
	return attachment, nil
}

// GetAttachment retrieves an attachment by ID
func (s *ForumService) GetAttachment(ctx context.Context, id uint) (*model.Attachment, error) {
	a := &model.Attachment{}
	err := s.DB.QueryRow(ctx,
		"SELECT id, COALESCE(post_id, 0), COALESCE(comment_id, 0), user_id, filename, file_type, file_path, file_size, created_at FROM schema_forum.attachments WHERE id = $1",
		id,
	).Scan(&a.ID, &a.PostID, &a.CommentID, &a.UserID, &a.Filename, &a.FileType, &a.FilePath, &a.FileSize, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("attachment not found")
	}
	return a, nil
}

// GetPostAttachments returns attachments for a post
func (s *ForumService) GetPostAttachments(ctx context.Context, postID uint) ([]*model.Attachment, error) {
	rows, err := s.DB.Query(ctx,
		"SELECT id, COALESCE(post_id, 0), COALESCE(comment_id, 0), user_id, filename, file_type, file_path, file_size, created_at FROM schema_forum.attachments WHERE post_id = $1",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []*model.Attachment
	for rows.Next() {
		a := &model.Attachment{}
		var pid, cid uint
		if err := rows.Scan(&a.ID, &pid, &cid, &a.UserID, &a.Filename, &a.FileType, &a.FilePath, &a.FileSize, &a.CreatedAt); err == nil {
			attachments = append(attachments, a)
		}
	}
	return attachments, nil
}

// GetUserLevel fetches user level from JWT context (set by auth middleware)
// Level checking is done via middleware on route level, not in service layer.
func (s *ForumService) GetUserLevel(_ context.Context, _ uint) (int, error) {
	return 0, fmt.Errorf("level check should be done via middleware, not service layer")
}
