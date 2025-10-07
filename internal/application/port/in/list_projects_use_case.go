package in

import (
	"ai-platform/internal/application/domain/entities"
)

type ProjectWithTrainingDataset struct {
	Project         entities.Project
	TrainingDataset *entities.TrainingDataset
	Finetune        *entities.Finetune
	Deployments     []entities.Deployment
}

type ListProjectsResult struct {
	Projects []ProjectWithTrainingDataset
}

type ListProjectsUseCase interface {
	ListProjects(command ListProjectsCommand) (*ListProjectsResult, error)
}