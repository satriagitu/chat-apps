package service

import (
	"testing"
	"time"

	"chat-apps/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock JobRepository
type MockJobRepository struct {
	mock.Mock
}

func (m *MockJobRepository) CreateJob(job domain.Job) (domain.Job, error) {
	args := m.Called(job)
	return args.Get(0).(domain.Job), args.Error(1)
}

func (m *MockJobRepository) GetJobByID(id int) (domain.Job, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Job), args.Error(1)
}

func (m *MockJobRepository) UpdateJobStatus(id int, status string, completedAt time.Time) error {
	args := m.Called(id, status, completedAt)
	return args.Error(0)
}

// Mock NotificationRepository
type MockNotificationsRepository struct {
	mock.Mock
}

func (m *MockNotificationsRepository) CreateNotification(notification domain.Notification) (domain.Notification, error) {
	args := m.Called(notification)
	return args.Get(0).(domain.Notification), args.Error(1)
}

func (m *MockNotificationsRepository) GetNotificationsByUserID(userID int) ([]domain.Notification, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Notification), args.Error(1)
}

// Mock UserRepository
type MockUserRepositorys struct {
	mock.Mock
}

func (m *MockUserRepositorys) ExistsByID(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepositorys) CreateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepositorys) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepositorys) GetUserByID(id int) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

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

func TestGetJobStatus(t *testing.T) {
	mockJobRepo := new(MockJobRepository)
	service := NewJobService(mockJobRepo, nil, nil, nil)

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
