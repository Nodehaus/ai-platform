package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/out/persistence"
)

type APIKeyMiddleware struct {
	DeploymentRepository persistence.DeploymentRepository
}

func NewAPIKeyMiddleware(deploymentRepo persistence.DeploymentRepository) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		DeploymentRepository: deploymentRepo,
	}
}

// AuthenticateAPIKey validates the API key from Authorization header
// and sets the deployment and finetune information in the context
func (m *APIKeyMiddleware) AuthenticateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Authorization header",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		apiKey := parts[1]

		// Extract project_id from URL
		projectIDStr := c.Param("project_id")
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid project ID format",
			})
			c.Abort()
			return
		}

		// Find deployment by API key
		deployment, err := m.DeploymentRepository.GetByAPIKey(apiKey)
		if err != nil || deployment == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		// Verify the deployment belongs to the specified project
		if deployment.ProjectID != projectID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key does not belong to this project",
			})
			c.Abort()
			return
		}

		// Store deployment information in context
		c.Set("deployment_id", deployment.ID)
		c.Set("model_name", deployment.ModelName)
		if deployment.FinetuneID != nil {
			c.Set("finetune_id", *deployment.FinetuneID)
		}
		c.Set("project_id", deployment.ProjectID)

		c.Next()
	}
}
