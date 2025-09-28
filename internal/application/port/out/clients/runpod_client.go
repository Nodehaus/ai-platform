package clients

import (
	"context"
)

type RunpodClient interface {
	StartFinetuneJob(ctx context.Context, s3Key string, documentsS3Path string, baseModelName string, modelName string, finetuneID string) error
}