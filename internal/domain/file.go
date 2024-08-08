package domain

import "time"

type File struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	UserID     int       `json:"user_id"`
	FileURL    string    `json:"file_url"`
	UploadedAt time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
}
