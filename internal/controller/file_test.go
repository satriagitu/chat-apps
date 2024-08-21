// controller/file_controller_test.go
package controller_test

import (
	"bytes"
	"chat-apps/internal/controller"
	"chat-apps/internal/domain"
	"chat-apps/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadFile(t *testing.T) {
	mockFileService := new(mocks.FileService)
	fileController := controller.NewFileController(mockFileService)

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
	mockFileService := new(mocks.FileService)
	fileController := controller.NewFileController(mockFileService)

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

func TestUploadFileV2(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful upload", func(t *testing.T) {
		body := domain.File{
			UserID:  1,
			FileURL: "http://example.com/file",
		}
		jsonValue, _ := json.Marshal(body)

		mockFileService := new(mocks.FileService)
		mockFileService.On("UploadFile", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(body, nil)
		fileController := controller.NewFileController(mockFileService)

		r := gin.Default()
		r.POST("/files/upload", fileController.UploadFile)

		req, _ := http.NewRequest(http.MethodPost, "/files/upload", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// start test
		t.Log("[response]:", resp.Body.String())
		var responseBody map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, float64(1), responseBody["user_id"])
		assert.Equal(t, "http://example.com/file", responseBody["file_url"])
	})

	t.Run("invalid json - bad request", func(t *testing.T) {
		mockFileService := new(mocks.FileService)
		fileController := controller.NewFileController(mockFileService)

		r := gin.Default()
		r.POST("/files/upload", fileController.UploadFile)

		req, _ := http.NewRequest(http.MethodPost, "/files/upload", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// start test
		t.Log("[response]:", resp.Body.String())
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "bad request")
	})

	t.Run("service error", func(t *testing.T) {
		body := gin.H{
			"user_id": 1,
			"file":    "http://example.com/file",
		}
		jsonValue, _ := json.Marshal(body)

		mockFileService := new(mocks.FileService)
		mockFileService.On("UploadFile", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(domain.File{}, assert.AnError)
		fileController := controller.NewFileController(mockFileService)

		r := gin.Default()
		r.POST("/files/upload", fileController.UploadFile)

		req, _ := http.NewRequest(http.MethodPost, "/files/upload", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// start test
		t.Log("[response]:", resp.Body.String())
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "internal server error")
	})
}

func TestGetFileV2(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid file id - bad request", func(t *testing.T) {
		mockFileService := new(mocks.FileService)
		fileController := controller.NewFileController(mockFileService)

		r := gin.Default()
		r.GET("/files/:id", fileController.GetFileByID)

		req, _ := http.NewRequest(http.MethodGet, "/files/abc", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// start test
		t.Log("[response]:", resp.Body.String())
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "Invalid file ID")
	})

	t.Run("service error - internal server error", func(t *testing.T) {
		mockFileService := new(mocks.FileService)
		mockFileService.On("GetFileByID", mock.AnythingOfType("int")).Return(domain.File{}, assert.AnError)
		fileController := controller.NewFileController(mockFileService)
		r := gin.Default()
		r.GET("/files/:id", fileController.GetFileByID)

		req, _ := http.NewRequest(http.MethodGet, "/files/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// start test
		t.Log("[response]:", resp.Body.String())
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "internal server error")
	})
}
