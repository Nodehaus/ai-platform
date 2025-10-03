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
		input_field, output_field, json_object_fields_json, expected_output_size_chars,
		total_generation_time_seconds,
		generate_prompt_history_ids_json, generate_prompt_id, corpus_id,
		language_iso, status, field_names_json, generate_examples_number, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)`

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
		model.JSONObjectFieldsJSON,
		model.ExpectedOutputSizeChars,
		model.TotalGenerationTimeSeconds,
		model.GeneratePromptHistoryIDsJSON,
		model.GeneratePromptID,
		model.CorpusID,
		model.LanguageISO,
		model.Status,
		model.FieldNamesJSON,
		model.GenerateExamplesNumber,
		model.CreatedAt,
		model.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Create training data items
	return r.createTrainingDataItems(ctx, trainingDataset)
}

func (r *TrainingDatasetRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.TrainingDataset, error) {
	query := `SELECT
		id, project_id, version, generate_model, generate_model_runner,
		generate_gpu_info_card, generate_gpu_info_total_gb, generate_gpu_info_cuda_version,
		input_field, output_field, json_object_fields_json, expected_output_size_chars,
		total_generation_time_seconds,
		generate_prompt_history_ids_json, generate_prompt_id, corpus_id,
		language_iso, status, field_names_json, generate_examples_number, created_at, updated_at
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
		&model.JSONObjectFieldsJSON,
		&model.ExpectedOutputSizeChars,
		&model.TotalGenerationTimeSeconds,
		&model.GeneratePromptHistoryIDsJSON,
		&model.GeneratePromptID,
		&model.CorpusID,
		&model.LanguageISO,
		&model.Status,
		&model.FieldNamesJSON,
		&model.GenerateExamplesNumber,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	entity, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	// Load training data items
	entity.Data, err = r.getTrainingDataItemsByDatasetID(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *TrainingDatasetRepositoryImpl) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TrainingDataset, error) {
	query := `SELECT
		id, project_id, version, generate_model, generate_model_runner,
		generate_gpu_info_card, generate_gpu_info_total_gb, generate_gpu_info_cuda_version,
		input_field, output_field, json_object_fields_json, expected_output_size_chars,
		total_generation_time_seconds,
		generate_prompt_history_ids_json, generate_prompt_id, corpus_id,
		language_iso, status, field_names_json, generate_examples_number, created_at, updated_at
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
			&model.JSONObjectFieldsJSON,
			&model.ExpectedOutputSizeChars,
			&model.TotalGenerationTimeSeconds,
			&model.GeneratePromptHistoryIDsJSON,
			&model.GeneratePromptID,
			&model.CorpusID,
			&model.LanguageISO,
			&model.Status,
			&model.FieldNamesJSON,
			&model.GenerateExamplesNumber,
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

		// Load training data items for each dataset
		entity.Data, err = r.getTrainingDataItemsByDatasetID(ctx, model.ID)
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
		input_field, output_field, json_object_fields_json, expected_output_size_chars,
		total_generation_time_seconds,
		generate_prompt_history_ids_json, generate_prompt_id, corpus_id,
		language_iso, status, field_names_json, generate_examples_number, created_at, updated_at
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
		&model.JSONObjectFieldsJSON,
		&model.ExpectedOutputSizeChars,
		&model.TotalGenerationTimeSeconds,
		&model.GeneratePromptHistoryIDsJSON,
		&model.GeneratePromptID,
		&model.CorpusID,
		&model.LanguageISO,
		&model.Status,
		&model.FieldNamesJSON,
		&model.GenerateExamplesNumber,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	entity, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	// Load training data items
	entity.Data, err = r.getTrainingDataItemsByDatasetID(ctx, model.ID)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *TrainingDatasetRepositoryImpl) Update(ctx context.Context, trainingDataset *entities.TrainingDataset) error {
	query := `UPDATE training_datasets SET
		generate_model = $3, generate_model_runner = $4, generate_gpu_info_card = $5,
		generate_gpu_info_total_gb = $6, generate_gpu_info_cuda_version = $7,
		input_field = $8, output_field = $9, json_object_fields_json = $10, expected_output_size_chars = $11,
		total_generation_time_seconds = $12,
		generate_prompt_history_ids_json = $13, generate_prompt_id = $14, corpus_id = $15,
		language_iso = $16, status = $17, field_names_json = $18, generate_examples_number = $19, updated_at = $20
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
		model.JSONObjectFieldsJSON,
		model.ExpectedOutputSizeChars,
		model.TotalGenerationTimeSeconds,
		model.GeneratePromptHistoryIDsJSON,
		model.GeneratePromptID,
		model.CorpusID,
		model.LanguageISO,
		model.Status,
		model.FieldNamesJSON,
		model.GenerateExamplesNumber,
		model.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Update training data items
	return r.updateTrainingDataItems(ctx, trainingDataset)
}

