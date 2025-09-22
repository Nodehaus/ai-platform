package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type UpdateTrainingDatasetStatusController struct {
	UpdateTrainingDatasetStatusUseCase in.UpdateTrainingDatasetStatusUseCase
}

func (c *UpdateTrainingDatasetStatusController) UpdateStatus(ctx *gin.Context) {
	// Extract training dataset ID from URL parameter
	trainingDatasetIDStr := ctx.Param("training_dataset_id")
	trainingDatasetID, err := uuid.Parse(trainingDatasetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid training dataset ID format",
		})
		return
	}

	// Parse request body
	var request UpdateTrainingDatasetStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Create command
	command := in.UpdateTrainingDatasetStatusCommand{
		TrainingDatasetID: trainingDatasetID,
		Status:            request.Status,
	}

	// Execute use case
	err = c.UpdateTrainingDatasetStatusUseCase.Execute(ctx.Request.Context(), command)
	if err != nil {
		if err.Error() == "training dataset not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Training dataset not found",
			})
			return
		}
		if err.Error()[:22] == "invalid status transition" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update training dataset status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Training dataset status updated successfully",
	})
}