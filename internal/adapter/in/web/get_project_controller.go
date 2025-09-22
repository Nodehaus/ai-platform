package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type GetProjectController struct {
	GetProjectUseCase in.GetProjectUseCase
}

func (c *GetProjectController) GetProject(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	projectIDStr := ctx.Param("project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID format",
		})
		return
	}

	command := in.GetProjectCommand{
		ProjectID: projectID,
		OwnerID:   userID,
	}

	result, err := c.GetProjectUseCase.GetProject(command)
	if err != nil {
		if err.Error() == "project not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		if err.Error() == "access denied" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch project",
		})
		return
	}

	response := NewGetProjectResponse(&result.Project)
	ctx.JSON(http.StatusOK, response)
}