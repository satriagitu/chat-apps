package repository

import (
	"chat-apps/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user domain.User) (domain.User, error)
	GetUserByID(id int) (domain.User, error)
	ExistsByID(id int) (bool, error)
	GetAllUsers() ([]domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user domain.User) (domain.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(id int) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) ExistsByID(id int) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
