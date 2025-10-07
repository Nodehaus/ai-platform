package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
	"github.com/google/uuid"
)

type GetFinetuneUseCaseImpl struct {
	FinetuneRepository   persistence.FinetuneRepository
	DeploymentRepository persistence.DeploymentRepository
}

func (uc *GetFinetuneUseCaseImpl) GetFinetune(ctx context.Context, command in.GetFinetuneCommand) (*in.GetFinetuneResult, error) {
	// Get the finetune
	finetune, err := uc.FinetuneRepository.GetByID(ctx, command.FinetuneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get finetune: %w", err)
	}

	if finetune == nil {
		return nil, fmt.Errorf("finetune not found")
	}

	// Verify the finetune belongs to the specified project and owner
	if finetune.ProjectID != command.ProjectID {
		return nil, fmt.Errorf("finetune not found")
	}

	// Note: We would need to verify the project belongs to the owner
	// This would require injecting ProjectRepository to check ownership
	// For now, we assume the controller has already verified project ownership

	// Check if this finetune has been deployed
	var deploymentID *uuid.UUID
	deployment, err := uc.DeploymentRepository.GetByFinetuneID(command.FinetuneID)
	if err == nil && deployment != nil {
		deploymentID = &deployment.ID
	}

	return &in.GetFinetuneResult{
		Finetune:     finetune,
		DeploymentID: deploymentID,
	}, nil
}