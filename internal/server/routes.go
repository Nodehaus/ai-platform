package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ai-platform/cmd/web"
	"io/fs"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Public routes (no authentication required)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/web/home")
	})
	r.GET("/health", s.healthHandler)
	r.POST("/api/login", s.loginController.Login)

	// Protected routes (authentication required)
	protected := r.Group("/api")
	protected.Use(s.authMiddleware.RequireAuth())

	staticFiles, _ := fs.Sub(web.Files, "assets")
	r.StaticFS("/assets", http.FS(staticFiles))

	// Frontend routes
	r.GET("/web", func(c *gin.Context) {
		// Redirect to home page (which will redirect to login if not authenticated)
		c.Redirect(302, "/web/home")
	})

	r.GET("/web/login", func(c *gin.Context) {
		web.LoginPageHandler(c.Writer, c.Request)
	})

	r.POST("/web/login", func(c *gin.Context) {
		web.LoginFormHandler(c.Writer, c.Request)
	})

	r.GET("/web/home", func(c *gin.Context) {
		web.HomepageHandler(c.Writer, c.Request)
	})

	r.GET("/web/logout", func(c *gin.Context) {
		web.LogoutHandler(c.Writer, c.Request)
	})


	// Example protected endpoint
	protected.GET("/profile", s.profileHandler)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) profileHandler(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	email, exists := GetUserEmailFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User email not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
		"message": "This is a protected endpoint",
	})
}
