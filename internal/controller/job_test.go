package controller

import (
	"strings"
	"testing"
	"time"

	"chat-apps/internal/domain"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock JobService
type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) QueueBroadcastNotification(message string) (domain.Job, error) {
	args := m.Called(message)
	return args.Get(0).(domain.Job), args.Error(1)
}

func (m *MockJobService) GetJobStatus(id int) (domain.Job, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Job), args.Error(1)
}

func TestQueueBroadcastNotification(t *testing.T) {
	mockJobService := new(MockJobService)
	controller := NewJobController(mockJobService)

	// Setup mock
	job := domain.Job{
		ID:       1,
		Message:  "Test message",
		Status:   "queued",
		QueuedAt: time.Now(),
	}
	mockJobService.On("QueueBroadcastNotification", "Test message").Return(job, nil)

	// Create request
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/jobs/broadcast", controller.QueueBroadcastNotification)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs/broadcast", strings.NewReader(`{"message":"Test message"}`))
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":1`)
	assert.Contains(t, w.Body.String(), `"message":"Test message"`)
	mockJobService.AssertExpectations(t)
}

func TestGetJobStatus(t *testing.T) {
	mockJobService := new(MockJobService)
	controller := NewJobController(mockJobService)

	// Setup mock
	job := domain.Job{
		ID:       1,
		Message:  "Test message",
		Status:   "completed",
		QueuedAt: time.Now(),
	}
	mockJobService.On("GetJobStatus", 1).Return(job, nil)

	// Create request
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/jobs/:id", controller.GetJobStatus)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":1`)
	assert.Contains(t, w.Body.String(), `"status":"completed"`)
	mockJobService.AssertExpectations(t)
}
