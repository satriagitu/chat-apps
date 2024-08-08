package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"chat-apps/internal/domain"
)

func TestCreateNotification(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&domain.Notification{})

	repo := NewNotificationRepository(db)

	notification := domain.Notification{
		UserID:  1,
		Message: "Test notification",
		SentAt:  time.Now().UTC(),
	}

	createdNotification, err := repo.CreateNotification(notification)
	assert.NoError(t, err)
	assert.Equal(t, notification.UserID, createdNotification.UserID)
	assert.Equal(t, notification.Message, createdNotification.Message)
	assert.WithinDuration(t, notification.SentAt, createdNotification.SentAt, time.Second)
}

func TestGetNotificationsByUserID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&domain.Notification{})

	repo := NewNotificationRepository(db)

	// Use fixed time for consistency
	now := time.Now().UTC().Truncate(time.Second)

	notification1 := domain.Notification{
		UserID:  1,
		Message: "Test notification 1",
		SentAt:  now,
	}

	notification2 := domain.Notification{
		UserID:  1,
		Message: "Test notification 2",
		SentAt:  now,
	}

	createdNotification1, err := repo.CreateNotification(notification1)
	if err != nil {
		t.Fatalf("failed to create notification1: %v", err)
	}

	createdNotification2, err := repo.CreateNotification(notification2)
	if err != nil {
		t.Fatalf("failed to create notification2: %v", err)
	}

	notifications, err := repo.GetNotificationsByUserID(1)
	assert.NoError(t, err)

	// Prepare expected notifications with IDs assigned by the database
	expectedNotifications := []domain.Notification{
		{ID: createdNotification1.ID, UserID: 1, Message: "Test notification 1", SentAt: now},
		{ID: createdNotification2.ID, UserID: 1, Message: "Test notification 2", SentAt: now},
	}

	// Normalize the time fields for comparison
	for i := range notifications {
		notifications[i].SentAt = notifications[i].SentAt.UTC().Truncate(time.Second)
		expectedNotifications[i].SentAt = expectedNotifications[i].SentAt.UTC().Truncate(time.Second)
	}

	// Check only fields that matter (excluding IDs)
	assert.ElementsMatch(t, expectedNotifications, notifications)
}
