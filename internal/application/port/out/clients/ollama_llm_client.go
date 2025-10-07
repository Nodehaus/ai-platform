package clients

import (
	"context"
)

type OllamaLLMClient interface {
	GenerateCompletion(ctx context.Context, finetuneID string, prompt string, model string, maxTokens int, temperature float64, topP float64) (string, error)
	GenerateChatCompletion(ctx context.Context, finetuneID string, messages []string, model string, maxTokens int, temperature float64, topP float64) (string, error)
}
