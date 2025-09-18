package web

import (
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type CreateTrainingDatasetResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
}

func ToCreateTrainingDatasetResponse(td *entities.TrainingDataset) *CreateTrainingDatasetResponse {
	return &CreateTrainingDatasetResponse{
		ID:        td.ID,
		ProjectID: td.ProjectID,
	}
}