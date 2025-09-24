package in

import "github.com/google/uuid"

type CreateFinetuneCommand struct {
	UserID                           uuid.UUID `json:"user_id"`
	ProjectID                        uuid.UUID `json:"project_id"`
	BaseModelName                    string    `json:"base_model_name"`
	TrainingDatasetID                uuid.UUID `json:"training_dataset_id"`
	TrainingDatasetNumberExamples    *int      `json:"training_dataset_number_examples,omitempty"`
	TrainingDatasetSelectRandom      bool      `json:"training_dataset_select_random"`
}