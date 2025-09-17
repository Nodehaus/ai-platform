package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-platform/internal/application/port/in"
)

type ListProjectsController struct {
	listProjectsUseCase in.ListProjectsUseCase
}

func NewListProjectsController(listProjectsUseCase in.ListProjectsUseCase) *ListProjectsController {
	return &ListProjectsController{
		listProjectsUseCase: listProjectsUseCase,
	}
}

func (c *ListProjectsController) ListProjects(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	command := in.ListProjectsCommand{
		OwnerID: userID,
	}

	result, err := c.listProjectsUseCase.ListProjects(command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch projects",
		})
		return
	}

	response := NewListProjectsResponse(result.Projects)
	ctx.JSON(http.StatusOK, response)
}