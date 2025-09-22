package web

import (
	"ai-platform/internal/application/domain/entities"
	"time"

	"github.com/google/uuid"
)


type ProjectResponse struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Status    entities.ProjectStatus `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type CreateProjectResponse struct {
	Project ProjectResponse `json:"project"`
	Message string          `json:"message"`
}

func NewCreateProjectResponse(project *entities.Project) *CreateProjectResponse {
	return &CreateProjectResponse{
		Project: ProjectResponse{
			ID:        project.ID,
			Name:      project.Name,
			Status:    project.Status,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
		},
		Message: "Project created successfully",
	}
}