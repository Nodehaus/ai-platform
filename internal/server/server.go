package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"ai-platform/internal/adapter/in/web"
	"ai-platform/internal/database"
)

type Server struct {
	port                                     int
	db                                       database.Service
	loginController                          *web.LoginController
	createProjectController                  *web.CreateProjectController
	getProjectController                     *web.GetProjectController
	listProjectsController                   *web.ListProjectsController
	createTrainingDatasetController          *web.CreateTrainingDatasetController
	getTrainingDatasetController             *web.GetTrainingDatasetController
	downloadTrainingDatasetController        *web.DownloadTrainingDatasetController
	updateTrainingDatasetStatusController    *web.UpdateTrainingDatasetStatusController
	updateFinetuneStatusController           *web.UpdateFinetuneStatusController
	createFinetuneController                 *web.CreateFinetuneController
	getFinetuneController                    *web.GetFinetuneController
	downloadModelController                  *web.DownloadModelController
	authMiddleware                           *AuthMiddleware
	externalAPIMiddleware                    *ExternalAPIMiddleware
}

func NewServer(db database.Service, loginController *web.LoginController, createProjectController *web.CreateProjectController, getProjectController *web.GetProjectController, listProjectsController *web.ListProjectsController, createTrainingDatasetController *web.CreateTrainingDatasetController, getTrainingDatasetController *web.GetTrainingDatasetController, downloadTrainingDatasetController *web.DownloadTrainingDatasetController, updateTrainingDatasetStatusController *web.UpdateTrainingDatasetStatusController, updateFinetuneStatusController *web.UpdateFinetuneStatusController, createFinetuneController *web.CreateFinetuneController, getFinetuneController *web.GetFinetuneController, downloadModelController *web.DownloadModelController, authMiddleware *AuthMiddleware, externalAPIMiddleware *ExternalAPIMiddleware) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	serverInstance := &Server{
		port:                                     port,
		db:                                       db,
		loginController:                          loginController,
		createProjectController:                  createProjectController,
		getProjectController:                     getProjectController,
		listProjectsController:                   listProjectsController,
		createTrainingDatasetController:          createTrainingDatasetController,
		getTrainingDatasetController:             getTrainingDatasetController,
		downloadTrainingDatasetController:        downloadTrainingDatasetController,
		updateTrainingDatasetStatusController:    updateTrainingDatasetStatusController,
		updateFinetuneStatusController:           updateFinetuneStatusController,
		createFinetuneController:                 createFinetuneController,
		getFinetuneController:                    getFinetuneController,
		downloadModelController:                  downloadModelController,
		authMiddleware:                           authMiddleware,
		externalAPIMiddleware:                    externalAPIMiddleware,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverInstance.port),
		Handler:      serverInstance.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
