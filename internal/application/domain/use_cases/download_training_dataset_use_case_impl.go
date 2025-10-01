package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type DownloadTrainingDatasetUseCaseImpl struct {
	TrainingDatasetService    *services.TrainingDatasetService
	TrainingDatasetRepository persistence.TrainingDatasetRepository
	ProjectRepository         persistence.ProjectRepository
}

func (uc *DownloadTrainingDatasetUseCaseImpl) DownloadTrainingDataset(command in.DownloadTrainingDatasetCommand) (*in.DownloadTrainingDatasetResult, error) {
	// Get the training dataset with all data items
	trainingDataset, err := uc.TrainingDatasetRepository.GetByID(context.Background(), command.TrainingDatasetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get training dataset: %w", err)
	}

	if trainingDataset == nil {
		return nil, fmt.Errorf("training dataset not found")
	}

	// Verify the training dataset belongs to the specified project
	if trainingDataset.ProjectID != command.ProjectID {
		return nil, fmt.Errorf("training dataset does not belong to the specified project")
	}

	// Get the project to generate filename
	project, err := uc.ProjectRepository.GetByID(command.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	// Convert data items to CSV format, filtering out deleted items
	var data [][]string
	for _, item := range trainingDataset.Data {
		if !item.Deleted {
			data = append(data, item.Values)
		}
	}

	// Generate filename: dataset_{project_name}_v{version}.csv
	filename := uc.TrainingDatasetService.GenerateCsvFilename(project.Name, trainingDataset.Version)

	return &in.DownloadTrainingDatasetResult{
		FieldNames: trainingDataset.FieldNames,
		Data:       data,
		Filename:   filename,
	}, nil
}
