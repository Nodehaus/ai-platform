package use_cases

import (
	"errors"
	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type LoginUseCaseImpl struct {
	userRepository persistence.UserRepository
	userService    *services.UserService
}

func NewLoginUseCase(userRepository persistence.UserRepository, userService *services.UserService) in.LoginUseCase {
	return &LoginUseCaseImpl{
		userRepository: userRepository,
		userService:    userService,
	}
}

func (uc *LoginUseCaseImpl) Login(command in.LoginCommand) (*entities.User, error) {
	if command.Email == "" {
		return nil, errors.New("email is required")
	}

	if command.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := uc.userRepository.FindByEmail(command.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = uc.userService.ValidatePassword(user.Password, command.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}