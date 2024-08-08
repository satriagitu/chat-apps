package service

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"
)

type MessageService interface {
	CreateMessage(conversationID int, senderID int, content string) (domain.Message, error)
	GetMessagesByConversationID(conversationID int) ([]domain.Message, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService(mr repository.MessageRepository) MessageService {
	return &messageService{messageRepo: mr}
}

func (s *messageService) CreateMessage(conversationID int, senderID int, content string) (domain.Message, error) {
	message := domain.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
	}
	return s.messageRepo.CreateMessage(message)
}

func (s *messageService) GetMessagesByConversationID(conversationID int) ([]domain.Message, error) {
	return s.messageRepo.GetMessagesByConversationID(conversationID)
}
