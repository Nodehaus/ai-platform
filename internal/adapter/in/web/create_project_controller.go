package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-platform/internal/application/port/in"
)

type CreateProjectController struct {
	CreateProjectUseCase in.CreateProjectUseCase
}


func (c *CreateProjectController) CreateProject(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	var request CreateProjectRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	command := in.CreateProjectCommand{
		Name:    request.Name,
		OwnerID: userID,
	}

	result, err := c.CreateProjectUseCase.CreateProject(command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := NewCreateProjectResponse(result.Project, "Project created successfully")
	ctx.JSON(http.StatusCreated, response)
}