package repository

import (
	"time"

	"chat-apps/internal/domain"

	"gorm.io/gorm"
)

type JobRepository interface {
	CreateJob(job domain.Job) (domain.Job, error)
	GetJobByID(id int) (domain.Job, error)
	UpdateJobStatus(id int, status string, completedAt time.Time) error
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) CreateJob(job domain.Job) (domain.Job, error) {
	if err := r.db.Create(&job).Error; err != nil {
		return domain.Job{}, err
	}
	return job, nil
}

func (r *jobRepository) GetJobByID(id int) (domain.Job, error) {
	var job domain.Job
	if err := r.db.First(&job, id).Error; err != nil {
		return domain.Job{}, err
	}
	return job, nil
}

func (r *jobRepository) UpdateJobStatus(id int, status string, completedAt time.Time) error {
	return r.db.Model(&domain.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"completed_at": completedAt,
	}).Error
}
