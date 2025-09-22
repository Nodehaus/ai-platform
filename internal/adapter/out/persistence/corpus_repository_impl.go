package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"ai-platform/internal/application/domain/entities"
)

type CorpusRepositoryImpl struct {
	Db *sql.DB
}


func (r *CorpusRepositoryImpl) Create(ctx context.Context, corpus *entities.Corpus) error {
	query := `INSERT INTO corpus (id, name, s3_path, files_subset, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	corpus.CreatedAt = now
	corpus.UpdatedAt = now

	var filesSubset interface{}
	if corpus.FilesSubset != nil {
		filesSubset = pq.Array(*corpus.FilesSubset)
	}

	_, err := r.Db.ExecContext(ctx, query,
		corpus.ID,
		corpus.Name,
		corpus.S3Path,
		filesSubset,
		corpus.CreatedAt,
		corpus.UpdatedAt,
	)

	return err
}

func (r *CorpusRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Corpus, error) {
	query := `SELECT id, name, s3_path, files_subset, created_at, updated_at FROM corpus WHERE id = $1`

	var model CorpusRepositoryModel
	var filesSubset pq.StringArray
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Name,
		&model.S3Path,
		&filesSubset,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(filesSubset) > 0 {
		subset := []string(filesSubset)
		model.FilesSubset = &subset
	}

	return model.ToEntity(), nil
}

func (r *CorpusRepositoryImpl) GetByName(ctx context.Context, name string) (*entities.Corpus, error) {
	query := `SELECT id, name, s3_path, files_subset, created_at, updated_at FROM corpus WHERE name = $1`

	var model CorpusRepositoryModel
	var filesSubset pq.StringArray
	err := r.Db.QueryRowContext(ctx, query, name).Scan(
		&model.ID,
		&model.Name,
		&model.S3Path,
		&filesSubset,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(filesSubset) > 0 {
		subset := []string(filesSubset)
		model.FilesSubset = &subset
	}

	return model.ToEntity(), nil
}

func (r *CorpusRepositoryImpl) Update(ctx context.Context, corpus *entities.Corpus) error {
	query := `UPDATE corpus SET name = $2, s3_path = $3, files_subset = $4, updated_at = $5 WHERE id = $1`

	corpus.UpdatedAt = time.Now()

	var filesSubset interface{}
	if corpus.FilesSubset != nil {
		filesSubset = pq.Array(*corpus.FilesSubset)
	}

	_, err := r.Db.ExecContext(ctx, query,
		corpus.ID,
		corpus.Name,
		corpus.S3Path,
		filesSubset,
		corpus.UpdatedAt,
	)

	return err
}

func (r *CorpusRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM corpus WHERE id = $1`
	_, err := r.Db.ExecContext(ctx, query, id)
	return err
}

func (r *CorpusRepositoryImpl) List(ctx context.Context) ([]*entities.Corpus, error) {
	query := `SELECT id, name, s3_path, files_subset, created_at, updated_at FROM corpus ORDER BY name`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var corpus []*entities.Corpus
	for rows.Next() {
		var model CorpusRepositoryModel
		var filesSubset pq.StringArray
		err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.S3Path,
			&filesSubset,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if len(filesSubset) > 0 {
			subset := []string(filesSubset)
			model.FilesSubset = &subset
		}

		corpus = append(corpus, model.ToEntity())
	}

	return corpus, nil
}