package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/out/persistence"
)

type CorpusRepositoryImpl struct {
	db *sql.DB
}

func NewCorpusRepository(db *sql.DB) persistence.CorpusRepository {
	return &CorpusRepositoryImpl{
		db: db,
	}
}

func (r *CorpusRepositoryImpl) Create(ctx context.Context, corpus *entities.Corpus) error {
	query := `INSERT INTO corpus (id, name, s3_path, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5)`

	now := time.Now()
	corpus.CreatedAt = now
	corpus.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		corpus.ID,
		corpus.Name,
		corpus.S3Path,
		corpus.CreatedAt,
		corpus.UpdatedAt,
	)

	return err
}

func (r *CorpusRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Corpus, error) {
	query := `SELECT id, name, s3_path, created_at, updated_at FROM corpus WHERE id = $1`

	var model CorpusRepositoryModel
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Name,
		&model.S3Path,
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

func (r *CorpusRepositoryImpl) GetByName(ctx context.Context, name string) (*entities.Corpus, error) {
	query := `SELECT id, name, s3_path, created_at, updated_at FROM corpus WHERE name = $1`

	var model CorpusRepositoryModel
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&model.ID,
		&model.Name,
		&model.S3Path,
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

func (r *CorpusRepositoryImpl) Update(ctx context.Context, corpus *entities.Corpus) error {
	query := `UPDATE corpus SET name = $2, s3_path = $3, updated_at = $4 WHERE id = $1`

	corpus.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		corpus.ID,
		corpus.Name,
		corpus.S3Path,
		corpus.UpdatedAt,
	)

	return err
}

func (r *CorpusRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM corpus WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *CorpusRepositoryImpl) List(ctx context.Context) ([]*entities.Corpus, error) {
	query := `SELECT id, name, s3_path, created_at, updated_at FROM corpus ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var corpus []*entities.Corpus
	for rows.Next() {
		var model CorpusRepositoryModel
		err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.S3Path,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		corpus = append(corpus, model.ToEntity())
	}

	return corpus, nil
}