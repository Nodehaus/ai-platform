package persistence

import (
	"ai-platform/internal/application/domain/entities"
	"github.com/google/uuid"
)

type DeploymentRepository interface {
	Create(deployment *entities.Deployment) error
	GetByID(id uuid.UUID) (*entities.Deployment, error)
	GetByProjectID(projectID uuid.UUID) ([]entities.Deployment, error)
	GetByFinetuneID(finetuneID uuid.UUID) (*entities.Deployment, error)
	GetByProjectIDAndModelName(projectID uuid.UUID, modelName string) (*entities.Deployment, error)
	Delete(id uuid.UUID) error
}
