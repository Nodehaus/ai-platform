package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type UpdateTrainingDatasetStatusUseCaseImpl struct {
	TrainingDatasetRepository persistence.TrainingDatasetRepository
}

func (uc *UpdateTrainingDatasetStatusUseCaseImpl) Execute(ctx context.Context, command in.UpdateTrainingDatasetStatusCommand) error {
	// Validate that the training dataset exists
	trainingDataset, err := uc.TrainingDatasetRepository.GetByID(ctx, command.TrainingDatasetID)
	if err != nil {
		return fmt.Errorf("failed to get training dataset: %w", err)
	}
	if trainingDataset == nil {
		return fmt.Errorf("training dataset not found")
	}

	// Validate status transition
	if !isValidStatusTransition(trainingDataset.Status, command.Status) {
		return fmt.Errorf("invalid status transition from %s to %s", trainingDataset.Status, command.Status)
	}

	// Update the status
	err = uc.TrainingDatasetRepository.UpdateStatus(ctx, command.TrainingDatasetID, command.Status)
	if err != nil {
		return fmt.Errorf("failed to update training dataset status: %w", err)
	}

	return nil
}

func isValidStatusTransition(current, new entities.TrainingDatasetStatus) bool {
	// Only allow updates to RUNNING, FAILED, and DONE statuses
	allowedStatuses := map[entities.TrainingDatasetStatus]bool{
		entities.TrainingDatasetStatusRunning: true,
		entities.TrainingDatasetStatusFailed:  true,
		entities.TrainingDatasetStatusDone:    true,
	}

	return allowedStatuses[new]
}