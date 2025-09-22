package clients

// TrainingDatasetResultsFileModel represents the structure of each JSON file from S3
type TrainingDatasetResultsFileModel struct {
	TotalGenerationTime float64                   `json:"total_generation_time"`
	Annotations         []AnnotationModel         `json:"annotations"`
}

// AnnotationModel represents a single annotation in the JSON file
type AnnotationModel struct {
	DocumentID           string                 `json:"document_id"`
	WithinStart          int                    `json:"within_start"`
	WithinEnd            int                    `json:"within_end"`
	InferenceTimeSeconds float64                `json:"inference_time_seconds"`

	// Dynamic fields based on FieldNames - stored as raw map
	Fields map[string]interface{} `json:"-"`
}