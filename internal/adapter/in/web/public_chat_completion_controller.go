package web

import (
	"encoding/json"
	"fmt"
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
			"error": fmt.Sprintf("Failed to generate chat completion: %v", err),
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

	// Convert messages to command ChatMessage
	messages := make([]in.ChatMessage, len(request.Messages))
	for i, msg := range request.Messages {
		content := ""
		if str, ok := msg.Content.(string); ok {
			content = str
		} else if arr, ok := msg.Content.([]interface{}); ok {
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					if typ, exists := m["type"]; exists && typ == "text" {
						if txt, exists := m["text"]; exists {
							if s, ok := txt.(string); ok {
								content += s
							}
						}
					}
				}
			}
		}
		messages[i] = in.ChatMessage{
			Role:    msg.Role,
			Content: content,
		}
	}

	// Set defaults for optional parameters
	temperature := 0.7
	if request.Temperature != nil {
		temperature = *request.Temperature
	}

	topP := 1.0
	if request.TopP != nil {
		topP = *request.TopP
	}

	stream := false
	if request.Stream != nil {
		stream = *request.Stream
	}

	command := in.PublicChatCompletionCommand{
		DeploymentID: deploymentID.(uuid.UUID),
		FinetuneID:   finetuneID,
		ModelName:    request.Model,
		Messages:     messages,
		MaxTokens:    request.MaxTokens,
		Temperature:  temperature,
		TopP:         topP,
		Stream:       stream,
	}

	// Handle streaming response
	if stream {
		c.handleStreamingResponse(ctx, command, deploymentID.(uuid.UUID), request.Model)
		return
	}

	// Handle non-streaming response
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

func (c *PublicChatCompletionController) handleStreamingResponse(ctx *gin.Context, command in.PublicChatCompletionCommand, deploymentID uuid.UUID, model string) {
	// Set SSE headers
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	// Get the streaming channel from use case
	streamChan, err := c.PublicChatCompletionUseCase.GenerateChatCompletionStream(ctx.Request.Context(), command)
	if err != nil {
		// Write error as SSE
		errorData := map[string]interface{}{
			"error": map[string]string{
				"message": err.Error(),
				"type":    "internal_error",
			},
		}
		errorJSON, _ := json.Marshal(errorData)
		fmt.Fprintf(ctx.Writer, "data: %s\n\n", errorJSON)
		ctx.Writer.(http.Flusher).Flush()
		return
	}

	// Stream the chunks
	for chunk := range streamChan {
		if chunk.Error != nil {
			// Write error and stop
			errorData := map[string]interface{}{
				"error": map[string]string{
					"message": chunk.Error.Error(),
					"type":    "internal_error",
				},
			}
			errorJSON, _ := json.Marshal(errorData)
			fmt.Fprintf(ctx.Writer, "data: %s\n\n", errorJSON)
			ctx.Writer.(http.Flusher).Flush()
			return
		}

		// Create OpenAI-compatible streaming chunk
		response := map[string]interface{}{
			"id":      deploymentID.String(),
			"object":  "chat.completion.chunk",
			"created": 0,
			"model":   model,
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"delta": map[string]interface{}{},
				},
			},
		}

		// Add content if present
		if chunk.Content != "" {
			response["choices"].([]map[string]interface{})[0]["delta"].(map[string]interface{})["content"] = chunk.Content
		}

		// Add finish_reason if present
		if chunk.FinishReason != nil {
			response["choices"].([]map[string]interface{})[0]["finish_reason"] = *chunk.FinishReason
		} else {
			response["choices"].([]map[string]interface{})[0]["finish_reason"] = nil
		}

		// Marshal to JSON and write
		chunkJSON, err := json.Marshal(response)
		if err != nil {
			continue
		}

		fmt.Fprintf(ctx.Writer, "data: %s\n\n", chunkJSON)
		ctx.Writer.(http.Flusher).Flush()
	}

	// Write [DONE] message
	fmt.Fprintf(ctx.Writer, "data: [DONE]\n\n")
	ctx.Writer.(http.Flusher).Flush()
}
