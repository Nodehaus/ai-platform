package persistence

import (
	"context"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type TrainingDatasetRepository interface {
	Create(ctx context.Context, trainingDataset *entities.TrainingDataset) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.TrainingDataset, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TrainingDataset, error)
	GetLatestByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.TrainingDataset, error)
	Update(ctx context.Context, trainingDataset *entities.TrainingDataset) error
	Delete(ctx context.Context, id uuid.UUID) error
}