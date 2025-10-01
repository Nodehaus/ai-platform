package in

import "github.com/google/uuid"

type UploadNewTrainingDatasetVersionResult struct {
	TrainingDatasetID uuid.UUID
	Version           int
	TotalItems        int
}

type UploadNewTrainingDatasetVersionUseCase interface {
	UploadNewTrainingDatasetVersion(command UploadNewTrainingDatasetVersionCommand) (*UploadNewTrainingDatasetVersionResult, error)
}
