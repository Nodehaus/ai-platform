package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type CreateDeploymentController struct {
	CreateDeploymentUseCase in.CreateDeploymentUseCase
}

func (c *CreateDeploymentController) CreateDeployment(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	projectIDStr := ctx.Param("project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var request CreateDeploymentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	command := in.CreateDeploymentCommand{
		ModelName:  request.ModelName,
		ProjectID:  projectID,
		FinetuneID: request.FinetuneID,
		OwnerID:    userID,
	}

	result, err := c.CreateDeploymentUseCase.CreateDeployment(command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := NewCreateDeploymentResponse(result.Deployment)
	ctx.JSON(http.StatusCreated, response)
}
