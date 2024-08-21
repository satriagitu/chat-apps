// service/file_service_test.go
package service

import (
	"chat-apps/internal/domain"
	"testing"

	"chat-apps/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockFileRepo := new(mocks.FileRepository)
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
	mockUserRepo := new(mocks.UserRepository)
	mockFileRepo := new(mocks.FileRepository)
	fileService := NewFileService(mockFileRepo, mockUserRepo)

	userID := 1
	fileURL := "http://example.com/file"

	mockUserRepo.On("ExistsByID", userID).Return(false, nil)

	file, err := fileService.UploadFile(userID, fileURL)

	assert.Error(t, err)
	assert.Empty(t, file)
}

func TestGetFileByID(t *testing.T) {
	mockFileRepo := new(mocks.FileRepository)
	fileService := NewFileService(mockFileRepo, nil)

	fileID := 1
	expectedFile := domain.File{ID: fileID, UserID: 1, FileURL: "http://example.com/file"}

	mockFileRepo.On("GetFileByID", fileID).Return(expectedFile, nil)

	file, err := fileService.GetFileByID(fileID)

	assert.NoError(t, err)
	assert.Equal(t, expectedFile, file)
	mockFileRepo.AssertExpectations(t)
}

func TestExistByID_Error(t *testing.T) {
	t.Run("test negatif - exist by id", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)

		userID := 1
		fileURL := "image.png"

		mockUserRepo.On("ExistsByID", userID).Return(false, assert.AnError)
		mockFileRepo.On("UploadFile", domain.File{
			UserID:  userID,
			FileURL: fileURL,
		}).Return(domain.File{}, assert.AnError)

		fileService := NewFileService(mockFileRepo, mockUserRepo)
		file, err := fileService.UploadFile(userID, fileURL)

		// test
		assert.NotNil(t, err)
		assert.Equal(t, domain.File{}, file)
	})
}
