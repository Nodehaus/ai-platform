package persistence

import (
	"ai-platform/internal/application/domain/entities"
	"time"

	"github.com/google/uuid"
)

type UserRepositoryModel struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *UserRepositoryModel) ToEntity() *entities.User {
	return &entities.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}