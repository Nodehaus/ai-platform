package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type CreateFinetuneController struct {
	CreateFinetuneUseCase in.CreateFinetuneUseCase
}

func (c *CreateFinetuneController) CreateFinetune(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	projectIDParam := ctx.Param("project_id")
	projectID, err := uuid.Parse(projectIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID format",
		})
		return
	}

	var request CreateFinetuneRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	trainingDatasetID, err := request.GetTrainingDatasetID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid training dataset ID format",
		})
		return
	}

	command := in.CreateFinetuneCommand{
		UserID:                           userID,
		ProjectID:                        projectID,
		BaseModelName:                    request.BaseModelName,
		TrainingDatasetID:                trainingDatasetID,
		TrainingDatasetNumberExamples:    request.TrainingDatasetNumberExamples,
		TrainingDatasetSelectRandom:      request.TrainingDatasetSelectRandom,
	}

	result, err := c.CreateFinetuneUseCase.Execute(ctx.Request.Context(), command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := ToCreateFinetuneResponse(result)
	ctx.JSON(http.StatusCreated, response)
}