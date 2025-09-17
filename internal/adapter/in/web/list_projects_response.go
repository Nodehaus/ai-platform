package web

import "ai-platform/internal/application/domain/entities"

type ListProjectsResponse struct {
	Projects []ProjectResponse `json:"projects"`
}

func NewListProjectsResponse(projects []entities.Project) *ListProjectsResponse {
	projectResponses := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = ProjectResponse{
			ID:        project.ID,
			Name:      project.Name,
			Status:    project.Status,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
		}
	}

	return &ListProjectsResponse{
		Projects: projectResponses,
	}
}