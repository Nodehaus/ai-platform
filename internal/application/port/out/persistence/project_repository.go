package persistence

import (
	"ai-platform/internal/application/domain/entities"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(project *entities.Project) error
	GetByID(id uuid.UUID) (*entities.Project, error)
	GetByOwnerID(ownerID uuid.UUID) ([]entities.Project, error)
	ExistsByNameAndOwnerID(name string, ownerID uuid.UUID) (bool, error)
	Update(project *entities.Project) error
	Delete(id uuid.UUID) error
}