package clients

import (
	"context"
)

type OllamaLLMClient interface {
	GenerateCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64, topP float64) (string, error)
}
