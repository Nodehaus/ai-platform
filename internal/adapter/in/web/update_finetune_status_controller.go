package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type UpdateFinetuneStatusController struct {
	UpdateFinetuneStatusUseCase in.UpdateFinetuneStatusUseCase
}

func (c *UpdateFinetuneStatusController) UpdateStatus(ctx *gin.Context) {
	// Extract finetune ID from URL parameter
	finetuneIDStr := ctx.Param("finetune_id")
	finetuneID, err := uuid.Parse(finetuneIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid finetune ID format",
		})
		return
	}

	// Parse request body
	var request UpdateFinetuneStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Create command
	command := in.UpdateFinetuneStatusCommand{
		FinetuneID: finetuneID,
		Status:     request.Status,
	}

	// Execute use case
	err = c.UpdateFinetuneStatusUseCase.Execute(ctx.Request.Context(), command)
	if err != nil {
		// Log the full error message to console
		fmt.Printf("Full error: %v\n", err)

		if err.Error() == "finetune not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Finetune not found",
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
			"error": "Failed to update finetune status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Finetune status updated successfully",
	})
}