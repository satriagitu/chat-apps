package repository

import (
	"chat-apps/internal/domain"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(notification domain.Notification) (domain.Notification, error)
	GetNotificationsByUserID(userID int) ([]domain.Notification, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(notification domain.Notification) (domain.Notification, error) {
	if err := r.db.Create(&notification).Error; err != nil {
		return domain.Notification{}, err
	}
	return notification, nil
}

func (r *notificationRepository) GetNotificationsByUserID(userID int) ([]domain.Notification, error) {
	var notifications []domain.Notification
	if err := r.db.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
