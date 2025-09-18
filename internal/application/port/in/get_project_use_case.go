package in

import "ai-platform/internal/application/domain/entities"

type GetProjectResult struct {
	Project entities.Project
}

type GetProjectUseCase interface {
	GetProject(command GetProjectCommand) (*GetProjectResult, error)
}