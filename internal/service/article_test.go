package service_test

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository/mocks"
	"chat-apps/internal/service"
	"chat-apps/internal/util"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArtikelService_GetArticleList(t *testing.T) {
	// Arrange
	mockArticleRepo := new(mocks.ArticleRepository)
	articleService := service.NewArtikelService(mockArticleRepo)

	mockArticles := []domain.ArticleList{
		{
			ID:           1,
			Menu:         "Tech",
			SubMenu:      "Programming",
			Title:        "Go Unit Testing",
			Image:        "image.png",
			Likes:        10,
			CommentCount: 5,
			CreatedAt:    time.Now().Add(-48 * time.Hour), // 2 days ago
		},
	}

	mockArticleRepo.On("GetArticleList").Return(mockArticles, nil)

	expectedResult := []domain.ArticleList{
		{
			ID:           1,
			Menu:         "Tech",
			SubMenu:      "Programming",
			Title:        "Go Unit Testing",
			Image:        "image.png",
			Likes:        10,
			CommentCount: 5,
			TimeAgo:      util.TimeAgo(mockArticles[0].CreatedAt),
		},
	}

	// Act
	result, err := articleService.GetArticleList()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockArticleRepo.AssertExpectations(t)
}

func TestArtikelService_GetArticleList_Error(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockError := errors.New("something went wrong")

	// Mengatur mock agar GetArticleList mengembalikan error
	mockArticleRepo.On("GetArticleList").Return(nil, mockError)

	articleService := service.NewArtikelService(mockArticleRepo)
	articles, err := articleService.GetArticleList()

	// Memastikan error diterima dan articles nil
	assert.Nil(t, articles)
	assert.EqualError(t, err, mockError.Error())

	mockArticleRepo.AssertExpectations(t)
}
