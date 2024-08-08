package controller_test

import (
	"chat-apps/internal/controller"
	"chat-apps/internal/domain"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService for testing purposes
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id int) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockUserService := new(MockUserService)
	userController := controller.NewUserController(mockUserService)

	router := gin.Default()
	router.POST("/users", userController.CreateUser)

	user := domain.User{Username: "testuser", Email: "test@example.com", Password: ""}

	mockUserService.On("CreateUser", user).Return(user, nil)

	w := httptest.NewRecorder()
	body := `{"username": "testuser", "email": "test@example.com", "password": "password"}`
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockUserService := new(MockUserService)
	userController := controller.NewUserController(mockUserService)

	router := gin.Default()
	router.GET("/users/:id", userController.GetUserByID)

	user := domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}

	mockUserService.On("GetUserByID", 1).Return(user, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockUserService := new(MockUserService)
	userController := controller.NewUserController(mockUserService)

	router := gin.Default()
	router.GET("/users/:id", userController.GetUserByID)

	mockUserService.On("GetUserByID", 1).Return(domain.User{}, errors.New("user not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserService.AssertExpectations(t)
}
