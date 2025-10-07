package in

import "context"

type PublicCompletionResult struct {
	Response string
}

type PublicCompletionUseCase interface {
	GenerateCompletion(ctx context.Context, command PublicCompletionCommand) (*PublicCompletionResult, error)
}
