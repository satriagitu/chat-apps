package controller

import (
	"chat-apps/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExternalAPIController struct {
	APIEndpoint string
}

func NewExternalPostController(apiEndpoint string) *ExternalAPIController {
	return &ExternalAPIController{APIEndpoint: apiEndpoint}
}

func (ctrl *ExternalAPIController) GetExternalPosts(c *gin.Context) {
	resp, err := http.Get(fmt.Sprint(ctrl.APIEndpoint, "posts"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	data := []domain.ExternalPost{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}

	c.JSON(http.StatusOK, data)
}
