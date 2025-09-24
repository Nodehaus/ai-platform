package in

import (
	"ai-platform/internal/application/domain/entities"
	"context"
)

type GetFinetuneResult struct {
	Finetune *entities.Finetune
}

type GetFinetuneUseCase interface {
	GetFinetune(ctx context.Context, command GetFinetuneCommand) (*GetFinetuneResult, error)
}