package util

import (
	"encoding/json"
	"log"
	"time"

	"chat-apps/internal/domain"
	"chat-apps/internal/repository"

	"github.com/streadway/amqp"
)

type NotificationWorker struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	jobRepo          repository.JobRepository
}

func NewNotificationWorker(nr repository.NotificationRepository, ur repository.UserRepository, jr repository.JobRepository) *NotificationWorker {
	return &NotificationWorker{
		notificationRepo: nr,
		userRepo:         ur,
		jobRepo:          jr,
	}
}

func (w *NotificationWorker) ProcessTask(delivery amqp.Delivery) {
	var payload map[string]interface{}
	if err := json.Unmarshal(delivery.Body, &payload); err != nil {
		log.Printf("failed to unmarshal payload: %v", err)
		delivery.Nack(false, false)
		return
	}

	jobIDFloat, ok := payload["job_id"].(float64)
	if !ok {
		log.Printf("job_id is missing or invalid")
		delivery.Nack(false, false)
		return
	}
	jobID := int(jobIDFloat)

	message, ok := payload["message"].(string)
	if !ok {
		log.Printf("message is missing or invalid")
		delivery.Nack(false, false)
		return
	}

	users, err := w.userRepo.GetAllUsers()
	if err != nil {
		log.Printf("failed to get users: %v", err)
		delivery.Nack(false, false)
		return
	}

	for _, user := range users {
		notification := domain.Notification{
			UserID:  user.ID,
			Message: message,
			SentAt:  time.Now(),
		}

		_, err := w.notificationRepo.CreateNotification(notification)
		if err != nil {
			log.Printf("failed to send notification to user %d: %v", user.ID, err)
		}
	}

	err = w.jobRepo.UpdateJobStatus(jobID, "completed", time.Now())
	if err != nil {
		log.Printf("failed to update job status: %v", err)
		delivery.Nack(false, false)
		return
	}

	delivery.Ack(false)
}
