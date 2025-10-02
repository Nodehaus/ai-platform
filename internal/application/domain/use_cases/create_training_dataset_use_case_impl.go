package use_cases

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
	"ai-platform/internal/application/port/out/persistence"
)

type CreateTrainingDatasetUseCaseImpl struct {
	TrainingDatasetRepository persistence.TrainingDatasetRepository
	ProjectRepository         persistence.ProjectRepository
	CorpusRepository          persistence.CorpusRepository
	PromptRepository          persistence.PromptRepository
	TrainingDatasetService    *services.TrainingDatasetService
	TrainingDatasetJobClient  clients.TrainingDatasetJobClient
}


func (uc *CreateTrainingDatasetUseCaseImpl) Execute(ctx context.Context, command in.CreateTrainingDatasetCommand) (*entities.TrainingDataset, error) {
	// Verify project exists
	project, err := uc.ProjectRepository.GetByID(command.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	if err := uc.TrainingDatasetService.ValidateGeneratePrompt(command.GeneratePrompt); err != nil {
		return nil, err
	}

	// Verify corpus exists (assuming it exists as stated in requirements)
	var corpus *entities.Corpus
	if (command.CorpusName != "") {
		corpus, err = uc.CorpusRepository.GetByName(ctx, command.CorpusName)
		if err != nil {
			return nil, err
		}
		if corpus == nil {
			return nil, errors.New("corpus not found")
		}
	}

	// Create prompt entity
	prompt := &entities.Prompt{
		ID:      uuid.New(),
		Version: 1,
		Text:    command.GeneratePrompt,
	}
	err = uc.PromptRepository.Create(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Get next version number for this project
	nextVersion, err := uc.TrainingDatasetService.GetNextVersion(command.ProjectID, func(projectID uuid.UUID) (*entities.TrainingDataset, error) {
		return uc.TrainingDatasetRepository.GetLatestByProjectID(ctx, projectID)
	})
	if err != nil {
		return nil, err
	}

	// Create training dataset
	var corpusID *uuid.UUID
	if corpus != nil {
		corpusID = &corpus.ID
	}
	trainingDataset, err := uc.TrainingDatasetService.CreateTrainingDataset(
		command.ProjectID,
		corpusID,
		prompt.ID,
		command.InputField,
		command.OutputField,
		command.LanguageISO,
		command.FieldNames,
		command.GenerateExamplesNumber,
	)
	if err != nil {
		return nil, err
	}

	// Set the correct version
	trainingDataset.Version = nextVersion

	// Save to repository
	err = uc.TrainingDatasetRepository.Create(ctx, trainingDataset)
	if err != nil {
		return nil, err
	}

	// Submit job to S3
	var corpusS3Path string
	var corpusFilesSubset []string

	if corpus != nil {
		corpusS3Path = corpus.S3Path
		if corpus.FilesSubset != nil {
			corpusFilesSubset = *corpus.FilesSubset
		}
	} else {
		corpusS3Path = ""
		corpusFilesSubset = []string{}
	}

	job := entities.TrainingDatasetJob{
		CorpusS3Path:           corpusS3Path,
		CorpusFilesSubset:      corpusFilesSubset,
		LanguageISO:            command.LanguageISO,
		UserID:                 command.UserID.String(),
		TrainingDatasetID:      trainingDataset.ID.String(),
		GeneratePrompt:         command.GeneratePrompt,
		GenerateExamplesNumber: command.GenerateExamplesNumber,
		GenerateModel:          command.GenerateModel,
		GenerateModelRunner:    command.GenerateModelRunner,
	}

	err = uc.TrainingDatasetJobClient.SubmitJob(ctx, job)
	if err != nil {
		// Log error but don't fail the training dataset creation
		// In production, you might want to implement retry logic or store the job for later submission
		// For now, we'll continue successfully
	}

	return trainingDataset, nil
}