package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type GetDeploymentController struct {
	GetDeploymentUseCase in.GetDeploymentUseCase
}

func (c *GetDeploymentController) GetDeployment(ctx *gin.Context) {
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

	deploymentIDStr := ctx.Param("deployment_id")
	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deployment ID format",
		})
		return
	}

	command := in.GetDeploymentCommand{
		DeploymentID: deploymentID,
		ProjectID:    projectID,
		OwnerID:      userID,
	}

	result, err := c.GetDeploymentUseCase.GetDeployment(command)
	if err != nil {
		if err.Error() == "deployment not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Deployment not found",
			})
			return
		}
		if err.Error() == "access denied" || err.Error() == "project not found" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}
		if err.Error() == "deployment does not belong to this project" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch deployment",
		})
		return
	}

	response := NewGetDeploymentResponse(result.Deployment, result.Logs)
	ctx.JSON(http.StatusOK, response)
}
