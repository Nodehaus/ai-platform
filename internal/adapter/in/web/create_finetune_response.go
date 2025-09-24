package web

import (
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type CreateFinetuneResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
}

func ToCreateFinetuneResponse(f *entities.Finetune) *CreateFinetuneResponse {
	return &CreateFinetuneResponse{
		ID:        f.ID,
		ProjectID: f.ProjectID,
	}
}