package web

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type UploadTrainingDatasetController struct {
	UploadTrainingDatasetUseCase in.UploadTrainingDatasetUseCase
}

func (c *UploadTrainingDatasetController) UploadTrainingDataset(ctx *gin.Context) {
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

	// Get the uploaded file
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded or invalid file",
		})
		return
	}
	defer file.Close()

	// Validate file is CSV
	if header.Header.Get("Content-Type") != "text/csv" && !isCsvFilename(header.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "File must be a CSV file",
		})
		return
	}

	// Read file content
	csvData, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read file",
		})
		return
	}

	command := in.UploadTrainingDatasetCommand{
		ProjectID:         projectID,
		TrainingDatasetID: trainingDatasetID,
		OwnerID:           userID,
		CsvData:           csvData,
	}

	result, err := c.UploadTrainingDatasetUseCase.UploadTrainingDataset(command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items_added": result.ItemsAdded,
		"total_items": result.TotalItems,
		"message":     "Training data uploaded successfully",
	})
}

func isCsvFilename(filename string) bool {
	return len(filename) > 4 && filename[len(filename)-4:] == ".csv"
}
