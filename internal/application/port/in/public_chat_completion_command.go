package in

import (
	"github.com/google/uuid"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PublicChatCompletionCommand struct {
	DeploymentID uuid.UUID
	FinetuneID   *uuid.UUID
	ModelName    string
	Messages     []ChatMessage
	MaxTokens    int
	Temperature  float64
	TopP         float64
}
