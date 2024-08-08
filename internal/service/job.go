package service

import (
	"encoding/json"
	"time"

	"chat-apps/internal/domain"
	"chat-apps/internal/repository"

	"github.com/streadway/amqp"
)

type JobService interface {
	QueueBroadcastNotification(message string) (domain.Job, error)
	GetJobStatus(id int) (domain.Job, error)
}

type jobService struct {
	jobRepo          repository.JobRepository
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	rabbitMQChannel  *amqp.Channel
}

func NewJobService(jr repository.JobRepository, nr repository.NotificationRepository, ur repository.UserRepository, ch *amqp.Channel) JobService {
	return &jobService{jobRepo: jr, notificationRepo: nr, userRepo: ur, rabbitMQChannel: ch}
}

func (s *jobService) QueueBroadcastNotification(message string) (domain.Job, error) {
	job := domain.Job{
		Message:  message,
		Status:   "queued",
		QueuedAt: time.Now(),
	}

	job, err := s.jobRepo.CreateJob(job)
	if err != nil {
		return domain.Job{}, err
	}

	payload, err := json.Marshal(map[string]interface{}{"job_id": job.ID, "message": message})
	if err != nil {
		return domain.Job{}, err
	}

	err = s.rabbitMQChannel.Publish(
		"",
		"broadcast_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
	if err != nil {
		return domain.Job{}, err
	}

	return job, nil
}

func (s *jobService) GetJobStatus(id int) (domain.Job, error) {
	return s.jobRepo.GetJobByID(id)
}
