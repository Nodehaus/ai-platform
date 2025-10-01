package in

import "github.com/google/uuid"

type UploadTrainingDatasetCommand struct {
	ProjectID         uuid.UUID
	TrainingDatasetID uuid.UUID
	OwnerID           uuid.UUID
	CsvData           []byte
}
