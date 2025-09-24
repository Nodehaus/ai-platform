package web

import "github.com/google/uuid"

type CreateFinetuneRequest struct {
	BaseModelName                    string    `json:"base_model_name" binding:"required"`
	TrainingDatasetID                string    `json:"training_dataset_id" binding:"required"`
	TrainingDatasetNumberExamples    *int      `json:"training_dataset_number_examples,omitempty"`
	TrainingDatasetSelectRandom      bool      `json:"training_dataset_select_random"`
}

func (r *CreateFinetuneRequest) GetTrainingDatasetID() (uuid.UUID, error) {
	return uuid.Parse(r.TrainingDatasetID)
}