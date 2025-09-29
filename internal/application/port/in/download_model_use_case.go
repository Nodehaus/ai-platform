package in

import (
	"context"
	"io"
)

type DownloadModelUseCase interface {
	DownloadModel(ctx context.Context, command DownloadModelCommand) (io.ReadCloser, int64, string, error)
}