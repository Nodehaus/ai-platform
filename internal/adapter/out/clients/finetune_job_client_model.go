package clients

type FinetuneJobClientModel struct {
	FinetuneID        string                   `json:"finetune_id"`
	TrainingDatasetID string                   `json:"training_dataset_id"`
	InputField        string                   `json:"input_field"`
	OutputField       string                   `json:"output_field"`
	UserID            string                   `json:"user_id"`
	TrainingData      []map[string]interface{} `json:"training_data"`
}