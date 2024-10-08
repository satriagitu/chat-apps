package controller

import (
	"chat-apps/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	ArticleService service.ArticleService
}

func NewArticleController(as service.ArticleService) *ArticleController {
	return &ArticleController{
		ArticleService: as,
	}
}

func (ac *ArticleController) GetArticleList(c *gin.Context) {
	search := c.Query("search")
	artikel, err := ac.ArticleService.GetArticleList(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, artikel)
}
