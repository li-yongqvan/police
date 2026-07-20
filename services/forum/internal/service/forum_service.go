package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"ai-forum/forum-service/internal/client"
	"ai-forum/forum-service/internal/model"

	"github.com/jackc/pgx/v5"
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

// ListPosts returns paginated posts, optionally filtered by board, keyword, and sort mode.
func (s *ForumService) ListPosts(ctx context.Context, boardID, authorID uint, page, limit int, keyword, sort string, viewerID uint) ([]*model.PostListItem, int, error) {
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
	argN := 1
	if boardID > 0 {
		countQuery += fmt.Sprintf(" AND board_id = $%d", argN)
		args = append(args, boardID)
		argN++
	}
	if authorID > 0 {
		countQuery += fmt.Sprintf(" AND author_id = $%d", argN)
		args = append(args, authorID)
		argN++
	}
	if keyword != "" {
		countQuery += fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d)", argN, argN)
		args = append(args, "%"+keyword+"%")
	}
	countQuery += postSortCountFilter(sort)
	var total int
	err := s.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Fetch posts
	query := `
		SELECT p.id, p.title, LEFT(p.content, 320), p.author_id, u.username, u.avatar, p.board_id, b.name, b.slug,
		       p.status, p.is_pinned, p.is_featured, p.like_count, p.dislike_count, p.comment_count, p.created_at
		FROM schema_forum.posts p
		JOIN schema_auth.users u ON p.author_id = u.id
		JOIN schema_forum.boards b ON p.board_id = b.id
		WHERE p.status = 'published'`
	qArgs := []interface{}{}
	qN := 1
	if boardID > 0 {
		query += fmt.Sprintf(" AND p.board_id = $%d", qN)
		qArgs = append(qArgs, boardID)
		qN++
	}
	if authorID > 0 {
		query += fmt.Sprintf(" AND p.author_id = $%d", qN)
		qArgs = append(qArgs, authorID)
		qN++
	}
	if keyword != "" {
		query += fmt.Sprintf(" AND (p.title ILIKE $%d OR p.content ILIKE $%d)", qN, qN)
		qArgs = append(qArgs, "%"+keyword+"%")
	}
	query += postSortListFilter(sort)
	query += postSortOrderBy(sort)
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := s.DB.Query(ctx, query, qArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*model.PostListItem
	for rows.Next() {
		p := &model.PostListItem{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.AuthorAvatar,
			&p.BoardID, &p.BoardName, &p.BoardSlug, &p.Status, &p.IsPinned, &p.IsFeatured,
			&p.LikeCount, &p.DislikeCount, &p.CommentCount, &p.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, p)
	}
	if viewerID > 0 && len(posts) > 0 {
		applyPostInteractions(ctx, s.DB, viewerID, posts)
	}
	return posts, total, nil
}

// ListUserCollections returns posts the user has bookmarked.
func (s *ForumService) ListUserCollections(ctx context.Context, userID uint, page, limit int) ([]*model.PostListItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int
	err := s.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM schema_forum.collections c
		JOIN schema_forum.posts p ON p.id = c.post_id
		WHERE c.user_id = $1 AND p.status = 'published'
	`, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count collections: %w", err)
	}

	rows, err := s.DB.Query(ctx, `
		SELECT p.id, p.title, LEFT(p.content, 320), p.author_id, u.username, u.avatar, p.board_id, b.name, b.slug,
		       p.status, p.is_pinned, p.is_featured, p.like_count, p.dislike_count, p.comment_count, p.created_at
		FROM schema_forum.collections c
		JOIN schema_forum.posts p ON p.id = c.post_id
		JOIN schema_auth.users u ON p.author_id = u.id
		JOIN schema_forum.boards b ON p.board_id = b.id
		WHERE c.user_id = $1 AND p.status = 'published'
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query collections: %w", err)
	}
	defer rows.Close()

	var posts []*model.PostListItem
	for rows.Next() {
		p := &model.PostListItem{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.AuthorAvatar,
			&p.BoardID, &p.BoardName, &p.BoardSlug, &p.Status, &p.IsPinned, &p.IsFeatured,
			&p.LikeCount, &p.DislikeCount, &p.CommentCount, &p.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan collection post: %w", err)
		}
		posts = append(posts, p)
	}
	if len(posts) > 0 {
		applyPostInteractions(ctx, s.DB, userID, posts)
	}
	return posts, total, nil
}

