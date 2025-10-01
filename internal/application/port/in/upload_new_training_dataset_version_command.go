package in

import "github.com/google/uuid"

type UploadNewTrainingDatasetVersionCommand struct {
	ProjectID uuid.UUID
	OwnerID   uuid.UUID
	CsvData   []byte
}
