package use_cases

import (
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
	"errors"
)

type GetDeploymentUseCaseImpl struct {
	DeploymentRepository persistence.DeploymentRepository
	DeploymentService    *services.DeploymentService
}

func (uc *GetDeploymentUseCaseImpl) GetDeployment(command in.GetDeploymentCommand) (*in.GetDeploymentResult, error) {
	// Validate project access
	err := uc.DeploymentService.ValidateProjectAccess(command.ProjectID, command.OwnerID)
	if err != nil {
		return nil, err
	}

	// Get deployment
	deployment, err := uc.DeploymentRepository.GetByID(command.DeploymentID)
	if err != nil {
		return nil, err
	}

	if deployment == nil {
		return nil, errors.New("deployment not found")
	}

	// Verify deployment belongs to project
	if deployment.ProjectID != command.ProjectID {
		return nil, errors.New("deployment does not belong to this project")
	}

	return &in.GetDeploymentResult{
		Deployment: deployment,
	}, nil
}
