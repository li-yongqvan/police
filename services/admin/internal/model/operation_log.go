package model

import "time"

// OperationLog represents an admin action log entry
type OperationLog struct {
	ID               uint                 `json:"id" db:"id"`
	OperatorID       int64                `json:"operator_id" db:"operator_id"`
	OperatorUsername string               `json:"operator_username" db:"operator_username"`
	Action           string               `json:"action" db:"action"`
	TargetType       string               `json:"target_type" db:"target_type"`
	TargetID         int64                `json:"target_id" db:"target_id"`
	Detail           map[string]interface{} `json:"detail" db:"detail"`
	CreatedAt        time.Time            `json:"created_at" db:"created_at"`
}
