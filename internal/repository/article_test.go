package repository_test

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetArticleList(t *testing.T) {
	// setup
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(&domain.Article{}, &domain.Menu{}, &domain.SubMenu{}, &domain.Comment{})
	if err != nil {
		t.Fatal("failed to migrate database")
	}

	menu := domain.Menu{Name: "Tech"}
	subMenu := domain.SubMenu{Name: "Programming"}
	article := domain.Article{
		MenuID:    1,
		SubMenuID: 1,
		Title:     "Go Unit Testing",
		Image:     "image.png",
		CreatedAt: time.Now().Add(-48 * time.Hour),
		Likes:     10,
	}
	comment := domain.Comment{
		ArticleID: 1,
		Content:   "Nice article!",
	}

	db.Create(&menu)
	db.Create(&subMenu)
	db.Create(&article)
	db.Create(&comment)

	repo := repository.NewArtikelRepository(db)
	articles, err := repo.GetArticleList()

	// start test
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
	assert.Equal(t, "Go Unit Testing", articles[0].Title)
	assert.Equal(t, 1, articles[0].CommentCount)
	assert.Equal(t, 10, articles[0].Likes)
	assert.Equal(t, "Tech", articles[0].Menu)
	assert.Equal(t, "image.png", articles[0].Image)
	assert.Equal(t, "Programming", articles[0].SubMenu)

	expectedTime := time.Now().Add(-48 * time.Hour)
	actualTime := articles[0].CreatedAt

	assert.WithinDuration(t, expectedTime, actualTime, time.Second, "CreatedAt should be within the expected time range")
}
