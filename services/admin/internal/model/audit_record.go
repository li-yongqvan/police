package model

import "time"

// AuditRecord represents a post pending audit (DTO for audit page)
type AuditRecord struct {
	PostID      uint      `json:"post_id"`
	PostTitle   string    `json:"post_title"`
	AuthorID    uint      `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	BoardID     uint      `json:"board_id"`
	BoardName   string    `json:"board_name"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	Reason      string    `json:"reason,omitempty"` // matched sensitive words
}
