package web

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type DownloadModelController struct {
	DownloadModelUseCase in.DownloadModelUseCase
}

func (controller *DownloadModelController) DownloadModel(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	finetuneIDStr := c.Param("finetune_id")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid project_id format",
		})
		return
	}

	finetuneID, err := uuid.Parse(finetuneIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid finetune_id format",
		})
		return
	}

	command := in.DownloadModelCommand{
		ProjectID:  projectID,
		FinetuneID: finetuneID,
	}

	reader, contentLength, filename, err := controller.DownloadModelUseCase.DownloadModel(c.Request.Context(), command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	defer reader.Close()

	// Set headers for file download
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))

	// Stream the file to the client
	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		// Cannot send JSON error after we've started streaming
		// Log the error but don't return anything to client
		return
	}
}