package service

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user domain.User) (domain.User, error)
	GetUserByID(id int) (domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{userRepo: ur}
}

func (s *userService) CreateUser(user domain.User) (domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.CreateUser(user)
}

func (s *userService) GetUserByID(id int) (domain.User, error) {
	return s.userRepo.GetUserByID(id)
}
