package use_cases

import (
	"errors"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type LoginUseCaseImpl struct {
	UserRepository persistence.UserRepository
	UserService    *services.UserService
	JwtService     *services.JWTService
}


func (uc *LoginUseCaseImpl) Login(command in.LoginCommand) (*in.LoginResult, error) {
	if command.Email == "" {
		return nil, errors.New("email is required")
	}

	if command.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := uc.UserRepository.FindByEmail(command.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = uc.UserService.ValidatePassword(user.Password, command.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := uc.JwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &in.LoginResult{
		User:  user,
		Token: token,
	}, nil
}