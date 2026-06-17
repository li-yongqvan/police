package model

import "time"

// SystemConfig represents a system configuration entry
type SystemConfig struct {
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy   int64     `json:"updated_by" db:"updated_by"`
}
