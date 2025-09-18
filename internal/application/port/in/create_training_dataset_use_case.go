package in

import (
	"context"

	"ai-platform/internal/application/domain/entities"
)

type CreateTrainingDatasetUseCase interface {
	Execute(ctx context.Context, command CreateTrainingDatasetCommand) (*entities.TrainingDataset, error)
}