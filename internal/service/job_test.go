package service_test

import (
	"encoding/json"
	"testing"
	"time"

	"chat-apps/internal/domain"
	"chat-apps/internal/repository/mocks"
	"chat-apps/internal/service"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func TestQueueBroadcastNotification(t *testing.T) {
// 	mockJobRepo := new(MockJobRepository)
// 	mockNotificationRepo := new(MockNotificationsRepository)
// 	mockUserRepo := new(MockUserRepositorys)
// 	mockChannel := new(amqp.Channel)

// 	service := NewJobService(mockJobRepo, mockNotificationRepo, mockUserRepo, mockChannel)

// 	// Mock data
// 	job := domain.Job{
// 		ID:       1,
// 		Message:  "Test message",
// 		Status:   "queued",
// 		QueuedAt: time.Now(),
// 	}
// 	mockJobRepo.On("CreateJob", mock.Anything).Return(job, nil)

// 	// Call the method
// 	result, err := service.QueueBroadcastNotification("Test message")

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.Equal(t, job.ID, result.ID)
// 	assert.Equal(t, job.Message, result.Message)
// 	mockJobRepo.AssertExpectations(t)
// }

func TestJobService_QueueBroadcastNotification(t *testing.T) {
	jobRepo := new(mocks.JobRepository)
	notificationRepo := new(mocks.NotificationRepository)
	userRepo := new(mocks.UserRepository)
	rabbitMQChannel := new(mocks.RabbitMQChannel)

	// Create an instance of the JobService with the mocked dependencies
	service := service.NewJobService(jobRepo, notificationRepo, userRepo, rabbitMQChannel)

	// Define the job and expected behavior
	message := "Test Broadcast Message"
	job := domain.Job{
		ID:       1,
		Message:  message,
		Status:   "queued",
		QueuedAt: time.Now(),
	}

	jobRepo.On("CreateJob", mock.AnythingOfType("domain.Job")).Return(job, nil)

	// Mock the Publish method
	payload, _ := json.Marshal(map[string]interface{}{"job_id": job.ID, "message": message})
	rabbitMQChannel.On("Publish", "", "broadcast_queue", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	}).Return(nil)

	// Call the QueueBroadcastNotification method
	result, err := service.QueueBroadcastNotification(message)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, job.ID, result.ID)
	assert.Equal(t, "queued", result.Status)
	jobRepo.AssertExpectations(t)
	rabbitMQChannel.AssertExpectations(t)
}

func TestGetJobStatus(t *testing.T) {
	mockJobRepo := new(mocks.JobRepository)
	service := service.NewJobService(mockJobRepo, nil, nil, nil)

	// Mock data
	job := domain.Job{
		ID:       1,
		Message:  "Test message",
		Status:   "completed",
		QueuedAt: time.Now(),
	}
	mockJobRepo.On("GetJobByID", 1).Return(job, nil)

	// Call the method
	result, err := service.GetJobStatus(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, job.ID, result.ID)
	assert.Equal(t, job.Status, result.Status)
	mockJobRepo.AssertExpectations(t)
}
