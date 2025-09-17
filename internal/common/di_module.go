package common

import (
	"os"

	"go.uber.org/fx"

	"ai-platform/internal/adapter/in/web"
	"ai-platform/internal/adapter/out/persistence"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/domain/use_cases"
	"ai-platform/internal/application/port/in"
	persistencePort "ai-platform/internal/application/port/out/persistence"
	"ai-platform/internal/database"
	"ai-platform/internal/server"
)

func NewUserRepository(dbService database.Service) persistencePort.UserRepository {
	return persistence.NewUserRepository(dbService.GetDB())
}

func NewProjectRepository(dbService database.Service) persistencePort.ProjectRepository {
	return persistence.NewProjectRepository(dbService.GetDB())
}

func NewUserService() *services.UserService {
	return services.NewUserService()
}

func NewProjectService() *services.ProjectService {
	return services.NewProjectService()
}

func NewJWTService() *services.JWTService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-this-in-production"
	}
	return services.NewJWTService(secretKey)
}

func NewLoginUseCase(userRepo persistencePort.UserRepository, userService *services.UserService, jwtService *services.JWTService) in.LoginUseCase {
	return use_cases.NewLoginUseCase(userRepo, userService, jwtService)
}

func NewCreateProjectUseCase(projectRepo persistencePort.ProjectRepository, projectService *services.ProjectService) in.CreateProjectUseCase {
	return use_cases.NewCreateProjectUseCase(projectRepo, projectService)
}

func NewLoginController(loginUseCase in.LoginUseCase) *web.LoginController {
	return web.NewLoginController(loginUseCase)
}

func NewCreateProjectController(createProjectUseCase in.CreateProjectUseCase) *web.CreateProjectController {
	return web.NewCreateProjectController(createProjectUseCase)
}

func NewAuthMiddleware(jwtService *services.JWTService) *server.AuthMiddleware {
	return server.NewAuthMiddleware(jwtService)
}

var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewProjectRepository),
	fx.Provide(NewUserService),
	fx.Provide(NewProjectService),
	fx.Provide(NewJWTService),
	fx.Provide(NewLoginUseCase),
	fx.Provide(NewCreateProjectUseCase),
	fx.Provide(NewLoginController),
	fx.Provide(NewCreateProjectController),
	fx.Provide(NewAuthMiddleware),
)