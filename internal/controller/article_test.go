package controller_test

import (
	"chat-apps/internal/controller"
	"chat-apps/internal/domain"
	"chat-apps/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestArticleController_GetArticleList(t *testing.T) {
	// setup
	mockArticles := []domain.ArticleList{
		{
			ID:           1,
			Menu:         "Tech",
			SubMenu:      "Programming",
			Title:        "Go Unit Testing",
			Image:        "image.png",
			Likes:        10,
			CommentCount: 5,
			TimeAgo:      "2 days ago",
		},
	}

	mockArticleService := new(mocks.ArticleService)
	mockArticleService.On("GetArticleList").Return(mockArticles, nil)
	articleController := controller.NewArticleController(mockArticleService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/article-list", articleController.GetArticleList)

	// start test
	req, _ := http.NewRequest(http.MethodGet, "/article-list", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[
		{
			"id": 1,
			"menu": "Tech",
			"sub_menu": "Programming",
			"title": "Go Unit Testing",
			"image": "image.png",
			"likes": 10,
			"comment_count": 5,
			"time_ago": "2 days ago"
		}
	]`, w.Body.String())
	mockArticleService.AssertExpectations(t)
}

func TestArticleController_GetArticleList_Error(t *testing.T) {
	// Setup
	mockArticleService := new(mocks.ArticleService)
	mockArticleService.On("GetArticleList").Return(nil, assert.AnError)
	articleController := controller.NewArticleController(mockArticleService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/article-list", articleController.GetArticleList)

	// Start test
	req, _ := http.NewRequest(http.MethodGet, "/article-list", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "internal server error"}`, w.Body.String())
	mockArticleService.AssertExpectations(t)
}
