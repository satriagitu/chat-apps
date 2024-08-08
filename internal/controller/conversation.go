package controller

import (
	"chat-apps/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConversationController struct {
	ConversationService service.ConversationService
}

func NewConversationController(cs service.ConversationService) *ConversationController {
	return &ConversationController{ConversationService: cs}
}

func (cc *ConversationController) CreateConversation(c *gin.Context) {
	var req struct {
		Participants []int `json:"participants"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conversation, err := cc.ConversationService.CreateConversation(req.Participants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func (cc *ConversationController) GetConversationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("conversationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	conversation, err := cc.ConversationService.GetConversationByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if conversation.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}
