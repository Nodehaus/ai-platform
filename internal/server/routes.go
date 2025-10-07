package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ai-platform/cmd/web"
	"ai-platform/cmd/web/training_datasets"
	"ai-platform/cmd/web/finetunes"
	"ai-platform/cmd/web/deployments"
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
	protected.POST("/analyze-training-dataset-prompt", s.analyzePromptController.AnalyzePrompt)
	protected.POST("/projects", s.createProjectController.CreateProject)
	protected.GET("/projects", s.listProjectsController.ListProjects)
	protected.GET("/projects/:project_id", s.getProjectController.GetProject)
	protected.POST("/projects/:project_id/training-datasets", s.createTrainingDatasetController.CreateTrainingDataset)
	protected.POST("/projects/:project_id/training-datasets/upload", s.uploadNewTrainingDatasetVersionController.UploadNewTrainingDatasetVersion)
	protected.GET("/projects/:project_id/training-datasets/:training_dataset_id", s.getTrainingDatasetController.GetTrainingDataset)
	protected.GET("/projects/:project_id/training-datasets/:training_dataset_id/download", s.downloadTrainingDatasetController.DownloadTrainingDataset)
	protected.POST("/projects/:project_id/training-datasets/:training_dataset_id/upload", s.uploadTrainingDatasetController.UploadTrainingDataset)
	protected.POST("/projects/:project_id/finetunes", s.createFinetuneController.CreateFinetune)
	protected.GET("/projects/:project_id/finetunes/:finetune_id", s.getFinetuneController.GetFinetune)
	protected.POST("/projects/:project_id/finetunes/:finetune_id/completion", s.finetuneCompletionController.GenerateCompletion)
	protected.GET("/projects/:project_id/finetunes/:finetune_id/download", s.downloadModelController.DownloadModel)
	protected.POST("/projects/:project_id/deployments", s.createDeploymentController.CreateDeployment)
	protected.GET("/projects/:project_id/deployments/:deployment_id", s.getDeploymentController.GetDeployment)

	// External API routes (API key protected)
	external := r.Group("/api/external")
	external.Use(s.externalAPIMiddleware.RequireAPIKey())
	external.PUT("/training-datasets/:training_dataset_id/update-status", s.updateTrainingDatasetStatusController.UpdateStatus)
	external.PUT("/finetunes/:finetune_id/update-status", s.updateFinetuneStatusController.UpdateStatus)

	// Public OpenAI-compatible API routes (deployment API key protected)
	publicAPI := r.Group("/public/:project_id")
	publicAPI.Use(s.apiKeyMiddleware.AuthenticateAPIKey())
	publicAPI.POST("/completions", s.publicCompletionController.GenerateCompletion)
	publicAPI.POST("/chat/completions", s.publicChatCompletionController.GenerateChatCompletion)

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
		web.HomeHandler(c.Writer, c.Request)
	})

	r.GET("/web/logout", func(c *gin.Context) {
		web.LogoutHandler(c.Writer, c.Request)
	})

	r.POST("/web/projects/create", func(c *gin.Context) {
		web.CreateProjectHandler(c.Writer, c.Request)
	})

	r.POST("/web/projects/:project_id/training-datasets/create", func(c *gin.Context) {
		training_datasets.CreateTrainingDatasetHandler(c.Writer, c.Request)
	})

	r.GET("/web/projects/:project_id/training-datasets/step1", func(c *gin.Context) {
		training_datasets.TrainingDatasetStep1Handler(c.Writer, c.Request)
	})

	r.GET("/web/projects/:project_id/training-datasets/step2", func(c *gin.Context) {
		training_datasets.TrainingDatasetStep2Handler(c.Writer, c.Request)
	})

	r.GET("/web/projects/:project_id/training-datasets/step3", func(c *gin.Context) {
		training_datasets.TrainingDatasetStep3Handler(c.Writer, c.Request)
	})

	r.GET("/web/projects/:project_id/training-datasets/step4", func(c *gin.Context) {
		training_datasets.TrainingDatasetStep4Handler(c.Writer, c.Request)
	})

	r.GET("/web/projects/:project_id/training-datasets/:training_dataset_id", func(c *gin.Context) {
		training_datasets.TrainingDatasetIndexHandler(c.Writer, c.Request)
	})

	r.POST("/web/projects/:project_id/finetunes/create", func(c *gin.Context) {
		training_datasets.CreateFinetuneHandler(c.Writer, c.Request)
	})

	r.GET("/web/projects/:project_id/finetunes/:finetune_id", func(c *gin.Context) {
		finetunes.FinetuneIndexHandler(c.Writer, c.Request)
	})

	r.GET("/web/projects/:project_id/deployments/:deployment_id", func(c *gin.Context) {
		deployments.DeploymentIndexHandler(c.Writer, c.Request)
	})

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

