package use_cases

import (
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type ListProjectsUseCaseImpl struct {
	ProjectRepository persistence.ProjectRepository
}


func (uc *ListProjectsUseCaseImpl) ListProjects(command in.ListProjectsCommand) (*in.ListProjectsResult, error) {
	projects, err := uc.ProjectRepository.GetActiveByOwnerID(command.OwnerID)
	if err != nil {
		return nil, err
	}

	return &in.ListProjectsResult{
		Projects: projects,
	}, nil
}