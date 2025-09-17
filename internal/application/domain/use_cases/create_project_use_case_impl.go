package use_cases

import (
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type CreateProjectUseCaseImpl struct {
	projectRepository persistence.ProjectRepository
	projectService    *services.ProjectService
}

func NewCreateProjectUseCase(projectRepository persistence.ProjectRepository, projectService *services.ProjectService) in.CreateProjectUseCase {
	return &CreateProjectUseCaseImpl{
		projectRepository: projectRepository,
		projectService:    projectService,
	}
}

func (uc *CreateProjectUseCaseImpl) CreateProject(command in.CreateProjectCommand) (*in.CreateProjectResult, error) {
	err := uc.projectService.ValidateProjectName(command.Name)
	if err != nil {
		return nil, err
	}

	err = uc.projectService.ValidateProjectNameUniqueness(command.Name, command.OwnerID, uc.projectRepository.ExistsByNameAndOwnerID)
	if err != nil {
		return nil, err
	}

	project := uc.projectService.CreateProject(command.Name, command.OwnerID)

	err = uc.projectRepository.Create(project)
	if err != nil {
		return nil, err
	}

	return &in.CreateProjectResult{
		Project: project,
	}, nil
}