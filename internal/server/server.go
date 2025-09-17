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
	port                     int
	db                       database.Service
	loginController          *web.LoginController
	createProjectController  *web.CreateProjectController
	authMiddleware           *AuthMiddleware
}

func NewServer(db database.Service, loginController *web.LoginController, createProjectController *web.CreateProjectController, authMiddleware *AuthMiddleware) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	serverInstance := &Server{
		port:                    port,
		db:                      db,
		loginController:         loginController,
		createProjectController: createProjectController,
		authMiddleware:          authMiddleware,
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
