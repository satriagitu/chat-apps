package controller

import (
	"chat-apps/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	FileService service.FileService
}

func NewFileController(fs service.FileService) *FileController {
	return &FileController{FileService: fs}
}

func (fc *FileController) UploadFile(c *gin.Context) {
	var req struct {
		UserID int    `json:"user_id"`
		File   string `json:"file"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := fc.FileService.UploadFile(req.UserID, req.File)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          file.ID,
		"user_id":     file.UserID,
		"file_url":    file.FileURL,
		"uploaded_at": file.UploadedAt,
	})
}

func (fc *FileController) GetFileByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	file, err := fc.FileService.GetFileByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          file.ID,
		"user_id":     file.UserID,
		"file_url":    file.FileURL,
		"uploaded_at": file.UploadedAt,
	})
}
