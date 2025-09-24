package entities

import (
	"time"

	"github.com/google/uuid"
)

type FinetuneStatus string

const (
	FinetuneStatusPlanning FinetuneStatus = "PLANNING"
	FinetuneStatusRunning  FinetuneStatus = "RUNNING"
	FinetuneStatusAborted  FinetuneStatus = "ABORTED"
	FinetuneStatusFailed   FinetuneStatus = "FAILED"
	FinetuneStatusDone     FinetuneStatus = "DONE"
	FinetuneStatusDeleted  FinetuneStatus = "DELETED"
)

type Finetune struct {
	ID                               uuid.UUID         `json:"id"`
	ProjectID                        uuid.UUID         `json:"project_id"`
	Version                          int               `json:"version"`
	ModelName                        string            `json:"model_name"`
	BaseModelName                    string            `json:"base_model_name"`
	ModelSizeGB                      *int              `json:"model_size_gb,omitempty"`
	ModelSizeParameter               *int              `json:"model_size_parameter,omitempty"`
	ModelDtype                       *string           `json:"model_dtype,omitempty"`
	ModelQuantization                *string           `json:"model_quantization,omitempty"`
	InferenceSamples                 []InferenceSample `json:"inference_samples"`
	TrainingDatasetID                uuid.UUID         `json:"training_dataset_id"`
	TrainingDatasetNumberExamples    *int              `json:"training_dataset_number_examples,omitempty"`
	TrainingDatasetSelectRandom      bool              `json:"training_dataset_select_random"`
	TrainingTimeSeconds              *float64          `json:"training_time_seconds,omitempty"`
	Status                           FinetuneStatus    `json:"status"`
	CreatedAt                        time.Time         `json:"created_at"`
	UpdatedAt                        time.Time         `json:"updated_at"`
}

type InferenceSample struct {
	AtStep int                         `json:"at_step"`
	Items  []InferenceSampleItem      `json:"items"`
}

type InferenceSampleItem struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}