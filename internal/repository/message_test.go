// repository/message_repository_test.go
package repository

import (
	"chat-apps/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateMessage(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&domain.Message{})

	repo := NewMessageRepository(db)

	message := domain.Message{
		ConversationID: 1,
		SenderID:       2,
		Content:        "Hello",
	}

	createdMessage, err := repo.CreateMessage(message)
	assert.NoError(t, err)
	assert.Equal(t, message.ConversationID, createdMessage.ConversationID)
	assert.Equal(t, message.SenderID, createdMessage.SenderID)
	assert.Equal(t, message.Content, createdMessage.Content)
	assert.WithinDuration(t, time.Now(), createdMessage.SentAt, time.Second) // Check that SentAt is close to now
}

func TestGetMessagesByConversationID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&domain.Message{})

	repo := NewMessageRepository(db)

	// Generate a fixed time for consistency
	now := time.Now().UTC()

	message1 := domain.Message{
		ConversationID: 1,
		SenderID:       2,
		Content:        "Hello",
		SentAt:         now,
	}

	message2 := domain.Message{
		ConversationID: 1,
		SenderID:       3,
		Content:        "Hi",
		SentAt:         now,
	}

	// Create messages
	_, err = repo.CreateMessage(message1)
	if err != nil {
		t.Fatalf("failed to create message1: %v", err)
	}

	_, err = repo.CreateMessage(message2)
	if err != nil {
		t.Fatalf("failed to create message2: %v", err)
	}

	// Get messages
	messages, err := repo.GetMessagesByConversationID(1)
	assert.NoError(t, err)

	// Use a tolerance when checking time
	for i, msg := range []domain.Message{message1, message2} {
		assert.Equal(t, msg.ConversationID, messages[i].ConversationID)
		assert.Equal(t, msg.SenderID, messages[i].SenderID)
		assert.Equal(t, msg.Content, messages[i].Content)
		assert.WithinDuration(t, msg.SentAt, messages[i].SentAt, time.Second)
	}
}
