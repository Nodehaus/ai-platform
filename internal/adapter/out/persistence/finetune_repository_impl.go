package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type FinetuneRepositoryImpl struct {
	Db *sql.DB
}

func (r *FinetuneRepositoryImpl) Create(ctx context.Context, finetune *entities.Finetune) error {
	query := `INSERT INTO finetunes (
		id, project_id, version, model_name, base_model_name,
		model_size_gb, model_size_parameter, model_dtype, model_quantization,
		inference_samples_json, training_dataset_id, training_dataset_number_examples,
		training_dataset_select_random, training_time_seconds, status, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	now := time.Now()
	finetune.CreatedAt = now
	finetune.UpdatedAt = now

	model, err := FromFinetuneEntity(finetune)
	if err != nil {
		return err
	}

	_, err = r.Db.ExecContext(ctx, query,
		model.ID,
		model.ProjectID,
		model.Version,
		model.ModelName,
		model.BaseModelName,
		model.ModelSizeGB,
		model.ModelSizeParameter,
		model.ModelDtype,
		model.ModelQuantization,
		model.InferenceSamplesJSON,
		model.TrainingDatasetID,
		model.TrainingDatasetNumberExamples,
		model.TrainingDatasetSelectRandom,
		model.TrainingTimeSeconds,
		model.Status,
		model.CreatedAt,
		model.UpdatedAt,
	)

	return err
}

func (r *FinetuneRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Finetune, error) {
	query := `SELECT
		id, project_id, version, model_name, base_model_name,
		model_size_gb, model_size_parameter, model_dtype, model_quantization,
		inference_samples_json, training_dataset_id, training_dataset_number_examples,
		training_dataset_select_random, training_time_seconds, status, created_at, updated_at
	FROM finetunes WHERE id = $1`

	var model FinetuneRepositoryModel
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Version,
		&model.ModelName,
		&model.BaseModelName,
		&model.ModelSizeGB,
		&model.ModelSizeParameter,
		&model.ModelDtype,
		&model.ModelQuantization,
		&model.InferenceSamplesJSON,
		&model.TrainingDatasetID,
		&model.TrainingDatasetNumberExamples,
		&model.TrainingDatasetSelectRandom,
		&model.TrainingTimeSeconds,
		&model.Status,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity()
}

func (r *FinetuneRepositoryImpl) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Finetune, error) {
	query := `SELECT
		id, project_id, version, model_name, base_model_name,
		model_size_gb, model_size_parameter, model_dtype, model_quantization,
		inference_samples_json, training_dataset_id, training_dataset_number_examples,
		training_dataset_select_random, training_time_seconds, status, created_at, updated_at
	FROM finetunes WHERE project_id = $1 ORDER BY version DESC`

	rows, err := r.Db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var finetunes []*entities.Finetune
	for rows.Next() {
		var model FinetuneRepositoryModel
		err := rows.Scan(
			&model.ID,
			&model.ProjectID,
			&model.Version,
			&model.ModelName,
			&model.BaseModelName,
			&model.ModelSizeGB,
			&model.ModelSizeParameter,
			&model.ModelDtype,
			&model.ModelQuantization,
			&model.InferenceSamplesJSON,
			&model.TrainingDatasetID,
			&model.TrainingDatasetNumberExamples,
			&model.TrainingDatasetSelectRandom,
			&model.TrainingTimeSeconds,
			&model.Status,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		finetune, err := model.ToEntity()
		if err != nil {
			return nil, err
		}

		finetunes = append(finetunes, finetune)
	}

	return finetunes, nil
}

func (r *FinetuneRepositoryImpl) GetLatestByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.Finetune, error) {
	query := `SELECT
		id, project_id, version, model_name, base_model_name,
		model_size_gb, model_size_parameter, model_dtype, model_quantization,
		inference_samples_json, training_dataset_id, training_dataset_number_examples,
		training_dataset_select_random, training_time_seconds, status, created_at, updated_at
	FROM finetunes WHERE project_id = $1 ORDER BY version DESC LIMIT 1`

	var model FinetuneRepositoryModel
	err := r.Db.QueryRowContext(ctx, query, projectID).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Version,
		&model.ModelName,
		&model.BaseModelName,
		&model.ModelSizeGB,
		&model.ModelSizeParameter,
		&model.ModelDtype,
		&model.ModelQuantization,
		&model.InferenceSamplesJSON,
		&model.TrainingDatasetID,
		&model.TrainingDatasetNumberExamples,
		&model.TrainingDatasetSelectRandom,
		&model.TrainingTimeSeconds,
		&model.Status,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity()
}

func (r *FinetuneRepositoryImpl) Update(ctx context.Context, finetune *entities.Finetune) error {
	query := `UPDATE finetunes SET
		model_name = $1, base_model_name = $2, model_size_gb = $3, model_size_parameter = $4,
		model_dtype = $5, model_quantization = $6, inference_samples_json = $7,
		training_dataset_number_examples = $8, training_dataset_select_random = $9,
		training_time_seconds = $10, status = $11, updated_at = $12
	WHERE id = $13`

	finetune.UpdatedAt = time.Now()

	model, err := FromFinetuneEntity(finetune)
	if err != nil {
		return err
	}

	_, err = r.Db.ExecContext(ctx, query,
		model.ModelName,
		model.BaseModelName,
		model.ModelSizeGB,
		model.ModelSizeParameter,
		model.ModelDtype,
		model.ModelQuantization,
		model.InferenceSamplesJSON,
		model.TrainingDatasetNumberExamples,
		model.TrainingDatasetSelectRandom,
		model.TrainingTimeSeconds,
		model.Status,
		model.UpdatedAt,
		model.ID,
	)

	return err
}

func (r *FinetuneRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.FinetuneStatus) error {
	query := `UPDATE finetunes SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.Db.ExecContext(ctx, query, string(status), time.Now(), id)
	return err
}

func (r *FinetuneRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE finetunes SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.Db.ExecContext(ctx, query, string(entities.FinetuneStatusDeleted), time.Now(), id)
	return err
}

func (r *FinetuneRepositoryImpl) GetNextVersion(ctx context.Context, projectID uuid.UUID) (int, error) {
	query := `SELECT COALESCE(MAX(version), 0) + 1 FROM finetunes WHERE project_id = $1`
	var nextVersion int
	err := r.Db.QueryRowContext(ctx, query, projectID).Scan(&nextVersion)
	return nextVersion, err
}