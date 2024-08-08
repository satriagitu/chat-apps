package domain

import "time"

type Conversation struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversationParticipant struct {
	ConversationID int `json:"conversation_id"`
	UserID         int `json:"user_id"`
}

type ConversationResponse struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Participants []int     `json:"participants"`
	CreatedAt    time.Time `json:"created_at"`
}
