package persistence

import (
	"database/sql"
	"ai-platform/internal/application/domain/entities"
)

type UserRepositoryImpl struct {
	Db *sql.DB
}


func (r *UserRepositoryImpl) FindByEmail(email string) (*entities.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	var model UserRepositoryModel
	err := r.Db.QueryRow(query, email).Scan(
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