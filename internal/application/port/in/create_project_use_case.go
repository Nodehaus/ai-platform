package in

import "ai-platform/internal/application/domain/entities"

type CreateProjectResult struct {
	Project *entities.Project
}

type CreateProjectUseCase interface {
	CreateProject(command CreateProjectCommand) (*CreateProjectResult, error)
}