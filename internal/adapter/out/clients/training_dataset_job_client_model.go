package clients

type TrainingDatasetJobClientModel struct {
	CorpusS3Path            string   `json:"corpus_s3_path"`
	CorpusFilesSubset       []string `json:"corpus_files_subset"`
	LanguageISO             string   `json:"language_iso"`
	UserID                  string   `json:"user_id"`
	TrainingDatasetID       string   `json:"training_dataset_id"`
	GeneratePrompt          string   `json:"generate_prompt"`
	GenerateExamplesNumber  int      `json:"generate_examples_number"`
	GenerateModel           string   `json:"generate_model"`
	GenerateModelRunner     string   `json:"generate_model_runner"`
	InputField              string   `json:"input_field"`
	OutputField             string   `json:"output_field"`
	JSONObjectFields        string   `json:"json_object_fields"`
	ExpectedOutputSizeChars int      `json:"expected_output_size_chars"`
}