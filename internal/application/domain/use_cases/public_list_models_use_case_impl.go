package use_cases

import (
	"context"

	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type PublicListModelsUseCaseImpl struct {
	DeploymentRepository persistence.DeploymentRepository
}

func NewPublicListModelsUseCaseImpl(deploymentRepo persistence.DeploymentRepository) *PublicListModelsUseCaseImpl {
	return &PublicListModelsUseCaseImpl{
		DeploymentRepository: deploymentRepo,
	}
}

func (uc *PublicListModelsUseCaseImpl) ListModels(ctx context.Context, command in.PublicListModelsCommand) (*in.PublicListModelsResult, error) {
	// Get all deployments for the project
	deployments, err := uc.DeploymentRepository.GetByProjectID(command.ProjectID)
	if err != nil {
		return nil, err
	}

	// Convert deployments to model info format (OpenAI compatible)
	models := make([]in.ModelInfo, 0, len(deployments))
	for _, deployment := range deployments {
		models = append(models, in.ModelInfo{
			ID:      deployment.ModelName,
			Object:  "model",
			Created: deployment.CreatedAt.Unix(),
			OwnedBy: "organization",
		})
	}

	return &in.PublicListModelsResult{
		Object: "list",
		Data:   models,
	}, nil
}
