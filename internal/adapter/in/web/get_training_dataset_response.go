package web

import (
	"ai-platform/internal/application/domain/entities"
)

type GetTrainingDatasetResponse struct {
	Version                int      `json:"version"`
	GeneratePrompt         string   `json:"generate_prompt"`
	InputField             string   `json:"input_field"`
	OutputField            string   `json:"output_field"`
	GenerateExamplesNumber int      `json:"generate_examples_number"`
	CorpusName             string   `json:"corpus_name"`
	LanguageISO            string   `json:"language_iso"`
	Status                 string   `json:"status"`
	FieldNames             []string `json:"field_names"`
	DataItemsSample        [][]string `json:"data_items_sample"`
}

func ToGetTrainingDatasetResponse(td *entities.TrainingDataset, prompt string, corpusName string) *GetTrainingDatasetResponse {
	response := &GetTrainingDatasetResponse{
		Version:                td.Version,
		GeneratePrompt:         prompt,
		InputField:             td.InputField,
		OutputField:            td.OutputField,
		GenerateExamplesNumber: td.GenerateExamplesNumber,
		CorpusName:             corpusName,
		LanguageISO:            td.LanguageISO,
		Status:                 string(td.Status),
		FieldNames:             td.FieldNames,
		DataItemsSample:        [][]string{},
	}

	// Only include sample data if status is DONE
	if td.Status == entities.TrainingDatasetStatusDone {
		// Get first 10 data items that are not deleted and don't correct other items
		sampleCount := 0
		for _, item := range td.Data {
			if sampleCount >= 10 {
				break
			}
			if !item.Deleted && item.CorrectsID == nil {
				response.DataItemsSample = append(response.DataItemsSample, item.Values)
				sampleCount++
			}
		}
	}

	return response
}