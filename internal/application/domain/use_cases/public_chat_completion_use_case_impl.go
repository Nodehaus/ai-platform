package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
)

type PublicChatCompletionUseCaseImpl struct {
	OllamaLLMClient clients.OllamaLLMClient
}

func (uc *PublicChatCompletionUseCaseImpl) GenerateChatCompletion(ctx context.Context, command in.PublicChatCompletionCommand) (*in.PublicChatCompletionResult, error) {
	// If no finetune_id, we cannot generate completions
	if command.FinetuneID == nil {
		return nil, fmt.Errorf("deployment does not have a finetune model")
	}

	// Call OllamaLLMClient
	response, err := uc.OllamaLLMClient.GenerateChatCompletion(
		ctx,
		command.FinetuneID.String(),
		command.Messages,
		command.ModelName,
		command.MaxTokens,
		command.Temperature,
		command.TopP,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chat completion: %w", err)
	}

	return &in.PublicChatCompletionResult{
		Response: response,
	}, nil
}
