package in

import (
	"ai-platform/internal/application/domain/entities"
	"context"

	"github.com/google/uuid"
)

type GetFinetuneResult struct {
	Finetune     *entities.Finetune
	DeploymentID *uuid.UUID
}

type GetFinetuneUseCase interface {
	GetFinetune(ctx context.Context, command GetFinetuneCommand) (*GetFinetuneResult, error)
}