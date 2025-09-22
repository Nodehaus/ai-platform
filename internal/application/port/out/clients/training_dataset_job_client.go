package clients

import (
	"context"

	clientModels "ai-platform/internal/adapter/out/clients"
)

type TrainingDatasetJobClient interface {
	SubmitJob(ctx context.Context, job clientModels.TrainingDatasetJobModel) error
}