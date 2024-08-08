// service/conversation_service_test.go
package service

import (
	"chat-apps/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConversationRepository
type MockConversationRepository struct {
	mock.Mock
}

func (m *MockConversationRepository) CreateConversation(participants []int) (domain.ConversationResponse, error) {
	args := m.Called(participants)
	return args.Get(0).(domain.ConversationResponse), args.Error(1)
}

func (m *MockConversationRepository) GetConversationByID(id int) (domain.ConversationResponse, error) {
	args := m.Called(id)
	return args.Get(0).(domain.ConversationResponse), args.Error(1)
}

func TestCreateConversation(t *testing.T) {
	mockRepo := new(MockConversationRepository)
	service := NewConversationService(mockRepo)

	participants := []int{1, 2, 3}
	expectedResp := domain.ConversationResponse{ID: 1, Participants: participants, CreatedAt: time.Now()}

	mockRepo.On("CreateConversation", participants).Return(expectedResp, nil)

	resp, err := service.CreateConversation(participants)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetConversationByID(t *testing.T) {
	mockRepo := new(MockConversationRepository)
	service := NewConversationService(mockRepo)

	conversationID := 1
	expectedResp := domain.ConversationResponse{ID: conversationID, Participants: []int{1, 2}, CreatedAt: time.Now()}

	mockRepo.On("GetConversationByID", conversationID).Return(expectedResp, nil)

	resp, err := service.GetConversationByID(conversationID)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockRepo.AssertExpectations(t)
}
