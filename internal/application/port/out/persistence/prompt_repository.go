package persistence

import (
	"context"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type PromptRepository interface {
	Create(ctx context.Context, prompt *entities.Prompt) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Prompt, error)
	Update(ctx context.Context, prompt *entities.Prompt) error
	Delete(ctx context.Context, id uuid.UUID) error
}