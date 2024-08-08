package repository_test

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&domain.User{})

	userRepo := repository.NewUserRepository(db)

	user := domain.User{Username: "testuser", Email: "test@example.com", Password: "password"}

	createdUser, err := userRepo.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, createdUser.Username)
	assert.Equal(t, user.Email, createdUser.Email)
}

func TestGetUserByID(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&domain.User{})

	userRepo := repository.NewUserRepository(db)

	user := domain.User{Username: "testuser", Email: "test@example.com", Password: "password"}
	db.Create(&user)

	foundUser, err := userRepo.GetUserByID(user.ID)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&domain.User{})

	userRepo := repository.NewUserRepository(db)

	_, err := userRepo.GetUserByID(1)

	assert.Error(t, err)
}
