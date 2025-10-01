package in

import "github.com/google/uuid"

type DownloadTrainingDatasetCommand struct {
	ProjectID         uuid.UUID
	TrainingDatasetID uuid.UUID
	OwnerID           uuid.UUID
}
