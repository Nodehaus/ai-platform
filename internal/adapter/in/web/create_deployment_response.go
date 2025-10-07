package web

import (
	"ai-platform/internal/application/domain/entities"

	"github.com/google/uuid"
)

type CreateDeploymentResponse struct {
	ID        uuid.UUID `json:"id"`
	ModelName string    `json:"model_name"`
	APIKey    string    `json:"api_key"`
}

func NewCreateDeploymentResponse(deployment *entities.Deployment) *CreateDeploymentResponse {
	return &CreateDeploymentResponse{
		ID:        deployment.ID,
		ModelName: deployment.ModelName,
		APIKey:    deployment.APIKey,
	}
}
