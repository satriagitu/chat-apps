package repository

import (
	"testing"
	"time"

	"chat-apps/internal/domain"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&domain.Job{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateJob(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewJobRepository(db)
	job := domain.Job{
		Message:  "Test job",
		Status:   "queued",
		QueuedAt: time.Now(),
	}

	createdJob, err := repo.CreateJob(job)
	assert.NoError(t, err)
	assert.Equal(t, job.Message, createdJob.Message)
	assert.Equal(t, job.Status, createdJob.Status)
}

func TestGetJobByID(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewJobRepository(db)
	job := domain.Job{
		Message:  "Test job",
		Status:   "queued",
		QueuedAt: time.Now(),
	}
	createdJob, err := repo.CreateJob(job)
	assert.NoError(t, err)

	retrievedJob, err := repo.GetJobByID(createdJob.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdJob.ID, retrievedJob.ID)
	assert.Equal(t, createdJob.Message, retrievedJob.Message)
}

func TestUpdateJobStatus(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewJobRepository(db)
	job := domain.Job{
		Message:  "Test job",
		Status:   "queued",
		QueuedAt: time.Now(),
	}
	createdJob, err := repo.CreateJob(job)
	assert.NoError(t, err)

	completedAt := time.Now()
	err = repo.UpdateJobStatus(createdJob.ID, "completed", completedAt)
	assert.NoError(t, err)

	updatedJob, err := repo.GetJobByID(createdJob.ID)
	assert.NoError(t, err)
	assert.Equal(t, "completed", updatedJob.Status)
	assert.NotNil(t, updatedJob.CompletedAt)
}
