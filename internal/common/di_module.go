package common

import (
	"go.uber.org/fx"

	"ai-platform/internal/adapter/in/web"
	"ai-platform/internal/adapter/out/persistence"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/domain/use_cases"
	"ai-platform/internal/application/port/in"
	persistencePort "ai-platform/internal/application/port/out/persistence"
	"ai-platform/internal/database"
)

func NewUserRepository(dbService database.Service) persistencePort.UserRepository {
	return persistence.NewUserRepository(dbService.GetDB())
}

func NewUserService() *services.UserService {
	return services.NewUserService()
}

func NewLoginUseCase(userRepo persistencePort.UserRepository, userService *services.UserService) in.LoginUseCase {
	return use_cases.NewLoginUseCase(userRepo, userService)
}

func NewLoginController(loginUseCase in.LoginUseCase) *web.LoginController {
	return web.NewLoginController(loginUseCase)
}

var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewUserService),
	fx.Provide(NewLoginUseCase),
	fx.Provide(NewLoginController),
)