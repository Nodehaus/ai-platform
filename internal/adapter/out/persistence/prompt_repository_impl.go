package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/out/persistence"
)

type PromptRepositoryImpl struct {
	db *sql.DB
}

func NewPromptRepository(db *sql.DB) persistence.PromptRepository {
	return &PromptRepositoryImpl{
		db: db,
	}
}

func (r *PromptRepositoryImpl) Create(ctx context.Context, prompt *entities.Prompt) error {
	query := `INSERT INTO prompts (id, version, text, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5)`

	now := time.Now()
	prompt.CreatedAt = now
	prompt.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		prompt.ID,
		prompt.Version,
		prompt.Text,
		prompt.CreatedAt,
		prompt.UpdatedAt,
	)

	return err
}

func (r *PromptRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Prompt, error) {
	query := `SELECT id, version, text, created_at, updated_at FROM prompts WHERE id = $1`

	var model PromptRepositoryModel
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Version,
		&model.Text,
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

func (r *PromptRepositoryImpl) Update(ctx context.Context, prompt *entities.Prompt) error {
	query := `UPDATE prompts SET version = $2, text = $3, updated_at = $4 WHERE id = $1`

	prompt.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		prompt.ID,
		prompt.Version,
		prompt.Text,
		prompt.UpdatedAt,
	)

	return err
}

func (r *PromptRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM prompts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}