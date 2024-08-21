package controller_test

import (
	"bytes"
	"chat-apps/internal/controller"
	"chat-apps/internal/service/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCacheController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SetCache - success", func(t *testing.T) {

		// Initialize mock CacheService
		mockCacheService := new(mocks.CacheService)
		cacheController := controller.NewCacheController(mockCacheService)

		// Create a test context
		r := gin.Default()
		// Set up the POST route
		r.POST("/cache-redis", cacheController.SetCache)

		// Prepare the request payload
		body := map[string]string{"key": "mykey", "value": "myvalue"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/cache-redis", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Mock the CacheData method
		mockCacheService.On("CacheData", mock.AnythingOfType("*gin.Context"), "mykey", "myvalue").Return(nil)

		// Execute the request
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Validate the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"status":"success"}`, w.Body.String())
	})

	t.Run("SetCache - bad request", func(t *testing.T) {
		// Create a test context
		r := gin.Default()

		// Initialize mock CacheService
		mockCacheService := new(mocks.CacheService)
		cacheController := controller.NewCacheController(mockCacheService)

		// Set up the POST route
		r.POST("/cache-redis", cacheController.SetCache)

		// Prepare the request payload with a missing "value"
		req, _ := http.NewRequest(http.MethodPost, "/cache-redis", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Execute the request
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		t.Log("w:", w)
		// Validate the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid request"}`, w.Body.String())
	})

	t.Run("SetCache - internal server error", func(t *testing.T) {
		// Create a test context
		r := gin.Default()

		// Initialize mock CacheService
		mockCacheService := new(mocks.CacheService)
		cacheController := controller.NewCacheController(mockCacheService)

		// Set up the POST route
		r.POST("/cache-redis", cacheController.SetCache)

		// Prepare the request payload with a missing "value"
		body := map[string]string{"key": "mykey", "value": "myvalue"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/cache-redis", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		mockCacheService.On("CacheData", mock.AnythingOfType("*gin.Context"), "mykey", "myvalue").Return(assert.AnError)
		// Execute the request
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		t.Log("w:", w)
		// Validate the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Failed to set cache"}`, w.Body.String())
	})

	t.Run("GetCache - success", func(t *testing.T) {
		// Initialize mock CacheService
		mockCacheService := new(mocks.CacheService)
		cacheController := controller.NewCacheController(mockCacheService)

		// setup return mock
		mockCacheService.On("RetrieveData", mock.AnythingOfType("*gin.Context"), "mykey").Return("myvalue", nil)

		// setup route
		r := gin.Default()
		r.GET("/cache-redis/:key", cacheController.GetCache)

		// prepare request
		req, _ := http.NewRequest(http.MethodGet, "/cache-redis/mykey", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		//check log response
		t.Log("[response]:", w)
		// validate the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"value":"myvalue"}`, w.Body.String())
	})

	t.Run("GetCache - internal server error", func(t *testing.T) {
		// Initialize mock CacheService
		mockCacheService := new(mocks.CacheService)
		cacheController := controller.NewCacheController(mockCacheService)

		// setup return mock
		mockCacheService.On("RetrieveData", mock.AnythingOfType("*gin.Context"), "mykey").Return("", assert.AnError)

		// setup route
		r := gin.Default()
		r.GET("/cache-redis/:key", cacheController.GetCache)

		// prepare request
		req, _ := http.NewRequest(http.MethodGet, "/cache-redis/mykey", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		//check log response
		t.Log("[response]:", w)
		// validate the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Failed to get cache"}`, w.Body.String())
	})

	t.Run("GetCache - not found", func(t *testing.T) {
		// Initialize mock CacheService
		mockCacheService := new(mocks.CacheService)
		cacheController := controller.NewCacheController(mockCacheService)

		// setup return mock
		mockCacheService.On("RetrieveData", mock.AnythingOfType("*gin.Context"), "mykey").Return("", nil)

		// setup route
		r := gin.Default()
		r.GET("/cache-redis/:key", cacheController.GetCache)

		// prepare request
		req, _ := http.NewRequest(http.MethodGet, "/cache-redis/mykey", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		//check log response
		t.Log("[response]:", w)
		// validate the response
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error": "Key not found"}`, w.Body.String())
	})
}
