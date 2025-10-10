package in

import "ai-platform/internal/application/domain/entities"

type GetDeploymentResult struct {
	Deployment *entities.Deployment
	Logs       []*entities.DeploymentLogs
}

type GetDeploymentUseCase interface {
	GetDeployment(command GetDeploymentCommand) (*GetDeploymentResult, error)
}
