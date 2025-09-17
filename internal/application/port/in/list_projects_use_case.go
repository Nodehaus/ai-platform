package in

import "ai-platform/internal/application/domain/entities"

type ListProjectsResult struct {
	Projects []entities.Project
}

type ListProjectsUseCase interface {
	ListProjects(command ListProjectsCommand) (*ListProjectsResult, error)
}