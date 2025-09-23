package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type GetTrainingDatasetController struct {
	GetTrainingDatasetUseCase in.GetTrainingDatasetUseCase
}

func (c *GetTrainingDatasetController) GetTrainingDataset(ctx *gin.Context) {
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
			"error": "Invalid project ID format",
		})
		return
	}

	trainingDatasetIDStr := ctx.Param("training_dataset_id")
	trainingDatasetID, err := uuid.Parse(trainingDatasetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid training dataset ID format",
		})
		return
	}

	command := in.GetTrainingDatasetCommand{
		ProjectID:         projectID,
		TrainingDatasetID: trainingDatasetID,
		OwnerID:           userID,
	}

	result, err := c.GetTrainingDatasetUseCase.GetTrainingDataset(command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch training dataset",
		})
		return
	}

	if result == nil || result.TrainingDataset == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Training dataset not found",
		})
		return
	}

	response := ToGetTrainingDatasetResponse(result.TrainingDataset, result.GeneratePrompt, result.CorpusName)
	ctx.JSON(http.StatusOK, response)
}