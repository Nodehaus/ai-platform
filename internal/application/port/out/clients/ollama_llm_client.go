package clients

import (
	"context"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaLLMClientResult struct {
	Response      string
	TokensIn      int
	TokensOut     int
	DelayTime     int
	ExecutionTime int
}

type StreamChunk struct {
	Content      string
	FinishReason *string
	Error        error
}

type OllamaLLMClient interface {
	GenerateCompletion(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens *int, temperature float64, topP float64) (*OllamaLLMClientResult, error)
	GenerateCompletionStream(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens *int, temperature float64, topP float64) (<-chan StreamChunk, error)
	GenerateChatCompletion(ctx context.Context, finetuneID *string, messages []ChatMessage, model string, maxTokens *int, temperature float64, topP float64) (*OllamaLLMClientResult, error)
	GenerateChatCompletionStream(ctx context.Context, finetuneID *string, messages []ChatMessage, model string, maxTokens *int, temperature float64, topP float64) (<-chan StreamChunk, error)
}
