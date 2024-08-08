package controller

import (
	"net/http"
	"strconv"

	"chat-apps/internal/domain"
	"chat-apps/internal/service"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	jobService service.JobService
}

func NewJobController(js service.JobService) *JobController {
	return &JobController{jobService: js}
}

func (c *JobController) QueueBroadcastNotification(ctx *gin.Context) {
	var request domain.BroadcastRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := c.jobService.QueueBroadcastNotification(request.Message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, job)
}

func (c *JobController) GetJobStatus(ctx *gin.Context) {
	jobIDStr := ctx.Param("id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := c.jobService.GetJobStatus(jobID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, job)
}
