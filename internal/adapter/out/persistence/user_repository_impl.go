package persistence

import (
	"database/sql"
	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/out/persistence"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) persistence.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*entities.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	var model UserRepositoryModel
	err := r.db.QueryRow(query, email).Scan(
		&model.ID,
		&model.Email,
		&model.Password,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}