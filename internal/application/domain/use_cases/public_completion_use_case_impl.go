package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
	"ai-platform/internal/application/port/out/persistence"
	"github.com/google/uuid"
)

type PublicCompletionUseCaseImpl struct {
	OllamaLLMClient          clients.OllamaLLMClient
	DeploymentLogsRepository persistence.DeploymentLogsRepository
}

func (uc *PublicCompletionUseCaseImpl) GenerateCompletion(ctx context.Context, command in.PublicCompletionCommand) (*in.PublicCompletionResult, error) {
	// If no finetune_id, we cannot generate completions
	if command.FinetuneID == nil {
		return nil, fmt.Errorf("deployment does not have a finetune model")
	}

	// Call OllamaLLMClient
	result, err := uc.OllamaLLMClient.GenerateCompletion(
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

	// Log the request and response
	log := &entities.DeploymentLogs{
		ID:            uuid.New(),
		DeploymentID:  command.DeploymentID,
		TokensIn:      result.TokensIn,
		TokensOut:     result.TokensOut,
		Input:         command.Prompt,
		Output:        result.Response,
		DelayTime:     result.DelayTime,
		ExecutionTime: result.ExecutionTime,
		Source:        "api",
	}

	if err := uc.DeploymentLogsRepository.Create(log); err != nil {
		return nil, fmt.Errorf("failed to log deployment request: %w", err)
	}

	return &in.PublicCompletionResult{
		Response: result.Response,
	}, nil
}
