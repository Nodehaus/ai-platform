package use_cases

import (
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type ListProjectsUseCaseImpl struct {
	projectRepository persistence.ProjectRepository
}

func NewListProjectsUseCase(projectRepository persistence.ProjectRepository) in.ListProjectsUseCase {
	return &ListProjectsUseCaseImpl{
		projectRepository: projectRepository,
	}
}

func (uc *ListProjectsUseCaseImpl) ListProjects(command in.ListProjectsCommand) (*in.ListProjectsResult, error) {
	projects, err := uc.projectRepository.GetActiveByOwnerID(command.OwnerID)
	if err != nil {
		return nil, err
	}

	return &in.ListProjectsResult{
		Projects: projects,
	}, nil
}