package repository

import (
	"chat-apps/internal/domain"
	"time"

	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(message domain.Message) (domain.Message, error)
	GetMessagesByConversationID(conversationID int) ([]domain.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(message domain.Message) (domain.Message, error) {
	message.SentAt = time.Now()
	if err := r.db.Create(&message).Error; err != nil {
		return domain.Message{}, err
	}
	return message, nil
}

func (r *messageRepository) GetMessagesByConversationID(conversationID int) ([]domain.Message, error) {
	var messages []domain.Message
	if err := r.db.Where("conversation_id = ?", conversationID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
