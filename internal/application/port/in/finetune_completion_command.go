package in

import "github.com/google/uuid"

type FinetuneCompletionCommand struct {
	ProjectID   uuid.UUID
	FinetuneID  uuid.UUID
	OwnerID     uuid.UUID
	Prompt      string
	MaxTokens   int
	Temperature float64
	TopP        float64
}
