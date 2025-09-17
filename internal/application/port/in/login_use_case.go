package in

import "ai-platform/internal/application/domain/entities"

type LoginUseCase interface {
	Login(command LoginCommand) (*entities.User, error)
}