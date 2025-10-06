package clients

import (
	"context"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type TrainingDatasetResult struct {
	TotalGenerationTimeSeconds float64                     `json:"total_generation_time_seconds"`
	TokensIn                   int                         `json:"tokens_in"`
	TokensOut                  int                         `json:"tokens_out"`
	TrainingDataItems          []entities.TrainingDataItem `json:"training_data_items"`
}

type TrainingDatasetResultsClient interface {
	GetTrainingDatasetResults(ctx context.Context, trainingDatasetID uuid.UUID, fieldNames []string) (*TrainingDatasetResult, error)
}