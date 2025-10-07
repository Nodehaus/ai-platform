package persistence

import (
	"ai-platform/internal/application/domain/entities"
	"time"

	"github.com/google/uuid"
)

type DeploymentRepositoryModel struct {
	ID         uuid.UUID  `db:"id"`
	ModelName  string     `db:"model_name"`
	APIKey     string     `db:"api_key"`
	ProjectID  uuid.UUID  `db:"project_id"`
	FinetuneID *uuid.UUID `db:"finetune_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

func (m *DeploymentRepositoryModel) ToEntity() *entities.Deployment {
	return &entities.Deployment{
		ID:         m.ID,
		ModelName:  m.ModelName,
		APIKey:     m.APIKey,
		ProjectID:  m.ProjectID,
		FinetuneID: m.FinetuneID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func DeploymentFromEntity(deployment *entities.Deployment) *DeploymentRepositoryModel {
	return &DeploymentRepositoryModel{
		ID:         deployment.ID,
		ModelName:  deployment.ModelName,
		APIKey:     deployment.APIKey,
		ProjectID:  deployment.ProjectID,
		FinetuneID: deployment.FinetuneID,
		CreatedAt:  deployment.CreatedAt,
		UpdatedAt:  deployment.UpdatedAt,
	}
}
