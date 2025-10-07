package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
)

type PublicCompletionUseCaseImpl struct {
	OllamaLLMClient clients.OllamaLLMClient
}

func (uc *PublicCompletionUseCaseImpl) GenerateCompletion(ctx context.Context, command in.PublicCompletionCommand) (*in.PublicCompletionResult, error) {
	// If no finetune_id, we cannot generate completions
	if command.FinetuneID == nil {
		return nil, fmt.Errorf("deployment does not have a finetune model")
	}

	// Call OllamaLLMClient
	response, err := uc.OllamaLLMClient.GenerateCompletion(
		ctx,
		command.FinetuneID.String(),
		command.Prompt,
		command.ModelName,
		command.MaxTokens,
		command.Temperature,
		command.TopP,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate completion: %w", err)
	}

	return &in.PublicCompletionResult{
		Response: response,
	}, nil
}
