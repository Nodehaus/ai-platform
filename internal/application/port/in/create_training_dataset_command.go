package in

import "github.com/google/uuid"

type CreateTrainingDatasetCommand struct {
	ProjectID      uuid.UUID `json:"project_id"`
	CorpusName     string    `json:"corpus_name"`
	InputField     string    `json:"input_field"`
	OutputField    string    `json:"output_field"`
	LanguageISO    string    `json:"language_iso"`
	FieldNames     []string  `json:"field_names"`
	GeneratePrompt string    `json:"generate_prompt"`
}