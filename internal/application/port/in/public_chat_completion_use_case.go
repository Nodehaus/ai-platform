package in

import "context"

type PublicChatCompletionResult struct {
	Response string
}

type StreamCallback func(chunk string, metadata *StreamMetadata)

type StreamMetadata struct {
	TokensIn      int
	TokensOut     int
	DelayTime     int
	ExecutionTime int
}

type PublicChatCompletionUseCase interface {
	GenerateChatCompletion(ctx context.Context, command PublicChatCompletionCommand) (*PublicChatCompletionResult, error)
	GenerateChatCompletionStream(ctx context.Context, command PublicChatCompletionCommand, callback StreamCallback) (*StreamMetadata, error)
}
