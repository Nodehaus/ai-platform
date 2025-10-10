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

type DownloadDeploymentLogsController struct {
	DownloadDeploymentLogsUseCase in.DownloadDeploymentLogsUseCase
}

func (c *DownloadDeploymentLogsController) DownloadDeploymentLogs(ctx *gin.Context) {
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

	deploymentIDStr := ctx.Param("deployment_id")
	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deployment ID format",
		})
		return
	}

	command := in.DownloadDeploymentLogsCommand{
		ProjectID:    projectID,
		DeploymentID: deploymentID,
		OwnerID:      userID,
	}

	result, err := c.DownloadDeploymentLogsUseCase.DownloadDeploymentLogs(command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to download deployment logs",
		})
		return
	}

	if result == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Deployment not found",
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
