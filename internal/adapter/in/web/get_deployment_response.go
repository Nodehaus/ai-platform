package web

import (
	"ai-platform/internal/application/domain/entities"
	"time"

	"github.com/google/uuid"
)

type GetDeploymentResponse struct {
	ID         uuid.UUID  `json:"id"`
	ModelName  string     `json:"model_name"`
	APIKey     string     `json:"api_key"`
	ProjectID  uuid.UUID  `json:"project_id"`
	FinetuneID *uuid.UUID `json:"finetune_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func NewGetDeploymentResponse(deployment *entities.Deployment) *GetDeploymentResponse {
	return &GetDeploymentResponse{
		ID:         deployment.ID,
		ModelName:  deployment.ModelName,
		APIKey:     deployment.APIKey,
		ProjectID:  deployment.ProjectID,
		FinetuneID: deployment.FinetuneID,
		CreatedAt:  deployment.CreatedAt,
		UpdatedAt:  deployment.UpdatedAt,
	}
}
