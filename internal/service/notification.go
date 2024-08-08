package service

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"
	"errors"
	"time"
)

type NotificationService interface {
	SendNotification(userID int, message string) (domain.Notification, error)
	GetNotificationsByUserID(userID int) ([]domain.Notification, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

func NewNotificationService(nr repository.NotificationRepository, ur repository.UserRepository) NotificationService {
	return &notificationService{notificationRepo: nr, userRepo: ur}
}

func (s *notificationService) SendNotification(userID int, message string) (domain.Notification, error) {
	exists, err := s.userRepo.ExistsByID(userID)
	if err != nil {
		return domain.Notification{}, err
	}
	if !exists {
		return domain.Notification{}, errors.New("user not found")
	}

	notification := domain.Notification{
		UserID:  userID,
		Message: message,
		SentAt:  time.Now(),
	}
	return s.notificationRepo.CreateNotification(notification)
}

func (s *notificationService) GetNotificationsByUserID(userID int) ([]domain.Notification, error) {
	return s.notificationRepo.GetNotificationsByUserID(userID)
}
