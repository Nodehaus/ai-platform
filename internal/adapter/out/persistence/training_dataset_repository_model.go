package persistence

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type TrainingDatasetRepositoryModel struct {
	ID                              uuid.UUID `db:"id"`
	ProjectID                       uuid.UUID `db:"project_id"`
	Version                         int       `db:"version"`
	GenerateModel                   *string   `db:"generate_model"`
	GenerateModelRunner             *string   `db:"generate_model_runner"`
	GenerateGPUInfoCard             *string   `db:"generate_gpu_info_card"`
	GenerateGPUInfoTotalGB          *float64  `db:"generate_gpu_info_total_gb"`
	GenerateGPUInfoCudaVersion      *string   `db:"generate_gpu_info_cuda_version"`
	InputField                      string    `db:"input_field"`
	OutputField                     string    `db:"output_field"`
	JSONObjectFieldsJSON            string    `db:"json_object_fields_json"`
	ExpectedOutputSizeChars         int       `db:"expected_output_size_chars"`
	TotalGenerationTimeSeconds      *float64  `db:"total_generation_time_seconds"`
	GeneratePromptHistoryIDsJSON    string     `db:"generate_prompt_history_ids_json"`
	GeneratePromptID                uuid.UUID  `db:"generate_prompt_id"`
	CorpusID                        *uuid.UUID `db:"corpus_id"`
	LanguageISO                     string     `db:"language_iso"`
	Status                          string    `db:"status"`
	FieldNamesJSON                  string    `db:"field_names_json"`
	GenerateExamplesNumber          int       `db:"generate_examples_number"`
	CreatedAt                       time.Time `db:"created_at"`
	UpdatedAt                       time.Time `db:"updated_at"`
}

func (m *TrainingDatasetRepositoryModel) ToEntity() (*entities.TrainingDataset, error) {
	var generatePromptHistoryIDs []uuid.UUID
	if m.GeneratePromptHistoryIDsJSON != "" {
		if err := json.Unmarshal([]byte(m.GeneratePromptHistoryIDsJSON), &generatePromptHistoryIDs); err != nil {
			return nil, err
		}
	}

	var fieldNames []string
	if err := json.Unmarshal([]byte(m.FieldNamesJSON), &fieldNames); err != nil {
		return nil, err
	}

	var jsonObjectFields map[string]string
	if m.JSONObjectFieldsJSON != "" {
		if err := json.Unmarshal([]byte(m.JSONObjectFieldsJSON), &jsonObjectFields); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json_object_fields: %w (value: %s)", err, m.JSONObjectFieldsJSON)
		}
	}
	if jsonObjectFields == nil {
		jsonObjectFields = make(map[string]string)
	}

	return &entities.TrainingDataset{
		ID:                              m.ID,
		ProjectID:                       m.ProjectID,
		Version:                         m.Version,
		GenerateModel:                   m.GenerateModel,
		GenerateModelRunner:             m.GenerateModelRunner,
		GenerateGPUInfoCard:             m.GenerateGPUInfoCard,
		GenerateGPUInfoTotalGB:          m.GenerateGPUInfoTotalGB,
		GenerateGPUInfoCudaVersion:      m.GenerateGPUInfoCudaVersion,
		InputField:                      m.InputField,
		OutputField:                     m.OutputField,
		JSONObjectFields:                jsonObjectFields,
		ExpectedOutputSizeChars:         m.ExpectedOutputSizeChars,
		TotalGenerationTimeSeconds:      m.TotalGenerationTimeSeconds,
		GeneratePromptHistoryIDs:        generatePromptHistoryIDs,
		GeneratePromptID:                m.GeneratePromptID,
		CorpusID:                        m.CorpusID,
		LanguageISO:                     m.LanguageISO,
		Status:                          entities.TrainingDatasetStatus(m.Status),
		FieldNames:                      fieldNames,
		GenerateExamplesNumber:          m.GenerateExamplesNumber,
		Data:                            []entities.TrainingDataItem{}, // Will be populated separately
		CreatedAt:                       m.CreatedAt,
		UpdatedAt:                       m.UpdatedAt,
	}, nil
}

func FromTrainingDatasetEntity(td *entities.TrainingDataset) (*TrainingDatasetRepositoryModel, error) {
	generatePromptHistoryIDsJSON, err := json.Marshal(td.GeneratePromptHistoryIDs)
	if err != nil {
		return nil, err
	}

	fieldNamesJSON, err := json.Marshal(td.FieldNames)
	if err != nil {
		return nil, err
	}

	jsonObjectFieldsJSON, err := json.Marshal(td.JSONObjectFields)
	if err != nil {
		return nil, err
	}

	return &TrainingDatasetRepositoryModel{
		ID:                              td.ID,
		ProjectID:                       td.ProjectID,
		Version:                         td.Version,
		GenerateModel:                   td.GenerateModel,
		GenerateModelRunner:             td.GenerateModelRunner,
		GenerateGPUInfoCard:             td.GenerateGPUInfoCard,
		GenerateGPUInfoTotalGB:          td.GenerateGPUInfoTotalGB,
		GenerateGPUInfoCudaVersion:      td.GenerateGPUInfoCudaVersion,
		InputField:                      td.InputField,
		OutputField:                     td.OutputField,
		JSONObjectFieldsJSON:            string(jsonObjectFieldsJSON),
		ExpectedOutputSizeChars:         td.ExpectedOutputSizeChars,
		TotalGenerationTimeSeconds:      td.TotalGenerationTimeSeconds,
		GeneratePromptHistoryIDsJSON:    string(generatePromptHistoryIDsJSON),
		GeneratePromptID:                td.GeneratePromptID,
		CorpusID:                        td.CorpusID,
		LanguageISO:                     td.LanguageISO,
		Status:                          string(td.Status),
		FieldNamesJSON:                  string(fieldNamesJSON),
		GenerateExamplesNumber:          td.GenerateExamplesNumber,
		CreatedAt:                       td.CreatedAt,
		UpdatedAt:                       td.UpdatedAt,
	}, nil
}