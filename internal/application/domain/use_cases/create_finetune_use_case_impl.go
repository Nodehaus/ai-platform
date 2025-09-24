package use_cases

import (
	"context"
	"errors"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type CreateFinetuneUseCaseImpl struct {
	FinetuneRepository        persistence.FinetuneRepository
	ProjectRepository         persistence.ProjectRepository
	TrainingDatasetRepository persistence.TrainingDatasetRepository
	FinetuneService           *services.FinetuneService
}

func (uc *CreateFinetuneUseCaseImpl) Execute(ctx context.Context, command in.CreateFinetuneCommand) (*entities.Finetune, error) {
	// Verify project exists and user has access
	project, err := uc.ProjectRepository.GetByID(command.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}
	if project.OwnerID != command.UserID {
		return nil, errors.New("access denied")
	}

	// Verify training dataset exists and belongs to the project
	trainingDataset, err := uc.TrainingDatasetRepository.GetByID(ctx, command.TrainingDatasetID)
	if err != nil {
		return nil, err
	}
	if trainingDataset == nil {
		return nil, errors.New("training dataset not found")
	}
	if trainingDataset.ProjectID != command.ProjectID {
		return nil, errors.New("training dataset does not belong to project")
	}
	if trainingDataset.Status != entities.TrainingDatasetStatusDone {
		return nil, errors.New("training dataset must be in DONE status")
	}

	// Validate base model name
	if err := uc.FinetuneService.ValidateBaseModelName(command.BaseModelName); err != nil {
		return nil, err
	}

	// Get next version number
	version, err := uc.FinetuneRepository.GetNextVersion(ctx, command.ProjectID)
	if err != nil {
		return nil, err
	}

	// Generate model name
	modelName := uc.FinetuneService.GenerateModelName(command.BaseModelName, project.Name, version)

	// Create finetune entity
	finetune := uc.FinetuneService.CreateFinetune(
		command.ProjectID,
		command.TrainingDatasetID,
		version,
		modelName,
		command.BaseModelName,
		command.TrainingDatasetNumberExamples,
		command.TrainingDatasetSelectRandom,
	)

	// Save to repository
	err = uc.FinetuneRepository.Create(ctx, finetune)
	if err != nil {
		return nil, err
	}

	return finetune, nil
}