package use_cases

import (
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"context"
)

type GetProjectUseCaseImpl struct {
	ProjectService *services.ProjectService
}

func (uc *GetProjectUseCaseImpl) GetProject(command in.GetProjectCommand) (*in.GetProjectResult, error) {
	project, err := uc.ProjectService.GetProject(context.Background(), command.ProjectID, command.OwnerID)
	if err != nil {
		return nil, err
	}

	return &in.GetProjectResult{
		Project: *project,
	}, nil
}