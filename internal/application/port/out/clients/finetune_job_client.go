package clients

import (
	"context"

	"ai-platform/internal/application/domain/entities"
)

type FinetuneJobClient interface {
	SubmitJob(ctx context.Context, job entities.FinetuneJob) (string, error)
}