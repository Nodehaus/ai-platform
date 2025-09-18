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
	protected.POST("/projects", s.createProjectController.CreateProject)
	protected.GET("/projects", s.listProjectsController.ListProjects)
	protected.POST("/projects/:project_id/training-datasets", s.createTrainingDatasetController.CreateTrainingDataset)

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

	r.POST("/web/projects/create", func(c *gin.Context) {
		web.CreateProjectHandler(c.Writer, c.Request)
	})

	r.GET("/web/projects/:project_id/training-datasets/step1", func(c *gin.Context) {
		web.TrainingDatasetStep1Handler(c.Writer, c.Request)
	})



	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

