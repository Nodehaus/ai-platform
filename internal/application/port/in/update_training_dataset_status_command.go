package in

import (
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type UpdateTrainingDatasetStatusCommand struct {
	TrainingDatasetID uuid.UUID                     `json:"training_dataset_id"`
	Status            entities.TrainingDatasetStatus `json:"status"`
}