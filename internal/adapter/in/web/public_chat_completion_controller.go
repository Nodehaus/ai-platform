package web

import (
	"fmt"
	"net/http"
	"time"

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

	// Convert messages to command ChatMessage
	messages := make([]in.ChatMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = in.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
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

	stream := false
	if request.Stream != nil {
		stream = *request.Stream
	}

	command := in.PublicChatCompletionCommand{
		DeploymentID: deploymentID.(uuid.UUID),
		FinetuneID:   finetuneID,
		ModelName:    request.Model,
		Messages:     messages,
		MaxTokens:    maxTokens,
		Temperature:  temperature,
		TopP:         topP,
		Stream:       stream,
	}

	// Handle streaming response
	if stream {
		c.handleStreamingResponse(ctx, command)
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
		"created": time.Now().Unix(),
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

func (c *PublicChatCompletionController) handleStreamingResponse(ctx *gin.Context, command in.PublicChatCompletionCommand) {
	// Set SSE headers
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	// Get the response writer
	w := ctx.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Streaming not supported",
		})
		return
	}

	// Create a channel to receive streamed content
	contentChan := make(chan string)
	errChan := make(chan error, 1)

	// Start streaming in a goroutine
	go func() {
		callback := func(chunk string, metadata *in.StreamMetadata) {
			if chunk != "" {
				contentChan <- chunk
			}
		}

		metadata, err := c.PublicChatCompletionUseCase.GenerateChatCompletionStream(
			ctx.Request.Context(),
			command,
			callback,
		)
		if err != nil {
			errChan <- err
		}

		// Send final response structure
		if metadata != nil && len(contentChan) == 0 {
			// This will signal completion
		}

		close(contentChan)
	}()

	// Stream the response
	deploymentID := command.DeploymentID.String()
	contentBuffer := ""

	for {
		select {
		case chunk, ok := <-contentChan:
			if !ok {
				// Stream ended, send done message
				fmt.Fprintf(w, "data: [DONE]\n\n")
				flusher.Flush()
				return
			}

			contentBuffer += chunk

			// Send content chunk in OpenAI format
			fmt.Fprintf(w, "data: {\"id\":\"%s\",\"object\":\"chat.completion.chunk\",\"created\":%d,\"model\":\"%s\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"%s\"},\"finish_reason\":null}]}\n\n",
				deploymentID,
				time.Now().Unix(),
				command.ModelName,
				escapeJSON(chunk),
			)
			flusher.Flush()

		case err := <-errChan:
			if err != nil {
				fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", escapeJSON(err.Error()))
				flusher.Flush()
			}
			return

		case <-ctx.Request.Context().Done():
			return
		}
	}
}

// escapeJSON escapes special characters in a string for JSON
func escapeJSON(s string) string {
	b := make([]byte, 0, len(s))
	for _, ch := range s {
		switch ch {
		case '"':
			b = append(b, '\\', '"')
		case '\\':
			b = append(b, '\\', '\\')
		case '\b':
			b = append(b, '\\', 'b')
		case '\f':
			b = append(b, '\\', 'f')
		case '\n':
			b = append(b, '\\', 'n')
		case '\r':
			b = append(b, '\\', 'r')
		case '\t':
			b = append(b, '\\', 't')
		default:
			b = append(b, byte(ch))
		}
	}
	return string(b)
}
