package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type UpdateFinetuneStatusUseCaseImpl struct {
	FinetuneRepository persistence.FinetuneRepository
}

func (uc *UpdateFinetuneStatusUseCaseImpl) Execute(ctx context.Context, command in.UpdateFinetuneStatusCommand) error {
	// Validate that the finetune exists
	finetune, err := uc.FinetuneRepository.GetByID(ctx, command.FinetuneID)
	if err != nil {
		return fmt.Errorf("failed to get finetune: %w", err)
	}
	if finetune == nil {
		return fmt.Errorf("finetune not found")
	}

	// Validate status transition
	if !isValidFinetuneStatusTransition(finetune.Status, command.Status) {
		return fmt.Errorf("invalid status transition from %s to %s", finetune.Status, command.Status)
	}

	// Update the status
	err = uc.FinetuneRepository.UpdateStatus(ctx, command.FinetuneID, command.Status)
	if err != nil {
		return fmt.Errorf("failed to update finetune status: %w", err)
	}

	return nil
}

func isValidFinetuneStatusTransition(current, new entities.FinetuneStatus) bool {
	// Only allow updates to RUNNING, FAILED, and DONE statuses
	allowedStatuses := map[entities.FinetuneStatus]bool{
		entities.FinetuneStatusRunning: true,
		entities.FinetuneStatusFailed:  true,
		entities.FinetuneStatusDone:    true,
	}

	return allowedStatuses[new]
}