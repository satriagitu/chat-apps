package service_test

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository/mocks"
	"chat-apps/internal/service"
	"chat-apps/internal/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArtikelService_GetArticleList(t *testing.T) {
	t.Run("test normal service get article list", func(t *testing.T) {

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

		mockArticleRepo.On("GetArticleList", mock.AnythingOfType("string")).Return(mockArticles, nil)

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
		result, err := articleService.GetArticleList("")
		t.Log("[result]:", result)

		t.Run("cek no error and expect result with time ago", func(t *testing.T) {
			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, result)
			mockArticleRepo.AssertExpectations(t)
		})
	})
}

func TestArtikelService_GetArticleList_Error(t *testing.T) {
	t.Run("test error service get list article", func(t *testing.T) {
		mockArticleRepo := new(mocks.ArticleRepository)

		// Mengatur mock agar GetArticleList mengembalikan error
		mockArticleRepo.On("GetArticleList", mock.AnythingOfType("string")).Return(nil, assert.AnError)

		articleService := service.NewArtikelService(mockArticleRepo)
		articles, err := articleService.GetArticleList("golang")
		articles2, err2 := articleService.GetArticleList("")

		// Memastikan error diterima dan articles nil
		assert.Nil(t, articles)
		assert.EqualError(t, err, assert.AnError.Error())
		assert.Nil(t, articles2)
		assert.EqualError(t, err2, assert.AnError.Error())

		mockArticleRepo.AssertExpectations(t)
	})

}
