package service

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"
)

type ConversationService interface {
	CreateConversation(participants []int) (domain.ConversationResponse, error)
	GetConversationByID(id int) (domain.ConversationResponse, error)
}

type conversationService struct {
	conversationRepo repository.ConversationRepository
}

func NewConversationService(cr repository.ConversationRepository) ConversationService {
	return &conversationService{conversationRepo: cr}
}

func (s *conversationService) CreateConversation(participants []int) (domain.ConversationResponse, error) {
	return s.conversationRepo.CreateConversation(participants)
}

func (s *conversationService) GetConversationByID(id int) (domain.ConversationResponse, error) {
	return s.conversationRepo.GetConversationByID(id)
}
