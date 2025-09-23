package in

import "github.com/google/uuid"

type GetTrainingDatasetCommand struct {
	ProjectID        uuid.UUID
	TrainingDatasetID uuid.UUID
	OwnerID          uuid.UUID
}