// repository/conversation_repository_test.go
package repository

import (
	"chat-apps/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateConversation(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(&domain.Conversation{}, &domain.ConversationParticipant{})
	if err != nil {
		t.Fatal("failed to migrate database")
	}

	repo := NewConversationRepository(db)
	participants := []int{1, 2, 3}
	resp, err := repo.CreateConversation(participants)

	assert.NoError(t, err)
	assert.NotZero(t, resp.ID)
	assert.ElementsMatch(t, participants, resp.Participants)
}

func TestGetConversationByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(&domain.Conversation{}, &domain.ConversationParticipant{})
	if err != nil {
		t.Fatal("failed to migrate database")
	}

	repo := NewConversationRepository(db)
	conversation := domain.Conversation{CreatedAt: time.Now()}
	db.Create(&conversation)
	participants := []domain.ConversationParticipant{
		{ConversationID: conversation.ID, UserID: 1},
		{ConversationID: conversation.ID, UserID: 2},
	}
	db.Create(&participants)

	resp, err := repo.GetConversationByID(conversation.ID)

	assert.NoError(t, err)
	assert.Equal(t, conversation.ID, resp.ID)
	assert.ElementsMatch(t, []int{1, 2}, resp.Participants)
}
