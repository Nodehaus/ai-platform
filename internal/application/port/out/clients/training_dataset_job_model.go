package clients

type TrainingDatasetJobModel struct {
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