package repository

import (
	"chat-apps/internal/domain"

	"gorm.io/gorm"
)

type ConversationRepository interface {
	CreateConversation(participants []int) (domain.ConversationResponse, error)
	GetConversationByID(id int) (domain.ConversationResponse, error)
}

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) CreateConversation(participants []int) (domain.ConversationResponse, error) {
	var conversation domain.Conversation

	tx := r.db.Begin()

	if err := tx.Create(&conversation).Error; err != nil {
		tx.Rollback()
		return domain.ConversationResponse{}, err
	}

	for _, userID := range participants {
		cp := domain.ConversationParticipant{
			ConversationID: conversation.ID,
			UserID:         userID,
		}
		if err := tx.Create(&cp).Error; err != nil {
			tx.Rollback()
			return domain.ConversationResponse{}, err
		}
	}

	tx.Commit()

	var participantIDs []int
	if err := r.db.Model(&domain.ConversationParticipant{}).
		Where("conversation_id = ?", conversation.ID).
		Pluck("user_id", &participantIDs).Error; err != nil {
		return domain.ConversationResponse{}, err
	}

	var conversationResp domain.ConversationResponse
	conversationResp.ID = conversation.ID
	conversationResp.Participants = participantIDs
	conversationResp.CreatedAt = conversation.CreatedAt

	return conversationResp, nil
}

func (r *conversationRepository) GetConversationByID(id int) (domain.ConversationResponse, error) {
	var conversation domain.Conversation
	var participants []domain.ConversationParticipant

	if err := r.db.Where("id = ?", id).Find(&conversation).Error; err != nil {
		return domain.ConversationResponse{}, err
	}

	if err := r.db.Where("conversation_id = ?", id).Find(&participants).Error; err != nil {
		return domain.ConversationResponse{}, err
	}

	var conversationResp domain.ConversationResponse
	for _, participant := range participants {
		conversationResp.Participants = append(conversationResp.Participants, participant.UserID)
	}
	conversationResp.ID = conversation.ID
	conversationResp.CreatedAt = conversation.CreatedAt

	return conversationResp, nil
}
