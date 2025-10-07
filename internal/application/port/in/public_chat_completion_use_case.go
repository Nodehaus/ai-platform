package in

import "context"

type PublicChatCompletionResult struct {
	Response string
}

type PublicChatCompletionUseCase interface {
	GenerateChatCompletion(ctx context.Context, command PublicChatCompletionCommand) (*PublicChatCompletionResult, error)
}
