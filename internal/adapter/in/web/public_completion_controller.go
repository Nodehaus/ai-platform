package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type PublicCompletionController struct {
	PublicCompletionUseCase in.PublicCompletionUseCase
}

func (c *PublicCompletionController) GenerateCompletion(ctx *gin.Context) {
	var request PublicCompletionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Extract deployment information from context (set by middleware)
	deploymentID, exists := ctx.Get("deployment_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Deployment ID not found in context",
		})
		return
	}

	deploymentModelName, exists := ctx.Get("model_name")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Model name not found in context",
		})
		return
	}

	// Validate that the requested model matches the deployment's model
	if request.Model != deploymentModelName.(string) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Model name in request does not match the deployment",
		})
		return
	}

	// Get finetune_id if it exists
	var finetuneID *uuid.UUID
	if finetuneIDValue, exists := ctx.Get("finetune_id"); exists {
		finetuneUUID := finetuneIDValue.(uuid.UUID)
		finetuneID = &finetuneUUID
	}

	// Set defaults for optional parameters
	maxTokens := 100
	if request.MaxTokens != nil {
		maxTokens = *request.MaxTokens
	}

	temperature := 0.5
	if request.Temperature != nil {
		temperature = *request.Temperature
	}

	topP := 1.0
	if request.TopP != nil {
		topP = *request.TopP
	}

	command := in.PublicCompletionCommand{
		DeploymentID: deploymentID.(uuid.UUID),
		FinetuneID:   finetuneID,
		ModelName:    request.Model,
		Prompt:       request.Prompt,
		MaxTokens:    maxTokens,
		Temperature:  temperature,
		TopP:         topP,
	}

	result, err := c.PublicCompletionUseCase.GenerateCompletion(ctx.Request.Context(), command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return OpenAI-compatible response format
	ctx.JSON(http.StatusOK, gin.H{
		"id":      deploymentID,
		"object":  "text_completion",
		"created": 0,
		"model":   request.Model,
		"choices": []gin.H{
			{
				"text":          result.Response,
				"index":         0,
				"finish_reason": "stop",
			},
		},
	})
}
