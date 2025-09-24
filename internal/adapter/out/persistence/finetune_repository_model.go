package persistence

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type FinetuneRepositoryModel struct {
	ID                               uuid.UUID  `db:"id"`
	ProjectID                        uuid.UUID  `db:"project_id"`
	Version                          int        `db:"version"`
	ModelName                        string     `db:"model_name"`
	BaseModelName                    string     `db:"base_model_name"`
	ModelSizeGB                      *int       `db:"model_size_gb"`
	ModelSizeParameter               *int       `db:"model_size_parameter"`
	ModelDtype                       *string    `db:"model_dtype"`
	ModelQuantization                *string    `db:"model_quantization"`
	InferenceSamplesJSON             string     `db:"inference_samples_json"`
	TrainingDatasetID                uuid.UUID  `db:"training_dataset_id"`
	TrainingDatasetNumberExamples    *int       `db:"training_dataset_number_examples"`
	TrainingDatasetSelectRandom      bool       `db:"training_dataset_select_random"`
	TrainingTimeSeconds              *float64   `db:"training_time_seconds"`
	Status                           string     `db:"status"`
	CreatedAt                        time.Time  `db:"created_at"`
	UpdatedAt                        time.Time  `db:"updated_at"`
}

func (m *FinetuneRepositoryModel) ToEntity() (*entities.Finetune, error) {
	var inferenceSamples []entities.InferenceSample
	if m.InferenceSamplesJSON != "" {
		if err := json.Unmarshal([]byte(m.InferenceSamplesJSON), &inferenceSamples); err != nil {
			return nil, err
		}
	}

	return &entities.Finetune{
		ID:                               m.ID,
		ProjectID:                        m.ProjectID,
		Version:                          m.Version,
		ModelName:                        m.ModelName,
		BaseModelName:                    m.BaseModelName,
		ModelSizeGB:                      m.ModelSizeGB,
		ModelSizeParameter:               m.ModelSizeParameter,
		ModelDtype:                       m.ModelDtype,
		ModelQuantization:                m.ModelQuantization,
		InferenceSamples:                 inferenceSamples,
		TrainingDatasetID:                m.TrainingDatasetID,
		TrainingDatasetNumberExamples:    m.TrainingDatasetNumberExamples,
		TrainingDatasetSelectRandom:      m.TrainingDatasetSelectRandom,
		TrainingTimeSeconds:              m.TrainingTimeSeconds,
		Status:                           entities.FinetuneStatus(m.Status),
		CreatedAt:                        m.CreatedAt,
		UpdatedAt:                        m.UpdatedAt,
	}, nil
}

func FromFinetuneEntity(f *entities.Finetune) (*FinetuneRepositoryModel, error) {
	inferenceSamplesJSON, err := json.Marshal(f.InferenceSamples)
	if err != nil {
		return nil, err
	}

	return &FinetuneRepositoryModel{
		ID:                               f.ID,
		ProjectID:                        f.ProjectID,
		Version:                          f.Version,
		ModelName:                        f.ModelName,
		BaseModelName:                    f.BaseModelName,
		ModelSizeGB:                      f.ModelSizeGB,
		ModelSizeParameter:               f.ModelSizeParameter,
		ModelDtype:                       f.ModelDtype,
		ModelQuantization:                f.ModelQuantization,
		InferenceSamplesJSON:             string(inferenceSamplesJSON),
		TrainingDatasetID:                f.TrainingDatasetID,
		TrainingDatasetNumberExamples:    f.TrainingDatasetNumberExamples,
		TrainingDatasetSelectRandom:      f.TrainingDatasetSelectRandom,
		TrainingTimeSeconds:              f.TrainingTimeSeconds,
		Status:                           string(f.Status),
		CreatedAt:                        f.CreatedAt,
		UpdatedAt:                        f.UpdatedAt,
	}, nil
}