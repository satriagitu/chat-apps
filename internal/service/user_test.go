package service_test

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository for testing purposes
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id int) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByID(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	// Initial user input with plain text password
	inputUser := domain.User{Username: "testuser", Email: "test@example.com", Password: "password"}

	// Hash the password as it would be in the service method
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(inputUser.Password), bcrypt.DefaultCost)

	// Prepare the expected user with the hashed password
	expectedUser := inputUser
	expectedUser.Password = string(hashedPassword)

	// Mock the repository to return the user with hashed password
	mockUserRepo.On("CreateUser", mock.AnythingOfType("domain.User")).Return(expectedUser, nil)

	// Call the CreateUser service method
	createdUser, err := userService.CreateUser(inputUser)

	// Check if there was no error during the service call
	assert.NoError(t, err)

	// Check if the username and email are as expected
	assert.Equal(t, inputUser.Username, createdUser.Username)
	assert.Equal(t, inputUser.Email, createdUser.Email)

	// Instead of comparing the actual password strings, verify the hash
	err = bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte("password"))
	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	user := domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}

	mockUserRepo.On("GetUserByID", 1).Return(user, nil)

	foundUser, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, user, foundUser)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	mockUserRepo.On("GetUserByID", 1).Return(domain.User{}, errors.New("user not found"))

	_, err := userService.GetUserByID(1)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockUserRepo.AssertExpectations(t)
}
