package web

import (
	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"time"

	"github.com/google/uuid"
)

type ListProjectsResponse struct {
	Projects []ListProjectResponse `json:"projects"`
}

type TrainingDatasetResponse struct {
	ID     uuid.UUID                     `json:"id"`
	Status entities.TrainingDatasetStatus `json:"status"`
}

type FinetuneResponse struct {
	ID     uuid.UUID              `json:"id"`
	Status entities.FinetuneStatus `json:"status"`
}

type ListProjectResponse struct {
	ID              uuid.UUID                `json:"id"`
	Name            string                   `json:"name"`
	Status          entities.ProjectStatus   `json:"status"`
	TrainingDataset *TrainingDatasetResponse `json:"training_dataset"`
	Finetune        *FinetuneResponse        `json:"finetune"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

func NewListProjectsResponse(projects []in.ProjectWithTrainingDataset) *ListProjectsResponse {
	projectResponses := make([]ListProjectResponse, len(projects))
	for i, projectWithDataset := range projects {
		var trainingDatasetResponse *TrainingDatasetResponse
		if projectWithDataset.TrainingDataset != nil {
			trainingDatasetResponse = &TrainingDatasetResponse{
				ID:     projectWithDataset.TrainingDataset.ID,
				Status: projectWithDataset.TrainingDataset.Status,
			}
		}

		var finetuneResponse *FinetuneResponse
		if projectWithDataset.Finetune != nil {
			finetuneResponse = &FinetuneResponse{
				ID:     projectWithDataset.Finetune.ID,
				Status: projectWithDataset.Finetune.Status,
			}
		}

		projectResponses[i] = ListProjectResponse{
			ID:              projectWithDataset.Project.ID,
			Name:            projectWithDataset.Project.Name,
			Status:          projectWithDataset.Project.Status,
			TrainingDataset: trainingDatasetResponse,
			Finetune:        finetuneResponse,
			CreatedAt:       projectWithDataset.Project.CreatedAt,
			UpdatedAt:       projectWithDataset.Project.UpdatedAt,
		}
	}

	return &ListProjectsResponse{
		Projects: projectResponses,
	}
}