package persistence

import (
	"ai-platform/internal/application/domain/entities"
)

type DeploymentLogsRepository interface {
	Create(log *entities.DeploymentLogs) error
}
