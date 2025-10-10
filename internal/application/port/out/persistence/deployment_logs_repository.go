package persistence

import (
	"ai-platform/internal/application/domain/entities"

	"github.com/google/uuid"
)

type DeploymentLogsRepository interface {
	Create(log *entities.DeploymentLogs) error
	GetLatest(deploymentID uuid.UUID, limit int) ([]*entities.DeploymentLogs, error)
}
