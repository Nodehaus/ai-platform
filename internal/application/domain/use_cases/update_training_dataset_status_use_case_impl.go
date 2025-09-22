package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
	"ai-platform/internal/application/port/out/persistence"
)

type UpdateTrainingDatasetStatusUseCaseImpl struct {
	TrainingDatasetRepository     persistence.TrainingDatasetRepository
	TrainingDatasetResultsClient clients.TrainingDatasetResultsClient
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

	// If setting status to DONE, fetch results from S3 and update training data
	if command.Status == entities.TrainingDatasetStatusDone {
		err = uc.processCompletedTrainingDataset(ctx, trainingDataset)
		if err != nil {
			return fmt.Errorf("failed to process completed training dataset: %w", err)
		}
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

func (uc *UpdateTrainingDatasetStatusUseCaseImpl) processCompletedTrainingDataset(ctx context.Context, trainingDataset *entities.TrainingDataset) error {
	// Fetch results from S3
	results, err := uc.TrainingDatasetResultsClient.GetTrainingDatasetResults(ctx, trainingDataset.ID, trainingDataset.FieldNames)
	if err != nil {
		return fmt.Errorf("failed to get training dataset results from S3: %w", err)
	}

	// Update the training dataset with the results
	trainingDataset.TotalGenerationTimeSeconds = &results.TotalGenerationTimeSeconds
	trainingDataset.Data = results.TrainingDataItems

	// Save the updated training dataset
	err = uc.TrainingDatasetRepository.Update(ctx, trainingDataset)
	if err != nil {
		return fmt.Errorf("failed to update training dataset with results: %w", err)
	}

	return nil
}