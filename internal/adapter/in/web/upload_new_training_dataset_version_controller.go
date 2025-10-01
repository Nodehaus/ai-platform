package web

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type UploadNewTrainingDatasetVersionController struct {
	UploadNewTrainingDatasetVersionUseCase in.UploadNewTrainingDatasetVersionUseCase
}

func (c *UploadNewTrainingDatasetVersionController) UploadNewTrainingDatasetVersion(ctx *gin.Context) {
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

	command := in.UploadNewTrainingDatasetVersionCommand{
		ProjectID: projectID,
		OwnerID:   userID,
		CsvData:   csvData,
	}

	result, err := c.UploadNewTrainingDatasetVersionUseCase.UploadNewTrainingDatasetVersion(command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"training_dataset_id": result.TrainingDatasetID,
		"version":             result.Version,
		"total_items":         result.TotalItems,
		"message":             "New training dataset version created successfully",
	})
}
