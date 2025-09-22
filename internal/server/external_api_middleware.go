package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ExternalAPIMiddleware struct{}

func (m *ExternalAPIMiddleware) RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key is required",
			})
			c.Abort()
			return
		}

		expectedAPIKey := os.Getenv("APP_EXTERNAL_API_KEY")
		if expectedAPIKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "External API key not configured",
			})
			c.Abort()
			return
		}

		if apiKey != expectedAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}