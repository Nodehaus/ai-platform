package clients

import (
	"context"
)

type OllamaLLMClientResult struct {
	Response      string
	TokensIn      int
	TokensOut     int
	DelayTime     int
	ExecutionTime int
}

type OllamaLLMClient interface {
	GenerateCompletion(ctx context.Context, finetuneID string, prompt string, model string, maxTokens int, temperature float64, topP float64) (*OllamaLLMClientResult, error)
	GenerateChatCompletion(ctx context.Context, finetuneID string, messages []string, model string, maxTokens int, temperature float64, topP float64) (*OllamaLLMClientResult, error)
}
