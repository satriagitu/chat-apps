package service

import (
	"errors"
	"testing"
	"time"

	"chat-apps/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) CreateNotification(notification domain.Notification) (domain.Notification, error) {
	args := m.Called(notification)
	return args.Get(0).(domain.Notification), args.Error(1)
}

func (m *MockNotificationRepository) GetNotificationsByUserID(userID int) ([]domain.Notification, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Notification), args.Error(1)
}

// MockUserRepository
type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) ExistsByID(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUsersRepository) CreateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUsersRepository) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUsersRepository) GetUserByID(id int) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func TestSendNotification_UserExists(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)
	mockUserRepo := new(MockUsersRepository)
	service := NewNotificationService(mockNotificationRepo, mockUserRepo)

	// Mock data
	userID := 1
	message := "Test notification"
	expectedNotification := domain.Notification{
		UserID:  userID,
		Message: message,
		SentAt:  time.Now().Truncate(time.Second),
	}

	// Set up expectations
	mockUserRepo.On("ExistsByID", userID).Return(true, nil)
	mockNotificationRepo.On("CreateNotification", mock.Anything).Return(expectedNotification, nil)

	// Call the method
	result, err := service.SendNotification(userID, message)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedNotification.ID, result.ID)
	assert.Equal(t, expectedNotification.UserID, result.UserID)
	assert.Equal(t, expectedNotification.Message, result.Message)
	assert.WithinDuration(t, expectedNotification.SentAt, result.SentAt, time.Second)
	mockUserRepo.AssertExpectations(t)
	mockNotificationRepo.AssertExpectations(t)
}

func TestSendNotification_UserDoesNotExist(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)
	mockUserRepo := new(MockUsersRepository)
	service := NewNotificationService(mockNotificationRepo, mockUserRepo)

	// Mock data
	userID := 1
	message := "Test notification"

	// Set up expectations
	mockUserRepo.On("ExistsByID", userID).Return(false, nil)

	// Call the method
	result, err := service.SendNotification(userID, message)

	// Assert results
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, domain.Notification{}, result)
	mockUserRepo.AssertExpectations(t)
	mockNotificationRepo.AssertExpectations(t)
}

func TestSendNotification_ErrorCheckingUser(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)
	mockUserRepo := new(MockUsersRepository)
	service := NewNotificationService(mockNotificationRepo, mockUserRepo)

	// Mock data
	userID := 1
	message := "Test notification"

	// Set up expectations
	mockUserRepo.On("ExistsByID", userID).Return(false, errors.New("database error"))

	// Call the method
	result, err := service.SendNotification(userID, message)

	// Assert results
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Equal(t, domain.Notification{}, result)
	mockUserRepo.AssertExpectations(t)
	mockNotificationRepo.AssertExpectations(t)
}

func TestGetNotificationsByUserID(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)
	mockUserRepo := new(MockUsersRepository)
	service := NewNotificationService(mockNotificationRepo, mockUserRepo)

	// Mock data
	userID := 1
	notifications := []domain.Notification{
		{UserID: userID, Message: "Test notification 1", SentAt: time.Now().UTC()},
		{UserID: userID, Message: "Test notification 2", SentAt: time.Now().UTC()},
	}

	// Set up expectations
	mockNotificationRepo.On("GetNotificationsByUserID", userID).Return(notifications, nil)

	// Call the method
	result, err := service.GetNotificationsByUserID(userID)

	// Assert results
	assert.NoError(t, err)
	assert.ElementsMatch(t, notifications, result)
	mockNotificationRepo.AssertExpectations(t)
}
