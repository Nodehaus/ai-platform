package in

import "context"

type FinetuneCompletionUseCase interface {
	GenerateCompletion(ctx context.Context, command FinetuneCompletionCommand) (*FinetuneCompletionResult, error)
}

type FinetuneCompletionResult struct {
	Completion string
}
