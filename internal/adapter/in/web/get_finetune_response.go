package web

import (
	"ai-platform/internal/application/domain/entities"
	"github.com/google/uuid"
)

type GetFinetuneResponse struct {
	ID                               uuid.UUID                    `json:"id"`
	Version                          int                          `json:"version"`
	Status                           entities.FinetuneStatus     `json:"status"`
	BaseModelName                    string                       `json:"base_model_name"`
	ModelName                        string                       `json:"model_name"`
	TrainingDatasetID                uuid.UUID                   `json:"training_dataset_id"`
	TrainingDatasetNumberExamples    *int                         `json:"training_dataset_number_examples"`
	TrainingDatasetSelectRandom      bool                         `json:"training_dataset_select_random"`
	ModelSizeGB                      *int                         `json:"model_size_gb"`
	ModelSizeParameter               *int                         `json:"model_size_parameter"`
	ModelDtype                       *string                      `json:"model_dtype"`
	ModelQuantization                *string                      `json:"model_quantization"`
	InferenceSamples                 []entities.InferenceSample   `json:"inference_samples"`
	TrainingTimeSeconds              *float64                     `json:"training_time_seconds"`
}

func ToGetFinetuneResponse(finetune *entities.Finetune) *GetFinetuneResponse {
	return &GetFinetuneResponse{
		ID:                               finetune.ID,
		Version:                          finetune.Version,
		Status:                           finetune.Status,
		BaseModelName:                    finetune.BaseModelName,
		ModelName:                        finetune.ModelName,
		TrainingDatasetID:                finetune.TrainingDatasetID,
		TrainingDatasetNumberExamples:    finetune.TrainingDatasetNumberExamples,
		TrainingDatasetSelectRandom:      finetune.TrainingDatasetSelectRandom,
		ModelSizeGB:                      finetune.ModelSizeGB,
		ModelSizeParameter:               finetune.ModelSizeParameter,
		ModelDtype:                       finetune.ModelDtype,
		ModelQuantization:                finetune.ModelQuantization,
		InferenceSamples:                 finetune.InferenceSamples,
		TrainingTimeSeconds:              finetune.TrainingTimeSeconds,
	}
}