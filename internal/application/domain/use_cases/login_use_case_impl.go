package use_cases

import (
	"errors"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type LoginUseCaseImpl struct {
	userRepository persistence.UserRepository
	userService    *services.UserService
	jwtService     *services.JWTService
}

func NewLoginUseCase(userRepository persistence.UserRepository, userService *services.UserService, jwtService *services.JWTService) in.LoginUseCase {
	return &LoginUseCaseImpl{
		userRepository: userRepository,
		userService:    userService,
		jwtService:     jwtService,
	}
}

func (uc *LoginUseCaseImpl) Login(command in.LoginCommand) (*in.LoginResult, error) {
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

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &in.LoginResult{
		User:  user,
		Token: token,
	}, nil
}