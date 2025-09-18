package in

import (
	"ai-platform/internal/application/domain/entities"
	"github.com/google/uuid"
)

type ProjectWithTrainingDataset struct {
	Project           entities.Project
	TrainingDatasetID *uuid.UUID
}

type ListProjectsResult struct {
	Projects []ProjectWithTrainingDataset
}

type ListProjectsUseCase interface {
	ListProjects(command ListProjectsCommand) (*ListProjectsResult, error)
}