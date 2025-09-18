package use_cases

import (
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type CreateProjectUseCaseImpl struct {
	ProjectRepository persistence.ProjectRepository
	ProjectService    *services.ProjectService
}


func (uc *CreateProjectUseCaseImpl) CreateProject(command in.CreateProjectCommand) (*in.CreateProjectResult, error) {
	err := uc.ProjectService.ValidateProjectName(command.Name)
	if err != nil {
		return nil, err
	}

	err = uc.ProjectService.ValidateProjectNameUniqueness(command.Name, command.OwnerID, uc.ProjectRepository.ExistsByNameAndOwnerID)
	if err != nil {
		return nil, err
	}

	project := uc.ProjectService.CreateProject(command.Name, command.OwnerID)

	err = uc.ProjectRepository.Create(project)
	if err != nil {
		return nil, err
	}

	return &in.CreateProjectResult{
		Project: project,
	}, nil
}