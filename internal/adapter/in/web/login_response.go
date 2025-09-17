package web

import (
	"ai-platform/internal/application/domain/entities"

	"github.com/google/uuid"
)

type LoginResponse struct {
	User    *UserResponse `json:"user"`
	Message string        `json:"message"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func NewLoginResponse(user *entities.User, message string) *LoginResponse {
	return &LoginResponse{
		User: &UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Message: message,
	}
}