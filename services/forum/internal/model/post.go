package model

import (
	"time"
)

// Post represents a forum post
type Post struct {
	ID           uint      `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Content      string    `json:"content" db:"content"`
	AuthorID     uint      `json:"author_id" db:"author_id"`
	BoardID      uint      `json:"board_id" db:"board_id"`
	Status       string    `json:"status" db:"status"`
	IsPinned     bool      `json:"is_pinned" db:"is_pinned"`
	IsFeatured   bool      `json:"is_featured" db:"is_featured"`
	LikeCount    int       `json:"like_count" db:"like_count"`
	CommentCount int       `json:"comment_count" db:"comment_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Board represents a forum board/category
type Board struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	Enabled     bool      `json:"enabled" db:"enabled"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Comment represents a reply to a post
type Comment struct {
	ID         uint      `json:"id" db:"id"`
	PostID     uint      `json:"post_id" db:"post_id"`
	ParentID   *uint     `json:"parent_id,omitempty" db:"parent_id"`
	Depth      int       `json:"depth" db:"depth"`
	AuthorID   uint      `json:"author_id" db:"author_id"`
	AuthorName string    `json:"author_name,omitempty" db:"-"`
	Content    string    `json:"content" db:"content"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Like represents a post like
type Like struct {
	ID        uint      `json:"id" db:"id"`
	PostID    uint      `json:"post_id" db:"post_id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Collection represents a bookmarked post
type Collection struct {
	ID        uint      `json:"id" db:"id"`
	PostID    uint      `json:"post_id" db:"post_id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Attachment represents an uploaded file or link
type Attachment struct {
	ID        uint      `json:"id" db:"id"`
	PostID    *uint     `json:"post_id" db:"post_id"`
	CommentID *uint     `json:"comment_id" db:"comment_id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	Filename  string    `json:"filename" db:"filename"`
	FileType  string    `json:"file_type" db:"file_type"`
	FilePath  string    `json:"file_path" db:"file_path"`
	FileSize  int64     `json:"file_size" db:"file_size"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// BoardResponse is the response for board listing
type BoardResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	PostCount   int    `json:"post_count"`
}

// PostListItem is a condensed post for list views
type PostListItem struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content,omitempty"`
	AuthorID     uint      `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	BoardID      uint      `json:"board_id"`
	BoardName    string    `json:"board_name"`
	BoardSlug    string    `json:"board_slug,omitempty"`
	Status       string    `json:"status"`
	IsPinned     bool      `json:"is_pinned"`
	IsFeatured   bool      `json:"is_featured"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	Liked        bool      `json:"liked,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// PostDetail is the full post detail view
type PostDetail struct {
	ID           uint         `json:"id"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	AuthorID     uint         `json:"author_id"`
	AuthorName   string       `json:"author_name"`
	BoardID      uint         `json:"board_id"`
	BoardName    string       `json:"board_name"`
	Status       string       `json:"status"`
	IsPinned     bool         `json:"is_pinned"`
	IsFeatured   bool         `json:"is_featured"`
	LikeCount    int          `json:"like_count"`
	CommentCount int          `json:"comment_count"`
	Liked        bool         `json:"liked,omitempty"`
	Collected    bool         `json:"collected,omitempty"`
	Attachments  []Attachment `json:"attachments"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// CreatePostRequest represents the input for creating a post
type CreatePostRequest struct {
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content" binding:"required"`
	BoardID       uint   `json:"board_id" binding:"required"`
	AttachmentIDs []uint `json:"attachment_ids"`
}

// UpdatePostRequest represents the input for updating a post
type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateCommentRequest represents the input for creating a comment
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

// LikeResponse represents the response for like toggle
type LikeResponse struct {
	Liked     bool `json:"liked"`
	LikeCount int  `json:"like_count"`
}

// CollectResponse represents the response for collection toggle
type CollectResponse struct {
	Collected bool `json:"collected"`
}
