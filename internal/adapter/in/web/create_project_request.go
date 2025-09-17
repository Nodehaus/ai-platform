package web

type CreateProjectRequest struct {
	Name string `json:"name" binding:"required"`
}