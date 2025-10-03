package web

type CreateTrainingDatasetRequest struct {
	CorpusName              string            `json:"corpus_name"`
	InputField              string            `json:"input_field" binding:"required"`
	OutputField             string            `json:"output_field" binding:"required"`
	JSONObjectFields        map[string]string `json:"json_object_fields" binding:"required"`
	ExpectedOutputSizeChars int               `json:"expected_output_size_chars" binding:"required"`
	LanguageISO             string            `json:"language_iso" binding:"required"`
	FieldNames              []string          `json:"field_names" binding:"required"`
	GeneratePrompt          string            `json:"generate_prompt" binding:"required"`
	GenerateExamplesNumber  int               `json:"generate_examples_number" binding:"required"`
	GenerateModel           string            `json:"generate_model"`
	GenerateModelRunner     string            `json:"generate_model_runner"`
}