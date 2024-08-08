package domain

import "time"

type Job struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	Message     string     `json:"message"`
	Status      string     `json:"status"`
	QueuedAt    time.Time  `json:"queued_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type BroadcastRequest struct {
	Message string `json:"message"`
}
