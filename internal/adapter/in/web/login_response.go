package web

import (
	"ai-platform/internal/application/domain/entities"

	"github.com/google/uuid"
)

type LoginResponse struct {
	User    *UserResponse `json:"user"`
	Token   string        `json:"token"`
	Message string        `json:"message"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func NewLoginResponse(user *entities.User, token string, message string) *LoginResponse {
	return &LoginResponse{
		User: &UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Token:   token,
		Message: message,
	}
}