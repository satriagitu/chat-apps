package domain

import "time"

type Notification struct {
	ID      int       `json:"id" gorm:"primaryKey"`
	UserID  int       `json:"user_id"`
	Message string    `json:"message"`
	SentAt  time.Time `json:"sent_at"`
}

type NotificationRequest struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}
