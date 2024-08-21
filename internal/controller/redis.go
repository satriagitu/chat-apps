package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chat-apps/internal/service"
)

type CacheController struct {
	cacheService service.CacheService
}

func NewCacheController(cacheService service.CacheService) *CacheController {
	return &CacheController{cacheService: cacheService}
}

func (cc *CacheController) SetCache(c *gin.Context) {
	var request struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := cc.cacheService.CacheData(c, request.Key, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (cc *CacheController) GetCache(c *gin.Context) {
	key := c.Param("key")

	value, err := cc.cacheService.RetrieveData(c, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache"})
		return
	}

	if value == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}
