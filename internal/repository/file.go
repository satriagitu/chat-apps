package repository

import (
	"chat-apps/internal/domain"
	"time"

	"gorm.io/gorm"
)

type FileRepository interface {
	UploadFile(file domain.File) (domain.File, error)
	GetFileByID(id int) (domain.File, error)
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) UploadFile(file domain.File) (domain.File, error) {
	file.UploadedAt = time.Now()
	if err := r.db.Create(&file).Error; err != nil {
		return domain.File{}, err
	}
	return file, nil
}

func (r *fileRepository) GetFileByID(id int) (domain.File, error) {
	var file domain.File
	if err := r.db.First(&file, id).Error; err != nil {
		return domain.File{}, err
	}
	return file, nil
}
