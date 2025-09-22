package web

type CreateTrainingDatasetRequest struct {
	CorpusName            string   `json:"corpus_name" binding:"required"`
	InputField            string   `json:"input_field" binding:"required"`
	OutputField           string   `json:"output_field" binding:"required"`
	LanguageISO           string   `json:"language_iso" binding:"required"`
	FieldNames            []string `json:"field_names" binding:"required"`
	GeneratePrompt        string   `json:"generate_prompt" binding:"required"`
	GenerateExamplesNumber int      `json:"generate_examples_number" binding:"required"`
}