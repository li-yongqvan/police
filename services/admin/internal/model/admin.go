package model

import (
	"time"
)

// SensitiveWord represents a word in the moderation filter
type SensitiveWord struct {
	ID        uint      `json:"id" db:"id"`
	Word      string    `json:"word" db:"word"`
	Category  string    `json:"category" db:"category"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ViolationRecord represents a user violation record
type ViolationRecord struct {
	ID         uint      `json:"id" db:"id"`
	UserID     uint      `json:"user_id" db:"user_id"`
	PostID     uint      `json:"post_id" db:"post_id"`
	Violation  string    `json:"violation" db:"violation"`
	Action     string    `json:"action" db:"action"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// DailyStats represents aggregated daily statistics
type DailyStats struct {
	ID            uint      `json:"id" db:"id"`
	Date          string    `json:"date" db:"date"`
	NewUsers      int       `json:"new_users" db:"new_users"`
	NewPosts      int       `json:"new_posts" db:"new_posts"`
	NewComments   int       `json:"new_comments" db:"new_comments"`
	ActiveUsers   int       `json:"active_users" db:"active_users"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
