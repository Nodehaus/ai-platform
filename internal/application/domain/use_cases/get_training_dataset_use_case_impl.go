package use_cases

import (
	"context"

	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type GetTrainingDatasetUseCaseImpl struct {
	TrainingDatasetService *services.TrainingDatasetService
	TrainingDatasetRepository persistence.TrainingDatasetRepository
	PromptRepository       persistence.PromptRepository
	CorpusRepository       persistence.CorpusRepository
}

func (uc *GetTrainingDatasetUseCaseImpl) GetTrainingDataset(command in.GetTrainingDatasetCommand) (*in.GetTrainingDatasetResult, error) {
	// Get the training dataset
	trainingDataset, err := uc.TrainingDatasetRepository.GetByID(context.Background(), command.TrainingDatasetID)
	if err != nil {
		return nil, err
	}

	if trainingDataset == nil {
		return nil, nil
	}

	// Verify the training dataset belongs to the specified project
	if trainingDataset.ProjectID != command.ProjectID {
		return nil, nil
	}

	// Get the generate prompt
	prompt, err := uc.PromptRepository.GetByID(context.Background(), trainingDataset.GeneratePromptID)
	if err != nil {
		return nil, err
	}

	generatePromptText := ""
	if prompt != nil {
		generatePromptText = prompt.Text
	}

	// Get the corpus name
	corpus, err := uc.CorpusRepository.GetByID(context.Background(), trainingDataset.CorpusID)
	if err != nil {
		return nil, err
	}

	corpusName := ""
	if corpus != nil {
		corpusName = corpus.Name
	}

	return &in.GetTrainingDatasetResult{
		TrainingDataset: trainingDataset,
		GeneratePrompt:  generatePromptText,
		CorpusName:      corpusName,
	}, nil
}