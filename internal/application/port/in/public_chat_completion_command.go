package in

import "github.com/google/uuid"

type PublicChatCompletionCommand struct {
	DeploymentID uuid.UUID
	FinetuneID   *uuid.UUID
	ModelName    string
	Messages     []string
	MaxTokens    int
	Temperature  float64
	TopP         float64
}
