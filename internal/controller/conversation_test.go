// controller/conversation_controller_test.go
package controller

import (
	"bytes"
	"chat-apps/internal/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConversationService
type MockConversationService struct {
	mock.Mock
}

func (m *MockConversationService) CreateConversation(participants []int) (domain.ConversationResponse, error) {
	args := m.Called(participants)
	return args.Get(0).(domain.ConversationResponse), args.Error(1)
}

func (m *MockConversationService) GetConversationByID(id int) (domain.ConversationResponse, error) {
	args := m.Called(id)
	return args.Get(0).(domain.ConversationResponse), args.Error(1)
}

func TestCreateConversation(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)
	router := gin.Default()

	router.POST("/conversations", controller.CreateConversation)

	participants := []int{1, 2, 3}
	expectedResp := domain.ConversationResponse{
		ID:           1,
		Participants: participants,
		CreatedAt:    time.Now().Truncate(time.Second), // Truncate to match precision
	}

	mockService.On("CreateConversation", participants).Return(expectedResp, nil)

	payload := map[string]interface{}{
		"participants": participants,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/conversations", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ConversationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Compare ID and Participants directly, and check CreatedAt separately
	assert.Equal(t, expectedResp.ID, response.ID)
	assert.ElementsMatch(t, expectedResp.Participants, response.Participants)
	assert.WithinDuration(t, expectedResp.CreatedAt, response.CreatedAt, time.Second) // Allow 1-second tolerance
	mockService.AssertExpectations(t)
}

func TestGetConversationByID(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)
	router := gin.Default()

	router.GET("/conversations/:conversationId", controller.GetConversationByID)

	conversationID := 1
	expectedResp := domain.ConversationResponse{
		ID:           conversationID,
		Participants: []int{1, 2},
		CreatedAt:    time.Now().Truncate(time.Second), // Truncate to second precision
	}

	mockService.On("GetConversationByID", conversationID).Return(expectedResp, nil)

	req, _ := http.NewRequest("GET", "/conversations/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ConversationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Truncate response.CreatedAt to second precision for comparison
	response.CreatedAt = response.CreatedAt.Truncate(time.Second)

	// Assert equality
	assert.Equal(t, expectedResp, response)
	mockService.AssertExpectations(t)
}
