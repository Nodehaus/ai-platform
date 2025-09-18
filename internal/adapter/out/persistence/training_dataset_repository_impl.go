package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type TrainingDatasetRepositoryImpl struct {
	Db *sql.DB
}


func (r *TrainingDatasetRepositoryImpl) Create(ctx context.Context, trainingDataset *entities.TrainingDataset) error {
	query := `INSERT INTO training_datasets (
		id, project_id, version, generate_model, generate_model_runner,
		generate_gpu_info_card, generate_gpu_info_total_gb, generate_gpu_info_cuda_version,
		input_field, output_field, total_generation_time_seconds,
		generate_prompt_history_ids_json, generate_prompt_id, corpus_id,
		language_iso, status, field_names_json, data_json, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	now := time.Now()
	trainingDataset.CreatedAt = now
	trainingDataset.UpdatedAt = now

	model, err := FromTrainingDatasetEntity(trainingDataset)
	if err != nil {
		return err
	}

	_, err = r.Db.ExecContext(ctx, query,
		model.ID,
		model.ProjectID,
		model.Version,
		model.GenerateModel,
		model.GenerateModelRunner,
		model.GenerateGPUInfoCard,
		model.GenerateGPUInfoTotalGB,
		model.GenerateGPUInfoCudaVersion,
		model.InputField,
		model.OutputField,
		model.TotalGenerationTimeSeconds,
		model.GeneratePromptHistoryIDsJSON,
		model.GeneratePromptID,
		model.CorpusID,
		model.LanguageISO,
		model.Status,
		model.FieldNamesJSON,
		model.DataJSON,
		model.CreatedAt,
		model.UpdatedAt,
	)

	return err
}

func (r *TrainingDatasetRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.TrainingDataset, error) {
	query := `SELECT
		id, project_id, version, generate_model, generate_model_runner,
		generate_gpu_info_card, generate_gpu_info_total_gb, generate_gpu_info_cuda_version,
		input_field, output_field, total_generation_time_seconds,
		generate_prompt_history_ids_json, generate_prompt_id, corpus_id,
		language_iso, status, field_names_json, data_json, created_at, updated_at
	FROM training_datasets WHERE id = $1`

	var model TrainingDatasetRepositoryModel
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Version,
		&model.GenerateModel,
		&model.GenerateModelRunner,
		&model.GenerateGPUInfoCard,
		&model.GenerateGPUInfoTotalGB,
		&model.GenerateGPUInfoCudaVersion,
		&model.InputField,
		&model.OutputField,
		&model.TotalGenerationTimeSeconds,
		&model.GeneratePromptHistoryIDsJSON,
		&model.GeneratePromptID,
		&model.CorpusID,
		&model.LanguageISO,
		&model.Status,
		&model.FieldNamesJSON,
		&model.DataJSON,
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

func (r *TrainingDatasetRepositoryImpl) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TrainingDataset, error) {
	query := `SELECT
		id, project_id, version, generate_model, generate_model_runner,
		generate_gpu_info_card, generate_gpu_info_total_gb, generate_gpu_info_cuda_version,
		input_field, output_field, total_generation_time_seconds,
		generate_prompt_history_ids_json, generate_prompt_id, corpus_id,
		language_iso, status, field_names_json, data_json, created_at, updated_at
	FROM training_datasets WHERE project_id = $1 ORDER BY version DESC`

	rows, err := r.Db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainingDatasets []*entities.TrainingDataset
	for rows.Next() {
		var model TrainingDatasetRepositoryModel
		err := rows.Scan(
			&model.ID,
			&model.ProjectID,
			&model.Version,
			&model.GenerateModel,
			&model.GenerateModelRunner,
			&model.GenerateGPUInfoCard,
			&model.GenerateGPUInfoTotalGB,
			&model.GenerateGPUInfoCudaVersion,
			&model.InputField,
			&model.OutputField,
			&model.TotalGenerationTimeSeconds,
			&model.GeneratePromptHistoryIDsJSON,
			&model.GeneratePromptID,
			&model.CorpusID,
			&model.LanguageISO,
			&model.Status,
			&model.FieldNamesJSON,
			&model.DataJSON,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		entity, err := model.ToEntity()
		if err != nil {
			return nil, err
		}

		trainingDatasets = append(trainingDatasets, entity)
	}

	return trainingDatasets, nil
}

func (r *TrainingDatasetRepositoryImpl) GetLatestByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.TrainingDataset, error) {
	query := `SELECT
		id, project_id, version, generate_model, generate_model_runner,
		generate_gpu_info_card, generate_gpu_info_total_gb, generate_gpu_info_cuda_version,
		input_field, output_field, total_generation_time_seconds,
		generate_prompt_history_ids_json, generate_prompt_id, corpus_id,
		language_iso, status, field_names_json, data_json, created_at, updated_at
	FROM training_datasets WHERE project_id = $1 ORDER BY version DESC LIMIT 1`

	var model TrainingDatasetRepositoryModel
	err := r.Db.QueryRowContext(ctx, query, projectID).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Version,
		&model.GenerateModel,
		&model.GenerateModelRunner,
		&model.GenerateGPUInfoCard,
		&model.GenerateGPUInfoTotalGB,
		&model.GenerateGPUInfoCudaVersion,
		&model.InputField,
		&model.OutputField,
		&model.TotalGenerationTimeSeconds,
		&model.GeneratePromptHistoryIDsJSON,
		&model.GeneratePromptID,
		&model.CorpusID,
		&model.LanguageISO,
		&model.Status,
		&model.FieldNamesJSON,
		&model.DataJSON,
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

func (r *TrainingDatasetRepositoryImpl) Update(ctx context.Context, trainingDataset *entities.TrainingDataset) error {
	query := `UPDATE training_datasets SET
		generate_model = $3, generate_model_runner = $4, generate_gpu_info_card = $5,
		generate_gpu_info_total_gb = $6, generate_gpu_info_cuda_version = $7,
		input_field = $8, output_field = $9, total_generation_time_seconds = $10,
		generate_prompt_history_ids_json = $11, generate_prompt_id = $12, corpus_id = $13,
		language_iso = $14, status = $15, field_names_json = $16, data_json = $17, updated_at = $18
	WHERE id = $1 AND project_id = $2`

	trainingDataset.UpdatedAt = time.Now()

	model, err := FromTrainingDatasetEntity(trainingDataset)
	if err != nil {
		return err
	}

	_, err = r.Db.ExecContext(ctx, query,
		model.ID,
		model.ProjectID,
		model.GenerateModel,
		model.GenerateModelRunner,
		model.GenerateGPUInfoCard,
		model.GenerateGPUInfoTotalGB,
		model.GenerateGPUInfoCudaVersion,
		model.InputField,
		model.OutputField,
		model.TotalGenerationTimeSeconds,
		model.GeneratePromptHistoryIDsJSON,
		model.GeneratePromptID,
		model.CorpusID,
		model.LanguageISO,
		model.Status,
		model.FieldNamesJSON,
		model.DataJSON,
		model.UpdatedAt,
	)

	return err
}

func (r *TrainingDatasetRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM training_datasets WHERE id = $1`
	_, err := r.Db.ExecContext(ctx, query, id)
	return err
}