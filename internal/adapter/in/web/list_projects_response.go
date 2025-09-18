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

type ListProjectResponse struct {
	ID                uuid.UUID              `json:"id"`
	Name              string                 `json:"name"`
	Status            entities.ProjectStatus `json:"status"`
	TrainingDatasetID *uuid.UUID             `json:"training_dataset_id"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

func NewListProjectsResponse(projects []in.ProjectWithTrainingDataset) *ListProjectsResponse {
	projectResponses := make([]ListProjectResponse, len(projects))
	for i, projectWithDataset := range projects {
		projectResponses[i] = ListProjectResponse{
			ID:                projectWithDataset.Project.ID,
			Name:              projectWithDataset.Project.Name,
			Status:            projectWithDataset.Project.Status,
			TrainingDatasetID: projectWithDataset.TrainingDatasetID,
			CreatedAt:         projectWithDataset.Project.CreatedAt,
			UpdatedAt:         projectWithDataset.Project.UpdatedAt,
		}
	}

	return &ListProjectsResponse{
		Projects: projectResponses,
	}
}