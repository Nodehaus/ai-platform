package use_cases

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
	"ai-platform/internal/application/port/out/persistence"
	"github.com/google/uuid"
)

type PublicChatCompletionUseCaseImpl struct {
	OllamaLLMClient           clients.OllamaLLMClient
	DeploymentLogsRepository  persistence.DeploymentLogsRepository
}

func (uc *PublicChatCompletionUseCaseImpl) GenerateChatCompletion(ctx context.Context, command in.PublicChatCompletionCommand) (*in.PublicChatCompletionResult, error) {
	// Check if finetune_id is required (only for nodehaus models)
	if command.FinetuneID == nil && strings.HasPrefix(command.ModelName, "nodehaus") {
		return nil, fmt.Errorf("deployment does not have a finetune model")
	}

	// Prepare finetuneID string pointer for the client
	var finetuneIDStr *string
	if command.FinetuneID != nil {
		idStr := command.FinetuneID.String()
		finetuneIDStr = &idStr
	}

	// Convert command messages directly to client models
	clientMessages := make([]clients.ChatMessage, len(command.Messages))
	for i, msg := range command.Messages {
		clientMessages[i] = clients.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Call OllamaLLMClient
	result, err := uc.OllamaLLMClient.GenerateChatCompletion(
		ctx,
		finetuneIDStr,
		clientMessages,
		command.ModelName,
		command.MaxTokens,
		command.Temperature,
		command.TopP,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chat completion: %w", err)
	}

	// Convert command messages to JSON string for logging
	messagesJSON, err := json.Marshal(command.Messages)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal messages: %w", err)
	}

	// Log the request and response
	log := &entities.DeploymentLogs{
		ID:            uuid.New(),
		DeploymentID:  command.DeploymentID,
		TokensIn:      result.TokensIn,
		TokensOut:     result.TokensOut,
		Input:         string(messagesJSON),
		Output:        result.Response,
		DelayTime:     result.DelayTime,
		ExecutionTime: result.ExecutionTime,
		Source:        "api",
	}

	if err := uc.DeploymentLogsRepository.Create(log); err != nil {
		return nil, fmt.Errorf("failed to log deployment request: %w", err)
	}

	return &in.PublicChatCompletionResult{
		Response: result.Response,
	}, nil
}