func (r *TrainingDatasetRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.TrainingDatasetStatus) error {
	query := `UPDATE training_datasets SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.Db.ExecContext(ctx, query, id, status, time.Now())
	return err
}

func (r *TrainingDatasetRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Training data items will be deleted automatically due to CASCADE
	query := `DELETE FROM training_datasets WHERE id = $1`
	_, err := r.Db.ExecContext(ctx, query, id)
	return err
}

// Helper methods for managing TrainingDataItems

func (r *TrainingDatasetRepositoryImpl) createTrainingDataItems(ctx context.Context, trainingDataset *entities.TrainingDataset) error {
	if len(trainingDataset.Data) == 0 {
		return nil
	}

	query := `INSERT INTO training_data_items (
		id, training_dataset_id, values_json, corrects_id, source_document,
		source_document_start, source_document_end, generation_time_seconds, deleted, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	for _, item := range trainingDataset.Data {
		if item.ID == uuid.Nil {
			item.ID = uuid.New()
		}
		item.CreatedAt = trainingDataset.CreatedAt
		item.UpdatedAt = trainingDataset.UpdatedAt

		model, err := FromTrainingDataItemEntity(&item, trainingDataset.ID)
		if err != nil {
			return err
		}

		_, err = r.Db.ExecContext(ctx, query,
			model.ID,
			model.TrainingDatasetID,
			model.ValuesJSON,
			model.CorrectsID,
			model.SourceDocument,
			model.SourceDocumentStart,
			model.SourceDocumentEnd,
			model.GenerationTimeSeconds,
			model.Deleted,
			model.CreatedAt,
			model.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TrainingDatasetRepositoryImpl) getTrainingDataItemsByDatasetID(ctx context.Context, datasetID uuid.UUID) ([]entities.TrainingDataItem, error) {
	query := `SELECT
		id, training_dataset_id, values_json, corrects_id, source_document,
		source_document_start, source_document_end, generation_time_seconds, deleted, created_at, updated_at
	FROM training_data_items WHERE training_dataset_id = $1 AND deleted = false ORDER BY created_at`

	rows, err := r.Db.QueryContext(ctx, query, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entities.TrainingDataItem
	for rows.Next() {
		var model TrainingDataItemRepositoryModel
		err := rows.Scan(
			&model.ID,
			&model.TrainingDatasetID,
			&model.ValuesJSON,
			&model.CorrectsID,
			&model.SourceDocument,
			&model.SourceDocumentStart,
			&model.SourceDocumentEnd,
			&model.GenerationTimeSeconds,
			&model.Deleted,
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

		items = append(items, *entity)
	}

	return items, nil
}

func (r *TrainingDatasetRepositoryImpl) updateTrainingDataItems(ctx context.Context, trainingDataset *entities.TrainingDataset) error {
	// Delete existing items (we'll recreate them)
	deleteQuery := `DELETE FROM training_data_items WHERE training_dataset_id = $1`
	_, err := r.Db.ExecContext(ctx, deleteQuery, trainingDataset.ID)
	if err != nil {
		return err
	}

	// Create new items
	return r.createTrainingDataItems(ctx, trainingDataset)
}