func applyPostInteractions(ctx context.Context, db *pgxpool.Pool, viewerID uint, posts []*model.PostListItem) {
	ids := make([]int64, len(posts))
	for i, p := range posts {
		ids[i] = int64(p.ID)
	}
	rows, err := db.Query(ctx, `
		SELECT post_id FROM schema_forum.likes
		WHERE user_id = $1 AND post_id = ANY($2)
	`, viewerID, ids)
	if err != nil {
		return
	}
	defer rows.Close()
	liked := make(map[uint]bool)
	for rows.Next() {
		var pid uint
		if rows.Scan(&pid) == nil {
			liked[pid] = true
		}
	}
	for _, p := range posts {
		p.Liked = liked[p.ID]
	}
}

func postSortCountFilter(sort string) string {
	switch sort {
	case "featured":
		return " AND (is_pinned = true OR is_featured = true)"
	case "today":
		return " AND created_at >= CURRENT_DATE"
	default:
		return ""
	}
}

func postSortListFilter(sort string) string {
	switch sort {
	case "featured":
		return " AND (p.is_pinned = true OR p.is_featured = true)"
	case "today":
		return " AND p.created_at >= CURRENT_DATE"
	default:
		return ""
	}
}

func postSortOrderBy(sort string) string {
	switch sort {
	case "new":
		return " ORDER BY p.is_pinned DESC, p.created_at DESC"
	case "featured":
		return " ORDER BY p.is_pinned DESC, p.created_at DESC"
	case "today", "hot":
		return " ORDER BY p.is_pinned DESC, (p.like_count * 2 + p.comment_count * 3) DESC, p.created_at DESC"
	default:
		return " ORDER BY p.is_pinned DESC, (p.like_count * 2 + p.comment_count * 3) DESC, p.created_at DESC"
	}
}

type PostNotifyMeta struct {
	AuthorID uint
	Title    string
}

func (s *ForumService) GetPostNotifyMeta(ctx context.Context, id uint) (*PostNotifyMeta, error) {
	meta := &PostNotifyMeta{}
	err := s.DB.QueryRow(ctx, `SELECT author_id, title FROM schema_forum.posts WHERE id = $1`, id).
		Scan(&meta.AuthorID, &meta.Title)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

// GetPost returns full post detail with attachments
func (s *ForumService) GetPost(ctx context.Context, id uint, viewerID uint) (*model.PostDetail, error) {
	detail := &model.PostDetail{}
	err := s.DB.QueryRow(ctx, `
		SELECT p.id, p.title, p.content, p.author_id, u.username, u.avatar, p.board_id, b.name,
		       p.status, p.is_pinned, p.is_featured, p.like_count, p.dislike_count, p.comment_count,
		       p.created_at, p.updated_at
		FROM schema_forum.posts p
		JOIN schema_auth.users u ON p.author_id = u.id
		JOIN schema_forum.boards b ON p.board_id = b.id
		WHERE p.id = $1 AND p.status IN ('published', 'pending_review')
	`, id).Scan(
		&detail.ID, &detail.Title, &detail.Content, &detail.AuthorID, &detail.AuthorName, &detail.AuthorAvatar,
		&detail.BoardID, &detail.BoardName, &detail.Status, &detail.IsPinned, &detail.IsFeatured,
		&detail.LikeCount, &detail.DislikeCount, &detail.CommentCount, &detail.CreatedAt, &detail.UpdatedAt,
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

	if viewerID > 0 {
		_ = s.DB.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM schema_forum.likes WHERE post_id = $1 AND user_id = $2)",
			id, viewerID,
		).Scan(&detail.Liked)
		_ = s.DB.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM schema_forum.dislikes WHERE post_id = $1 AND user_id = $2)",
			id, viewerID,
		).Scan(&detail.Disliked)
		_ = s.DB.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM schema_forum.collections WHERE post_id = $1 AND user_id = $2)",
			id, viewerID,
		).Scan(&detail.Collected)
	}

	return detail, nil
}

