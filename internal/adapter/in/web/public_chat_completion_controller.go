package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type PublicChatCompletionController struct {
	PublicChatCompletionUseCase in.PublicChatCompletionUseCase
}

func (c *PublicChatCompletionController) GenerateChatCompletion(ctx *gin.Context) {
	var request PublicChatCompletionRequest
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

	// Convert messages to string array (simplified for Ollama)
	messages := make([]string, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = msg.Content
	}

	// Set defaults for optional parameters
	maxTokens := 100
	if request.MaxTokens != nil {
		maxTokens = *request.MaxTokens
	}

	temperature := 0.7
	if request.Temperature != nil {
		temperature = *request.Temperature
	}

	topP := 1.0
	if request.TopP != nil {
		topP = *request.TopP
	}

	command := in.PublicChatCompletionCommand{
		DeploymentID: deploymentID.(uuid.UUID),
		FinetuneID:   finetuneID,
		ModelName:    request.Model,
		Messages:     messages,
		MaxTokens:    maxTokens,
		Temperature:  temperature,
		TopP:         topP,
	}

	result, err := c.PublicChatCompletionUseCase.GenerateChatCompletion(ctx.Request.Context(), command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return OpenAI-compatible response format
	ctx.JSON(http.StatusOK, gin.H{
		"id":      deploymentID,
		"object":  "chat.completion",
		"created": 0,
		"model":   request.Model,
		"choices": []gin.H{
			{
				"index": 0,
				"message": gin.H{
					"role":    "assistant",
					"content": result.Response,
				},
				"finish_reason": "stop",
			},
		},
	})
}
