package in

import "github.com/google/uuid"

type PublicCompletionCommand struct {
	DeploymentID uuid.UUID
	FinetuneID   *uuid.UUID
	ModelName    string
	Prompt       string
	MaxTokens    int
	Temperature  float64
	TopP         float64
	Stream       bool
}
