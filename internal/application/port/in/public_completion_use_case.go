package in

import (
	"context"

	"ai-platform/internal/application/port/out/clients"
)

type PublicCompletionResult struct {
	Response string
}

type PublicCompletionUseCase interface {
	GenerateCompletion(ctx context.Context, command PublicCompletionCommand) (*PublicCompletionResult, error)
	GenerateCompletionStream(ctx context.Context, command PublicCompletionCommand) (<-chan clients.StreamChunk, error)
}
