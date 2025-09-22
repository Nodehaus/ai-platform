package clients

import (
	"context"
)

type TrainingDatasetJobClient interface {
	SubmitJob(ctx context.Context, job TrainingDatasetJobModel) error
}