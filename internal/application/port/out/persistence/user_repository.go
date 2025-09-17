package persistence

import "ai-platform/internal/application/domain/entities"

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
}