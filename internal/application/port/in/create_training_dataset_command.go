package in

import "github.com/google/uuid"

type CreateTrainingDatasetCommand struct {
	UserID                 uuid.UUID `json:"user_id"`
	ProjectID              uuid.UUID `json:"project_id"`
	CorpusName             string    `json:"corpus_name"`
	InputField             string    `json:"input_field"`
	OutputField            string    `json:"output_field"`
	LanguageISO            string    `json:"language_iso"`
	FieldNames             []string  `json:"field_names"`
	GeneratePrompt         string    `json:"generate_prompt"`
	GenerateExamplesNumber int       `json:"generate_examples_number"`
	GenerateModel          string    `json:"generate_model"`
	GenerateModelRunner    string    `json:"generate_model_runner"`
}