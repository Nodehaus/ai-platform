package persistence

import (
	"context"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type FinetuneRepository interface {
	Create(ctx context.Context, finetune *entities.Finetune) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Finetune, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Finetune, error)
	GetLatestByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.Finetune, error)
	Update(ctx context.Context, finetune *entities.Finetune) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.FinetuneStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetNextVersion(ctx context.Context, projectID uuid.UUID) (int, error)
}