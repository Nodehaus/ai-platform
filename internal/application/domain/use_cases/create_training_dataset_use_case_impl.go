package use_cases

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type CreateTrainingDatasetUseCaseImpl struct {
	trainingDatasetRepository persistence.TrainingDatasetRepository
	projectRepository         persistence.ProjectRepository
	corpusRepository          persistence.CorpusRepository
	promptRepository          persistence.PromptRepository
	trainingDatasetService    *services.TrainingDatasetService
}

func NewCreateTrainingDatasetUseCase(
	trainingDatasetRepository persistence.TrainingDatasetRepository,
	projectRepository persistence.ProjectRepository,
	corpusRepository persistence.CorpusRepository,
	promptRepository persistence.PromptRepository,
	trainingDatasetService *services.TrainingDatasetService,
) in.CreateTrainingDatasetUseCase {
	return &CreateTrainingDatasetUseCaseImpl{
		trainingDatasetRepository: trainingDatasetRepository,
		projectRepository:         projectRepository,
		corpusRepository:          corpusRepository,
		promptRepository:          promptRepository,
		trainingDatasetService:    trainingDatasetService,
	}
}

func (uc *CreateTrainingDatasetUseCaseImpl) Execute(ctx context.Context, command in.CreateTrainingDatasetCommand) (*entities.TrainingDataset, error) {
	// Verify project exists
	project, err := uc.projectRepository.GetByID(command.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	// Validate corpus name and generate prompt
	if err := uc.trainingDatasetService.ValidateCorpusName(command.CorpusName); err != nil {
		return nil, err
	}
	if err := uc.trainingDatasetService.ValidateGeneratePrompt(command.GeneratePrompt); err != nil {
		return nil, err
	}

	// Verify corpus exists (assuming it exists as stated in requirements)
	corpus, err := uc.corpusRepository.GetByName(ctx, command.CorpusName)
	if err != nil {
		return nil, err
	}
	if corpus == nil {
		return nil, errors.New("corpus not found")
	}

	// Create prompt entity
	prompt := &entities.Prompt{
		ID:      uuid.New(),
		Version: 1,
		Text:    command.GeneratePrompt,
	}
	err = uc.promptRepository.Create(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Get next version number for this project
	nextVersion, err := uc.trainingDatasetService.GetNextVersion(command.ProjectID, func(projectID uuid.UUID) (*entities.TrainingDataset, error) {
		return uc.trainingDatasetRepository.GetLatestByProjectID(ctx, projectID)
	})
	if err != nil {
		return nil, err
	}

	// Create training dataset
	trainingDataset, err := uc.trainingDatasetService.CreateTrainingDataset(
		command.ProjectID,
		corpus.ID,
		prompt.ID,
		command.InputField,
		command.OutputField,
		command.LanguageISO,
		command.FieldNames,
	)
	if err != nil {
		return nil, err
	}

	// Set the correct version
	trainingDataset.Version = nextVersion

	// Save to repository
	err = uc.trainingDatasetRepository.Create(ctx, trainingDataset)
	if err != nil {
		return nil, err
	}

	return trainingDataset, nil
}