package in

import (
	"context"

	"ai-platform/internal/application/domain/entities"
)

type CreateFinetuneUseCase interface {
	Execute(ctx context.Context, command CreateFinetuneCommand) (*entities.Finetune, error)
}