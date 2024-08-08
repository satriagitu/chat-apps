// repository/file_repository_test.go
package repository

import (
	"chat-apps/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUploadFile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	db.AutoMigrate(&domain.File{})

	fileRepo := NewFileRepository(db)
	file := domain.File{
		UserID:  1,
		FileURL: "http://example.com/file",
	}

	createdFile, err := fileRepo.UploadFile(file)
	if err != nil {
		t.Fatalf("Failed to upload file: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, file.UserID, createdFile.UserID)
	assert.Equal(t, file.FileURL, createdFile.FileURL)
	assert.WithinDuration(t, time.Now(), createdFile.UploadedAt, time.Second)
}

func TestGetFileByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	db.AutoMigrate(&domain.File{})

	fileRepo := NewFileRepository(db)
	file := domain.File{
		UserID:  1,
		FileURL: "http://example.com/file",
	}
	createdFile, _ := fileRepo.UploadFile(file)

	retrievedFile, err := fileRepo.GetFileByID(createdFile.ID)
	if err != nil {
		t.Fatalf("Failed to get file by ID: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, createdFile.ID, retrievedFile.ID)
	assert.Equal(t, createdFile.UserID, retrievedFile.UserID)
	assert.Equal(t, createdFile.FileURL, retrievedFile.FileURL)
}
