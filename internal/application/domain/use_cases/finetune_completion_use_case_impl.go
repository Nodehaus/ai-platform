package use_cases

import (
	"context"

	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
)

type FinetuneCompletionUseCaseImpl struct {
	finetuneCompletionService *services.FinetuneCompletionService
}

func NewFinetuneCompletionUseCaseImpl(
	finetuneCompletionService *services.FinetuneCompletionService,
) *FinetuneCompletionUseCaseImpl {
	return &FinetuneCompletionUseCaseImpl{
		finetuneCompletionService: finetuneCompletionService,
	}
}

func (uc *FinetuneCompletionUseCaseImpl) GenerateCompletion(ctx context.Context, command in.FinetuneCompletionCommand) (*in.FinetuneCompletionResult, error) {
	// Validate ownership
	if err := uc.finetuneCompletionService.ValidateOwnership(ctx, command.ProjectID, command.OwnerID); err != nil {
		return nil, err
	}

	// Get the model name from the finetune
	modelName, err := uc.finetuneCompletionService.GetFinetuneModelName(ctx, command.FinetuneID)
	if err != nil {
		return nil, err
	}

	// Generate completion
	completion, err := uc.finetuneCompletionService.GenerateCompletion(
		ctx,
		command.FinetuneID,
		modelName,
		command.Prompt,
		command.MaxTokens,
		command.Temperature,
		command.TopP,
	)
	if err != nil {
		return nil, err
	}

	return &in.FinetuneCompletionResult{
		Completion: completion,
	}, nil
}
