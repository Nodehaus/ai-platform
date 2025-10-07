package in

import "ai-platform/internal/application/domain/entities"

type CreateDeploymentResult struct {
	Deployment *entities.Deployment
}

type CreateDeploymentUseCase interface {
	CreateDeployment(command CreateDeploymentCommand) (*CreateDeploymentResult, error)
}
