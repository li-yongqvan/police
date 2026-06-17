package model

import (
	"time"
)

// User represents a user account in the system
type User struct {
	ID           uint      `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Nickname     string    `json:"nickname" db:"nickname"`
	Bio          string    `json:"bio" db:"bio"`
	Avatar       string    `json:"avatar" db:"avatar"`
	Level            int       `json:"level" db:"level"`
	Status           string    `json:"status" db:"status"`
	Department       string    `json:"department" db:"department"`
	Squad            string    `json:"squad" db:"squad"`
	Grade            string    `json:"grade" db:"grade"`
	ProfileCompleted bool      `json:"profile_completed" db:"profile_completed"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UserResponse is the sanitized user data returned to clients (excludes password_hash)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	Level            int       `json:"level"`
	Status           string    `json:"status"`
	Department       string    `json:"department"`
	Squad            string    `json:"squad"`
	Grade            string    `json:"grade"`
	ProfileCompleted bool      `json:"profile_completed"`
	CreatedAt        time.Time `json:"created_at"`
}

// ToResponse converts a User to a UserResponse (excludes password_hash)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Bio:       u.Bio,
		Avatar:    u.Avatar,
		Level:            u.Level,
		Status:           u.Status,
		Department:       u.Department,
		Squad:            u.Squad,
		Grade:            u.Grade,
		ProfileCompleted: u.ProfileCompleted,
		CreatedAt:        u.CreatedAt,
	}
}

// RegisterRequest represents the registration input
type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	InvitationCode string `json:"invitation_code" binding:"required"`
	Department     string `json:"department"`
	Squad          string `json:"squad"`
	Grade          string `json:"grade"`
}

// LoginRequest represents the login input
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login output with JWT tokens
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateProfileRequest represents profile update input
type UpdateProfileRequest struct {
	Username         string `json:"username"`
	Nickname         string `json:"nickname"`
	Avatar           string `json:"avatar"`
	Bio              string `json:"bio"`
	Department       string `json:"department"`
	Squad            string `json:"squad"`
	Grade            string `json:"grade"`
	ProfileCompleted *bool  `json:"profile_completed"`
}
