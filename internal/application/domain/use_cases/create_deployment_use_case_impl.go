package use_cases

import (
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
	"context"
)

type CreateDeploymentUseCaseImpl struct {
	DeploymentRepository persistence.DeploymentRepository
	DeploymentService    *services.DeploymentService
}

func (uc *CreateDeploymentUseCaseImpl) CreateDeployment(command in.CreateDeploymentCommand) (*in.CreateDeploymentResult, error) {
	ctx := context.Background()

	// Validate model name
	err := uc.DeploymentService.ValidateModelName(command.ModelName)
	if err != nil {
		return nil, err
	}

	// Validate project access
	err = uc.DeploymentService.ValidateProjectAccess(command.ProjectID, command.OwnerID)
	if err != nil {
		return nil, err
	}

	// Validate finetune if provided
	if command.FinetuneID != nil {
		err = uc.DeploymentService.ValidateFinetuneExists(ctx, *command.FinetuneID, command.ProjectID)
		if err != nil {
			return nil, err
		}

		err = uc.DeploymentService.ValidateFinetuneNotAlreadyDeployed(*command.FinetuneID)
		if err != nil {
			return nil, err
		}
	}

	// Validate model name is unique in project
	err = uc.DeploymentService.ValidateModelNameUnique(command.ProjectID, command.ModelName)
	if err != nil {
		return nil, err
	}

	// Create deployment
	deployment := uc.DeploymentService.CreateDeployment(command.ModelName, command.ProjectID, command.FinetuneID)

	err = uc.DeploymentRepository.Create(deployment)
	if err != nil {
		return nil, err
	}

	return &in.CreateDeploymentResult{
		Deployment: deployment,
	}, nil
}
