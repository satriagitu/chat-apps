// controller/file_controller_test.go
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

// MockFileService
type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) UploadFile(userID int, fileURL string) (domain.File, error) {
	args := m.Called(userID, fileURL)
	return args.Get(0).(domain.File), args.Error(1)
}

func (m *MockFileService) GetFileByID(id int) (domain.File, error) {
	args := m.Called(id)
	return args.Get(0).(domain.File), args.Error(1)
}

func TestUploadFile(t *testing.T) {
	mockFileService := new(MockFileService)
	fileController := NewFileController(mockFileService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/upload", fileController.UploadFile)

	file := domain.File{
		ID:         1,
		UserID:     1,
		FileURL:    "http://example.com/file",
		UploadedAt: time.Now(),
	}
	mockFileService.On("UploadFile", 1, "http://example.com/file").Return(file, nil)

	jsonData, _ := json.Marshal(map[string]interface{}{
		"user_id": 1,
		"file":    "http://example.com/file",
	})
	req, _ := http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, float64(1), response["user_id"])
	assert.Equal(t, "http://example.com/file", response["file_url"])
	mockFileService.AssertExpectations(t)
}

func TestGetFileByID(t *testing.T) {
	mockFileService := new(MockFileService)
	fileController := NewFileController(mockFileService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/file/:id", fileController.GetFileByID)

	file := domain.File{
		ID:         1,
		UserID:     1,
		FileURL:    "http://example.com/file",
		UploadedAt: time.Now(),
	}
	mockFileService.On("GetFileByID", 1).Return(file, nil)

	req, _ := http.NewRequest("GET", "/file/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, float64(1), response["user_id"])
	assert.Equal(t, "http://example.com/file", response["file_url"])
	mockFileService.AssertExpectations(t)
}
