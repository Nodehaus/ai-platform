package in

import "ai-platform/internal/application/domain/entities"

type LoginResult struct {
	User  *entities.User
	Token string
}

type LoginUseCase interface {
	Login(command LoginCommand) (*LoginResult, error)
}