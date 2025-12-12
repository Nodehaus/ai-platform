package in

import (
	"context"

	"ai-platform/internal/application/port/out/clients"
)

type PublicChatCompletionResult struct {
	Response string
}

type PublicChatCompletionUseCase interface {
	GenerateChatCompletion(ctx context.Context, command PublicChatCompletionCommand) (*PublicChatCompletionResult, error)
	GenerateChatCompletionStream(ctx context.Context, command PublicChatCompletionCommand) (<-chan clients.StreamChunk, error)
}
