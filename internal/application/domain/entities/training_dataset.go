package entities

import (
	"time"

	"github.com/google/uuid"
)

type TrainingDatasetStatus string

const (
	TrainingDatasetStatusPlanning TrainingDatasetStatus = "PLANNING"
	TrainingDatasetStatusRunning  TrainingDatasetStatus = "RUNNING"
	TrainingDatasetStatusAborted  TrainingDatasetStatus = "ABORTED"
	TrainingDatasetStatusFailed   TrainingDatasetStatus = "FAILED"
	TrainingDatasetStatusDone     TrainingDatasetStatus = "DONE"
	TrainingDatasetStatusDeleted  TrainingDatasetStatus = "DELETED"
)

type TrainingDataset struct {
	ID                              uuid.UUID             `json:"id"`
	ProjectID                       uuid.UUID             `json:"project_id"`
	Version                         int                   `json:"version"`
	GenerateModel                   *string               `json:"generate_model,omitempty"`
	GenerateModelRunner             *string               `json:"generate_model_runner,omitempty"`
	GenerateGPUInfoCard             *string               `json:"generate_gpu_info_card,omitempty"`
	GenerateGPUInfoTotalGB          *float64              `json:"generate_gpu_info_total_gb,omitempty"`
	GenerateGPUInfoCudaVersion      *string               `json:"generate_gpu_info_cuda_version,omitempty"`
	InputField                      string                `json:"input_field"`
	OutputField                     string                `json:"output_field"`
	TotalGenerationTimeSeconds      *float64              `json:"total_generation_time_seconds,omitempty"`
	GeneratePromptHistoryIDs        []uuid.UUID           `json:"generate_prompt_history_ids"`
	GeneratePromptID                uuid.UUID             `json:"generate_prompt_id"`
	CorpusID                        uuid.UUID             `json:"corpus_id"`
	LanguageISO                     string                `json:"language_iso"`
	Status                          TrainingDatasetStatus `json:"status"`
	FieldNames                      []string              `json:"field_names"`
	GenerateExamplesNumber          int                   `json:"generate_examples_number"`
	Data                            []TrainingDataItem    `json:"data"`
	CreatedAt                       time.Time             `json:"created_at"`
	UpdatedAt                       time.Time             `json:"updated_at"`
}

type TrainingDataItem struct {
	ID                       uuid.UUID `json:"id"`
	Values                   []string  `json:"values"`
	CorrectsID               *uuid.UUID `json:"corrects_id,omitempty"`
	SourceDocument           *string   `json:"source_document,omitempty"`
	SourceDocumentStart      *string   `json:"source_document_start,omitempty"`
	SourceDocumentEnd        *string   `json:"source_document_end,omitempty"`
	GenerationTimeSeconds    float64   `json:"generation_time_seconds"`
	Deleted                  bool      `json:"deleted"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type TrainingDatasetJob struct {
	CorpusS3Path            string   `json:"corpus_s3_path"`
	CorpusFilesSubset       []string `json:"corpus_files_subset"`
	LanguageISO             string   `json:"language_iso"`
	UserID                  string   `json:"user_id"`
	TrainingDatasetID       string   `json:"training_dataset_id"`
	GeneratePrompt          string   `json:"generate_prompt"`
	GenerateExamplesNumber  int      `json:"generate_examples_number"`
	GenerateModel           string   `json:"generate_model"`
	GenerateModelRunner     string   `json:"generate_model_runner"`
}

