package service

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"
	"errors"
)

type FileService interface {
	UploadFile(userID int, fileURL string) (domain.File, error)
	GetFileByID(id int) (domain.File, error)
}

type fileService struct {
	fileRepo repository.FileRepository
	userRepo repository.UserRepository
}

func NewFileService(fr repository.FileRepository, ur repository.UserRepository) FileService {
	return &fileService{fileRepo: fr, userRepo: ur}
}

func (s *fileService) UploadFile(userID int, fileURL string) (domain.File, error) {
	exists, err := s.userRepo.ExistsByID(userID)
	if err != nil {
		return domain.File{}, err
	}
	if !exists {
		return domain.File{}, errors.New("user not found")
	}

	file := domain.File{
		UserID:  userID,
		FileURL: fileURL,
	}
	return s.fileRepo.UploadFile(file)
}

func (s *fileService) GetFileByID(id int) (domain.File, error) {
	return s.fileRepo.GetFileByID(id)
}
