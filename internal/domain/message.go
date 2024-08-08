package domain

import "time"

type Message struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	Content        string    `json:"content"`
	SentAt         time.Time `json:"sent_at"`
}
