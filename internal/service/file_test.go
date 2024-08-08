// service/file_service_test.go
package service

import (
	"chat-apps/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) ExistsByID(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id int) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

// MockFileRepository
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) UploadFile(file domain.File) (domain.File, error) {
	args := m.Called(file)
	return args.Get(0).(domain.File), args.Error(1)
}

func (m *MockFileRepository) GetFileByID(id int) (domain.File, error) {
	args := m.Called(id)
	return args.Get(0).(domain.File), args.Error(1)
}

func TestUploadFile(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockFileRepo := new(MockFileRepository)
	fileService := NewFileService(mockFileRepo, mockUserRepo)

	userID := 1
	fileURL := "http://example.com/file"

	mockUserRepo.On("ExistsByID", userID).Return(true, nil)
	mockFileRepo.On("UploadFile", domain.File{UserID: userID, FileURL: fileURL}).Return(domain.File{ID: 1, UserID: userID, FileURL: fileURL}, nil)

	file, err := fileService.UploadFile(userID, fileURL)

	assert.NoError(t, err)
	assert.Equal(t, userID, file.UserID)
	assert.Equal(t, fileURL, file.FileURL)
	mockUserRepo.AssertExpectations(t)
	mockFileRepo.AssertExpectations(t)
}

func TestUploadFile_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockFileRepo := new(MockFileRepository)
	fileService := NewFileService(mockFileRepo, mockUserRepo)

	userID := 1
	fileURL := "http://example.com/file"

	mockUserRepo.On("ExistsByID", userID).Return(false, nil)

	file, err := fileService.UploadFile(userID, fileURL)

	assert.Error(t, err)
	assert.Empty(t, file)
}

func TestGetFileByID(t *testing.T) {
	mockFileRepo := new(MockFileRepository)
	fileService := NewFileService(mockFileRepo, nil)

	fileID := 1
	expectedFile := domain.File{ID: fileID, UserID: 1, FileURL: "http://example.com/file"}

	mockFileRepo.On("GetFileByID", fileID).Return(expectedFile, nil)

	file, err := fileService.GetFileByID(fileID)

	assert.NoError(t, err)
	assert.Equal(t, expectedFile, file)
	mockFileRepo.AssertExpectations(t)
}
