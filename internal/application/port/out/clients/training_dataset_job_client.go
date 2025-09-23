package clients

import (
	"context"

	"ai-platform/internal/application/domain/entities"
)

type TrainingDatasetJobClient interface {
	SubmitJob(ctx context.Context, job entities.TrainingDatasetJob) error
}