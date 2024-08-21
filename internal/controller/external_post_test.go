package controller_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"chat-apps/internal/controller"
	"chat-apps/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestExternalAPIController_GetPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("test positif - get posts from JSONPlaceholder", func(t *testing.T) {
		url := "https://jsonplaceholder.typicode.com/"
		externalAPIController := controller.NewExternalPostController(url)

		// Setup router
		r := gin.Default()
		r.GET("/external-post", externalAPIController.GetExternalPosts)

		// Create request and response recorder
		req, _ := http.NewRequest(http.MethodGet, "/external-post", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		body, _ := io.ReadAll(w.Body)
		data := []domain.ExternalPost{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Log("error unmarshall")
		}

		// check total
		assert.Greater(t, len(data), 0)
		assert.Equal(t, len(body), 24519)

		// check first post
		firstPost := data[0]
		assert.Equal(t, 1, firstPost.Id)
		assert.Equal(t, 1, firstPost.UserId)

		// check element
		for _, post := range data {
			assert.NotZero(t, post.UserId, "UserId should not be zero")
			assert.NotZero(t, post.Id, "Id should not be zero")
			assert.NotEmpty(t, post.Title, "Title should not be empty")
			assert.NotEmpty(t, post.Body, "Body should not be empty")
		}
	})
}
