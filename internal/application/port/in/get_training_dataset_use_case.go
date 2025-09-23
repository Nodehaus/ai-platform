package in

import (
	"ai-platform/internal/application/domain/entities"
)

type GetTrainingDatasetResult struct {
	TrainingDataset *entities.TrainingDataset
	GeneratePrompt  string
	CorpusName      string
}

type GetTrainingDatasetUseCase interface {
	GetTrainingDataset(command GetTrainingDatasetCommand) (*GetTrainingDatasetResult, error)
}