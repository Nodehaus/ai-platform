package clients

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type DownloadModelClient interface {
	DownloadModel(ctx context.Context, finetuneID uuid.UUID, modelName string) (io.ReadCloser, int64, error)
}