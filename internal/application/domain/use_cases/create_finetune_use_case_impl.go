package use_cases

import (
	"context"
	"errors"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
	"ai-platform/internal/application/port/out/persistence"
)

type CreateFinetuneUseCaseImpl struct {
	FinetuneRepository        persistence.FinetuneRepository
	ProjectRepository         persistence.ProjectRepository
	TrainingDatasetRepository persistence.TrainingDatasetRepository
	CorpusRepository          persistence.CorpusRepository
	FinetuneService           *services.FinetuneService
	TrainingDatasetService    *services.TrainingDatasetService
	FinetuneJobClient         clients.FinetuneJobClient
	RunpodClient              clients.RunpodClient
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

	// Select subset of training data
	selectedData := uc.TrainingDatasetService.SelectTrainingDataSubset(
		trainingDataset.Data,
		command.TrainingDatasetNumberExamples,
		command.TrainingDatasetSelectRandom,
	)

	// Convert training data to finetune job format
	jobData := uc.TrainingDatasetService.ConvertToFinetuneJobData(
		selectedData,
		trainingDataset.FieldNames,
		trainingDataset.InputField,
	)

	// Create finetune job
	finetuneJob := entities.FinetuneJob{
		FinetuneID:        finetune.ID.String(),
		TrainingDatasetID: trainingDataset.ID.String(),
		InputField:        trainingDataset.InputField,
		OutputField:       trainingDataset.OutputField,
		UserID:            command.UserID.String(),
		TrainingData:      jobData,
	}

	// Submit job to S3
	s3Key, err := uc.FinetuneJobClient.SubmitJob(ctx, finetuneJob)
	if err != nil {
		return nil, err
	}

	// Get corpus information for documents S3 path
	corpus, err := uc.CorpusRepository.GetByID(ctx, trainingDataset.CorpusID)
	if err != nil {
		return nil, err
	}
	if corpus == nil {
		return nil, errors.New("corpus not found")
	}

	// Start finetune job on Runpod
	err = uc.RunpodClient.StartFinetuneJob(ctx, s3Key, corpus.S3Path, command.BaseModelName, modelName, finetune.ID.String())
	if err != nil {
		return nil, err
	}

	return finetune, nil
}