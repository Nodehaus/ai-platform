package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type FinetuneCompletionController struct {
	FinetuneCompletionUseCase in.FinetuneCompletionUseCase
}

func (c *FinetuneCompletionController) GenerateCompletion(ctx *gin.Context) {
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

	finetuneIDStr := ctx.Param("finetune_id")
	finetuneID, err := uuid.Parse(finetuneIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid finetune ID format",
		})
		return
	}

	var request FinetuneCompletionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	command := in.FinetuneCompletionCommand{
		ProjectID:   projectID,
		FinetuneID:  finetuneID,
		OwnerID:     userID,
		Prompt:      request.Prompt,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
	}

	result, err := c.FinetuneCompletionUseCase.GenerateCompletion(ctx.Request.Context(), command)
	if err != nil {
		if err.Error() == "finetune not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Finetune not found",
			})
			return
		}
		if err.Error() == "finetune is not ready for inference" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Finetune is not ready for inference",
			})
			return
		}
		if err.Error() == "unauthorized: project does not belong to user" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		if err.Error() == "project not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate completion",
		})
		return
	}

	response := FinetuneCompletionResponse{
		Completion: result.Completion,
	}
	ctx.JSON(http.StatusOK, response)
}
