package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type PublicListModelsController struct {
	PublicListModelsUseCase in.PublicListModelsUseCase
}

func (c *PublicListModelsController) ListModels(ctx *gin.Context) {
	// Extract project_id from context (set by middleware)
	projectIDValue, exists := ctx.Get("project_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Project ID not found in context",
		})
		return
	}

	projectID, ok := projectIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid project ID format in context",
		})
		return
	}

	command := in.PublicListModelsCommand{
		ProjectID: projectID,
	}

	result, err := c.PublicListModelsUseCase.ListModels(ctx.Request.Context(), command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return OpenAI-compatible response format
	ctx.JSON(http.StatusOK, result)
}
