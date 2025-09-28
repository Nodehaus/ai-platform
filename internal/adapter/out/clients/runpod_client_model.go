package clients

type RunpodClientModel struct {
	S3Bucket               string `json:"s3_bucket"`
	TrainingDatasetS3Path  string `json:"training_dataset_s3_path"`
	DocumentsS3Path        string `json:"documents_s3_path"`
	BaseModelName          string `json:"base_model_name"`
	ModelName              string `json:"model_name"`
}