// service/message_service_test.go
package service

import (
	"chat-apps/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) CreateMessage(message domain.Message) (domain.Message, error) {
	args := m.Called(message)
	return args.Get(0).(domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetMessagesByConversationID(conversationID int) ([]domain.Message, error) {
	args := m.Called(conversationID)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func TestCreateMessage(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	conversationID := 1
	senderID := 2
	content := "Hello"
	expectedMessage := domain.Message{
		ID:             1,
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		SentAt:         time.Now(),
	}

	mockRepo.On("CreateMessage", domain.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
	}).Return(expectedMessage, nil)

	message, err := service.CreateMessage(conversationID, senderID, content)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, message)
	mockRepo.AssertExpectations(t)
}

func TestGetMessagesByConversationID(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	conversationID := 1
	expectedMessages := []domain.Message{
		{ID: 1, ConversationID: conversationID, SenderID: 2, Content: "Hello", SentAt: time.Now()},
		{ID: 2, ConversationID: conversationID, SenderID: 3, Content: "Hi", SentAt: time.Now()},
	}

	mockRepo.On("GetMessagesByConversationID", conversationID).Return(expectedMessages, nil)

	messages, err := service.GetMessagesByConversationID(conversationID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}
