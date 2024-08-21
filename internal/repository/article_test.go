package repository_test

import (
	"chat-apps/internal/repository"
	"chat-apps/internal/util"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetArticleList_Success(t *testing.T) {
	db, mock := util.SetupMockDB()

	t.Run("test normal get article list", func(t *testing.T) {
		createdAt := time.Now()
		// Mock the expected query result
		rows := sqlmock.NewRows([]string{"id", "menu", "sub_menu", "title", "image", "likes", "comment_count", "created_at"}).
			AddRow(2, "Menu2", "SubMenu2", "Title2", "Image2", 10, 5, createdAt).
			AddRow(1, "Menu1", "SubMenu1", "Title1", "Image1", 15, 3, createdAt)

		mock.ExpectQuery(`SELECT a.id, m.name AS menu, sm.name AS sub_menu, a.title, a.image, a.likes, c.comment_count, 
							a.created_at
							FROM articles AS a
							LEFT JOIN menus m ON a.menu_id = m.id
							LEFT JOIN sub_menus sm ON a.sub_menu_id = sm.id
							LEFT JOIN (
								SELECT article_id, COUNT(id) AS comment_count
								FROM comments
								GROUP BY article_id
								) c ON a.id = c.article_id
								ORDER BY a.id DESC`).
			WillReturnRows(rows)

		// Call the method
		repo := repository.NewArtikelRepository(db)
		articles, err := repo.GetArticleList("")

		t.Run("test without query search", func(t *testing.T) {
			// Assertions
			assert.NoError(t, err)
			assert.Len(t, articles, 2)
			assert.Equal(t, "Title2", articles[0].Title)
			assert.Equal(t, "Title1", articles[1].Title)
		})

	})
}

func TestGetArticleListWithSearch(t *testing.T) {
	db, mock := util.SetupMockDB()

	t.Run("test normal get article list", func(t *testing.T) {

		createdAt := time.Now()
		rows2 := sqlmock.NewRows([]string{"id", "menu", "sub_menu", "title", "image", "likes", "comment_count", "created_at"}).
			AddRow(5, "Education", "Online Learning", "Best Online Learning Platforms", "https://example.com/learning.png", 130, 0, createdAt)

		search := "online"
		searchQuery := "%" + search + "%" // Menambahkan wildcard untuk LIKE
		query := fmt.Sprintf(`SELECT a.id, m.name AS menu, sm.name AS sub_menu, a.title, a.image, a.likes, c.comment_count, a.created_at 
		FROM articles AS a
		LEFT JOIN menus m ON a.menu_id = m.id
		LEFT JOIN sub_menus sm ON a.sub_menu_id = sm.id
							LEFT JOIN (
								SELECT article_id, COUNT(id) AS comment_count
								FROM comments
								GROUP BY article_id
								) c ON a.id = c.article_id
								WHERE a.title ILIKE '%s' OR m.name ILIKE '%s' OR sm.name ILIKE '%s'
								ORDER BY a.id DESC`, searchQuery, searchQuery, searchQuery)

		mock.ExpectQuery(query).
			WillReturnRows(rows2)

		// Call the method
		repo2 := repository.NewArtikelRepository(db)
		articles2, err2 := repo2.GetArticleList(search)

		t.Run("test with query search", func(t *testing.T) {
			// Assertions
			assert.NoError(t, err2)
			assert.Len(t, articles2, 1)
			assert.Equal(t, "Best Online Learning Platforms", articles2[0].Title)
		})
	})
}

func TestGetArticleList_Error(t *testing.T) {
	db, mock := util.SetupMockDB()

	t.Run("test error get article list", func(t *testing.T) {
		repo := repository.NewArtikelRepository(db)

		// Simulate an error during query execution
		mock.ExpectQuery("SELECT a.id, m.name AS menu, sm.name AS sub_menu, a.title, a.image, a.likes, c.comment_count, a.created_at").
			WillReturnError(assert.AnError)

		// Call the method
		articles, err := repo.GetArticleList("")
		t.Run("test response error and article nil", func(t *testing.T) {
			// Assertions
			assert.Error(t, err)
			assert.Nil(t, articles)
		})
	})
}
