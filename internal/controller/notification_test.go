// controller/notification_controller_test.go
package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"chat-apps/internal/domain"
)

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendNotification(userID int, message string) (domain.Notification, error) {
	args := m.Called(userID, message)
	return args.Get(0).(domain.Notification), args.Error(1)
}

func (m *MockNotificationService) GetNotificationsByUserID(userID int) ([]domain.Notification, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Notification), args.Error(1)
}

func TestSendNotification(t *testing.T) {
	mockService := new(MockNotificationService)
	controller := NewNotificationController(mockService)
	router := gin.Default()

	router.POST("/notifications", controller.SendNotification)

	expectedNotification := domain.Notification{
		ID:      1,
		UserID:  1,
		Message: "Test notification",
		SentAt:  time.Now().UTC(),
	}

	mockService.On("SendNotification", 1, "Test notification").Return(expectedNotification, nil)

	req, _ := http.NewRequest("POST", "/notifications", jsonBody(t, map[string]interface{}{
		"user_id": 1,
		"message": "Test notification",
	}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Notification
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedNotification.UserID, response.UserID)
	assert.Equal(t, expectedNotification.Message, response.Message)
	mockService.AssertExpectations(t)
}

func TestGetNotificationsByUserID(t *testing.T) {
	mockService := new(MockNotificationService)
	controller := NewNotificationController(mockService)
	router := gin.Default()

	router.GET("/notifications/:userId", controller.GetNotificationsByUserID)

	expectedNotifications := []domain.Notification{
		{
			ID:      1,
			UserID:  1,
			Message: "Test notification 1",
			SentAt:  time.Now().UTC(),
		},
		{
			ID:      2,
			UserID:  1,
			Message: "Test notification 2",
			SentAt:  time.Now().UTC(),
		},
	}

	mockService.On("GetNotificationsByUserID", 1).Return(expectedNotifications, nil)

	req, _ := http.NewRequest("GET", "/notifications/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.Notification
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedNotifications, response)
	mockService.AssertExpectations(t)
}

func jsonBody(t *testing.T, obj interface{}) *strings.Reader {
	data, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	return strings.NewReader(string(data))
}
