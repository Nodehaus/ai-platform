package web

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type DownloadTrainingDatasetController struct {
	DownloadTrainingDatasetUseCase in.DownloadTrainingDatasetUseCase
}

func (c *DownloadTrainingDatasetController) DownloadTrainingDataset(ctx *gin.Context) {
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

	command := in.DownloadTrainingDatasetCommand{
		ProjectID:         projectID,
		TrainingDatasetID: trainingDatasetID,
		OwnerID:           userID,
	}

	result, err := c.DownloadTrainingDatasetUseCase.DownloadTrainingDataset(command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to download training dataset",
		})
		return
	}

	if result == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Training dataset not found",
		})
		return
	}

	// Generate CSV
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header row with field names
	if err := writer.Write(result.FieldNames); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate CSV",
		})
		return
	}

	// Write data rows
	for _, row := range result.Data {
		if err := writer.Write(row); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate CSV",
			})
			return
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate CSV",
		})
		return
	}

	// Set headers for file download
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", result.Filename))
	ctx.Data(http.StatusOK, "text/csv", buf.Bytes())
}