// CreatePost creates a new post with sensitive word check
func (s *ForumService) CreatePost(ctx context.Context, authorID uint, req *model.CreatePostRequest) (*model.Post, error) {
	// Check sensitive words
	clean, matchedWords, err := s.AdminClient.CheckSensitiveWords(req.Title + " " + req.Content)
	if err != nil {
		// If moderation is unavailable, send the post to manual review.
		clean = false
	}

	status := "published"
	if !clean {
		status = "pending_review"
	}

	post := &model.Post{}
	err = s.DB.QueryRow(ctx, `
		INSERT INTO schema_forum.posts (title, content, author_id, board_id, status, matched_words)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, content, author_id, board_id, status, is_pinned, is_featured, like_count, comment_count, matched_words, created_at, updated_at
	`, req.Title, req.Content, authorID, req.BoardID, status, matchedWords).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.BoardID,
		&post.Status, &post.IsPinned, &post.IsFeatured, &post.LikeCount,
		&post.CommentCount, &post.MatchedWords, &post.CreatedAt, &post.UpdatedAt,
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
			log.Printf("forum: failed to link attachments for post %d: %v", post.ID, err)
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
		return nil, fmt.Errorf("閺冪姵娼堟穱顔芥暭娴犳牔姹夌敮鏍х摍")
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
	clean, matchedWords, _ := s.AdminClient.CheckSensitiveWords(title + " " + content)
	newStatus := existing.Status
	if !clean {
		newStatus = "pending_review"
	}

	err = s.DB.QueryRow(ctx, `
		UPDATE schema_forum.posts
		SET title = COALESCE(NULLIF($1, ''), title),
		    content = COALESCE(NULLIF($2, ''), content),
		    status = $3,
		    matched_words = $4,
		    updated_at = NOW()
		WHERE id = $5
		RETURNING id, title, content, author_id, board_id, status, is_pinned, is_featured, like_count, comment_count, matched_words, created_at, updated_at
	`, title, content, newStatus, matchedWords, postID).Scan(
		&existing.ID, &existing.Title, &existing.Content, &existing.AuthorID,
		&existing.BoardID, &existing.Status, &existing.IsPinned, &existing.IsFeatured,
		&existing.LikeCount, &existing.CommentCount, &existing.MatchedWords, &existing.CreatedAt, &existing.UpdatedAt,
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
		return fmt.Errorf("閺冪姵娼堥崚鐘绘珟娴犳牔姹夌敮鏍х摍")
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
func (s *ForumService) ListComments(ctx context.Context, postID uint, page, limit int, viewerID uint) ([]*model.Comment, int, error) {
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

	viewer := int64(viewerID)
	rows, err := s.DB.Query(ctx, `
		SELECT c.id, c.post_id, c.parent_id, c.depth, c.author_id, u.username, u.avatar, c.content,
		       c.like_count, c.dislike_count,
		       EXISTS(SELECT 1 FROM schema_forum.comment_likes    WHERE comment_id = c.id AND user_id = $4) AS liked,
		       EXISTS(SELECT 1 FROM schema_forum.comment_dislikes WHERE comment_id = c.id AND user_id = $4) AS disliked,
		       c.created_at
		FROM schema_forum.comments c
		JOIN schema_auth.users u ON c.author_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3`,
		postID, limit, offset, viewer,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Depth, &c.AuthorID, &c.AuthorName, &c.AuthorAvatar, &c.Content,
			&c.LikeCount, &c.DislikeCount, &c.Liked, &c.Disliked, &c.CreatedAt); err != nil {
			return nil, 0, err
		}
		comments = append(comments, c)
	}
	return comments, total, nil
}

// CreateComment creates a new comment with sensitive word check and optional parent reply.
func (s *ForumService) CreateComment(ctx context.Context, authorID, postID uint, content string, parentID *uint) (*model.Comment, error) {
	clean, _, err := s.AdminClient.CheckSensitiveWords(content)
	if err != nil {
		return nil, fmt.Errorf("moderation service unavailable")
	}
	if !clean {
		return nil, fmt.Errorf("comment contains sensitive words, rejected")
	}

	depth := 0
	if parentID != nil && *parentID > 0 {
		var parentPostID uint
		var parentDepth int
		err = s.DB.QueryRow(ctx,
			`SELECT post_id, depth FROM schema_forum.comments WHERE id = $1`,
			*parentID,
		).Scan(&parentPostID, &parentDepth)
		if err != nil || parentPostID != postID {
			return nil, fmt.Errorf("invalid parent comment")
		}
		if parentDepth >= 8 {
			return nil, fmt.Errorf("reply depth limit reached")
		}
		depth = parentDepth + 1
	}

	comment := &model.Comment{ParentID: parentID, Depth: depth}
	err = s.DB.QueryRow(ctx, `
		INSERT INTO schema_forum.comments (post_id, parent_id, depth, author_id, content)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, post_id, parent_id, depth, author_id, content, created_at
	`, postID, parentID, depth, authorID, content).Scan(
		&comment.ID, &comment.PostID, &comment.ParentID, &comment.Depth, &comment.AuthorID, &comment.Content, &comment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	_ = s.DB.QueryRow(ctx, `SELECT username, avatar FROM schema_auth.users WHERE id = $1`, authorID).Scan(&comment.AuthorName, &comment.AuthorAvatar)

	_, _ = s.DB.Exec(ctx, "UPDATE schema_forum.posts SET comment_count = comment_count + 1 WHERE id = $1", postID)

	return comment, nil
}

// LikePost toggles a like on a post
func (s *ForumService) LikePost(ctx context.Context, userID, postID uint) (*model.LikeResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var likeID uint
	err = tx.QueryRow(ctx, "SELECT id FROM schema_forum.likes WHERE post_id = $1 AND user_id = $2", postID, userID).Scan(&likeID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to query like state: %w", err)
	}

	resp := &model.LikeResponse{}
	if err == nil {
		_, err = tx.Exec(ctx, "DELETE FROM schema_forum.likes WHERE id = $1", likeID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete like: %w", err)
		}
		if _, err = tx.Exec(ctx, "UPDATE schema_forum.posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1", postID); err != nil {
			return nil, fmt.Errorf("failed to decrement like count: %w", err)
		}
		resp.Liked = false
	} else {
		_, err = tx.Exec(ctx, "INSERT INTO schema_forum.likes (post_id, user_id) VALUES ($1, $2)", postID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert like: %w", err)
		}
		if _, err = tx.Exec(ctx, "UPDATE schema_forum.posts SET like_count = like_count + 1 WHERE id = $1", postID); err != nil {
			return nil, fmt.Errorf("failed to increment like count: %w", err)
		}
		resp.Liked = true
		tag, err := tx.Exec(ctx, "DELETE FROM schema_forum.dislikes WHERE post_id = $1 AND user_id = $2", postID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to clear dislike: %w", err)
		}
		if tag.RowsAffected() > 0 {
			if _, err = tx.Exec(ctx, "UPDATE schema_forum.posts SET dislike_count = GREATEST(dislike_count - 1, 0) WHERE id = $1", postID); err != nil {
				return nil, fmt.Errorf("failed to decrement dislike count: %w", err)
			}
		}
	}

	if err := tx.QueryRow(ctx, `
		SELECT like_count, dislike_count,
		       EXISTS(SELECT 1 FROM schema_forum.dislikes WHERE post_id = $1 AND user_id = $2)
		FROM schema_forum.posts WHERE id = $1
	`, postID, userID).Scan(&resp.LikeCount, &resp.DislikeCount, &resp.Disliked); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return resp, nil
}

// DislikePost toggles a dislike on a post
func (s *ForumService) DislikePost(ctx context.Context, userID, postID uint) (*model.DislikeResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var dislikeID uint
	err = tx.QueryRow(ctx, "SELECT id FROM schema_forum.dislikes WHERE post_id = $1 AND user_id = $2", postID, userID).Scan(&dislikeID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to query dislike state: %w", err)
	}

	resp := &model.DislikeResponse{}
	if err == nil {
		_, err = tx.Exec(ctx, "DELETE FROM schema_forum.dislikes WHERE id = $1", dislikeID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete dislike: %w", err)
		}
		if _, err = tx.Exec(ctx, "UPDATE schema_forum.posts SET dislike_count = GREATEST(dislike_count - 1, 0) WHERE id = $1", postID); err != nil {
			return nil, fmt.Errorf("failed to decrement dislike count: %w", err)
		}
		resp.Disliked = false
	} else {
		_, err = tx.Exec(ctx, "INSERT INTO schema_forum.dislikes (post_id, user_id) VALUES ($1, $2)", postID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert dislike: %w", err)
		}
		if _, err = tx.Exec(ctx, "UPDATE schema_forum.posts SET dislike_count = dislike_count + 1 WHERE id = $1", postID); err != nil {
			return nil, fmt.Errorf("failed to increment dislike count: %w", err)
		}
		resp.Disliked = true
		tag, err := tx.Exec(ctx, "DELETE FROM schema_forum.likes WHERE post_id = $1 AND user_id = $2", postID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to clear like: %w", err)
		}
		if tag.RowsAffected() > 0 {
			if _, err = tx.Exec(ctx, "UPDATE schema_forum.posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1", postID); err != nil {
				return nil, fmt.Errorf("failed to decrement like count: %w", err)
			}
		}
	}

	if err := tx.QueryRow(ctx, `
		SELECT like_count, dislike_count,
		       EXISTS(SELECT 1 FROM schema_forum.likes WHERE post_id = $1 AND user_id = $2)
		FROM schema_forum.posts WHERE id = $1
	`, postID, userID).Scan(&resp.LikeCount, &resp.DislikeCount, &resp.Liked); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
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

// LikeComment toggles a like on a comment
func (s *ForumService) LikeComment(ctx context.Context, userID, commentID uint) (*model.CommentLikeResponse, error) {
	var likeID uint
	err := s.DB.QueryRow(ctx, "SELECT id FROM schema_forum.comment_likes WHERE comment_id = $1 AND user_id = $2", commentID, userID).Scan(&likeID)

	resp := &model.CommentLikeResponse{}
	if err == nil {
		_, err = s.DB.Exec(ctx, "DELETE FROM schema_forum.comment_likes WHERE id = $1", likeID)
		if err != nil {
			return nil, err
		}
		_, _ = s.DB.Exec(ctx, "UPDATE schema_forum.comments SET like_count = like_count - 1 WHERE id = $1", commentID)
		resp.Liked = false
	} else {
		_, err = s.DB.Exec(ctx, "INSERT INTO schema_forum.comment_likes (comment_id, user_id) VALUES ($1, $2)", commentID, userID)
		if err != nil {
			return nil, err
		}
		_, _ = s.DB.Exec(ctx, "UPDATE schema_forum.comments SET like_count = like_count + 1 WHERE id = $1", commentID)
		resp.Liked = true
	}

	_ = s.DB.QueryRow(ctx, "SELECT like_count FROM schema_forum.comments WHERE id = $1", commentID).Scan(&resp.LikeCount)
	return resp, nil
}

// DislikeComment toggles a dislike on a comment
func (s *ForumService) DislikeComment(ctx context.Context, userID, commentID uint) (*model.CommentDislikeResponse, error) {
	var dislikeID uint
	err := s.DB.QueryRow(ctx, "SELECT id FROM schema_forum.comment_dislikes WHERE comment_id = $1 AND user_id = $2", commentID, userID).Scan(&dislikeID)

	resp := &model.CommentDislikeResponse{}
	if err == nil {
		_, err = s.DB.Exec(ctx, "DELETE FROM schema_forum.comment_dislikes WHERE id = $1", dislikeID)
		if err != nil {
			return nil, err
		}
		_, _ = s.DB.Exec(ctx, "UPDATE schema_forum.comments SET dislike_count = dislike_count - 1 WHERE id = $1", commentID)
		resp.Disliked = false
	} else {
		_, err = s.DB.Exec(ctx, "INSERT INTO schema_forum.comment_dislikes (comment_id, user_id) VALUES ($1, $2)", commentID, userID)
		if err != nil {
			return nil, err
		}
		_, _ = s.DB.Exec(ctx, "UPDATE schema_forum.comments SET dislike_count = dislike_count + 1 WHERE id = $1", commentID)
		resp.Disliked = true
	}

	_ = s.DB.QueryRow(ctx, "SELECT dislike_count FROM schema_forum.comments WHERE id = $1", commentID).Scan(&resp.DislikeCount)
	return resp, nil
}
