package use_cases

import (
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"context"
)

type ListProjectsUseCaseImpl struct {
	ProjectService *services.ProjectService
}


func (uc *ListProjectsUseCaseImpl) ListProjects(command in.ListProjectsCommand) (*in.ListProjectsResult, error) {
	projectsWithDatasets, err := uc.ProjectService.ListProjects(context.Background(), command.OwnerID)
	if err != nil {
		return nil, err
	}

	result := make([]in.ProjectWithTrainingDataset, len(projectsWithDatasets))
	for i, projectWithDataset := range projectsWithDatasets {
		result[i] = in.ProjectWithTrainingDataset{
			Project:         projectWithDataset.Project,
			TrainingDataset: projectWithDataset.TrainingDataset,
		}
	}

	return &in.ListProjectsResult{
		Projects: result,
	}, nil
}