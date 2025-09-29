package use_cases

import (
	"context"
	"fmt"
	"io"

	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
	"ai-platform/internal/application/port/out/persistence"
)

type DownloadModelUseCaseImpl struct {
	FinetuneRepository   persistence.FinetuneRepository
	DownloadModelClient  clients.DownloadModelClient
}

func (u *DownloadModelUseCaseImpl) DownloadModel(ctx context.Context, command in.DownloadModelCommand) (io.ReadCloser, int64, string, error) {
	// Verify that the finetune exists and belongs to the project
	finetune, err := u.FinetuneRepository.GetByID(ctx, command.FinetuneID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to get finetune: %w", err)
	}

	if finetune.ProjectID != command.ProjectID {
		return nil, 0, "", fmt.Errorf("finetune %s does not belong to project %s", command.FinetuneID, command.ProjectID)
	}

	// Get model name from finetune entity
	modelName := finetune.ModelName
	if modelName == "" {
		return nil, 0, "", fmt.Errorf("finetune %s has no model name", command.FinetuneID)
	}

	// Download the model from S3
	reader, contentLength, err := u.DownloadModelClient.DownloadModel(ctx, command.FinetuneID, modelName)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to download model: %w", err)
	}

	// Generate filename for download
	filename := fmt.Sprintf("%s.gguf", modelName)

	return reader, contentLength, filename, nil
}