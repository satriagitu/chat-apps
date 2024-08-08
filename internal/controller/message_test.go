// controller/message_controller_test.go
package controller

import (
	"chat-apps/internal/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMessageService
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) CreateMessage(conversationID int, senderID int, content string) (domain.Message, error) {
	args := m.Called(conversationID, senderID, content)
	return args.Get(0).(domain.Message), args.Error(1)
}

func (m *MockMessageService) GetMessagesByConversationID(conversationID int) ([]domain.Message, error) {
	args := m.Called(conversationID)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func TestCreateMessage(t *testing.T) {
	mockService := new(MockMessageService)
	controller := NewMessageController(mockService)
	router := gin.Default()

	router.POST("/conversations/:conversationId/messages", controller.CreateMessage)

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

	mockService.On("CreateMessage", conversationID, senderID, content).Return(expectedMessage, nil)

	reqBody := `{"sender_id":2, "content":"Hello"}`
	req, _ := http.NewRequest("POST", "/conversations/1/messages", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage.ID, response.ID)
	assert.Equal(t, expectedMessage.ConversationID, response.ConversationID)
	assert.Equal(t, expectedMessage.SenderID, response.SenderID)
	assert.Equal(t, expectedMessage.Content, response.Content)
	mockService.AssertExpectations(t)
}

func TestGetMessagesByConversationID(t *testing.T) {
	mockService := new(MockMessageService)
	controller := NewMessageController(mockService)
	router := gin.Default()

	router.GET("/conversations/:conversationId/messages", controller.GetMessagesByConversationID)

	conversationID := 1
	expectedMessages := []domain.Message{
		{ID: 1, ConversationID: conversationID, SenderID: 2, Content: "Hello", SentAt: time.Now().UTC()},
		{ID: 2, ConversationID: conversationID, SenderID: 3, Content: "Hi", SentAt: time.Now().UTC()},
	}

	mockService.On("GetMessagesByConversationID", conversationID).Return(expectedMessages, nil)

	req, _ := http.NewRequest("GET", "/conversations/1/messages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	// Check that the messages are equal
	for i, msg := range expectedMessages {
		assert.Equal(t, msg.ID, response[i].ID)
		assert.Equal(t, msg.ConversationID, response[i].ConversationID)
		assert.Equal(t, msg.SenderID, response[i].SenderID)
		assert.Equal(t, msg.Content, response[i].Content)
		assert.WithinDuration(t, msg.SentAt, response[i].SentAt, time.Second)
	}
	mockService.AssertExpectations(t)
}
