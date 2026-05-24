package model

import "time"

// Role represents an admin role with permissions
type Role struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Permissions []string  `json:"permissions" db:"permissions"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UserRole represents a user-to-role assignment
type UserRole struct {
	UserID     int64     `json:"user_id" db:"user_id"`
	RoleID     uint      `json:"role_id" db:"role_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
	AssignedBy int64     `json:"assigned_by" db:"assigned_by"`
}
