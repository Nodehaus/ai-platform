package use_cases

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type UploadTrainingDatasetUseCaseImpl struct {
	TrainingDatasetService    *services.TrainingDatasetService
	TrainingDatasetRepository persistence.TrainingDatasetRepository
}

func (uc *UploadTrainingDatasetUseCaseImpl) UploadTrainingDataset(command in.UploadTrainingDatasetCommand) (*in.UploadTrainingDatasetResult, error) {
	// Get the training dataset
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

	// Parse CSV data
	reader := csv.NewReader(strings.NewReader(string(command.CsvData)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// First row should be the header
	header := records[0]

	// Validate that header matches field names
	if len(header) != len(trainingDataset.FieldNames) {
		return nil, fmt.Errorf("CSV header has %d columns but expected %d columns", len(header), len(trainingDataset.FieldNames))
	}

	for i, fieldName := range trainingDataset.FieldNames {
		if header[i] != fieldName {
			return nil, fmt.Errorf("CSV header column %d is '%s' but expected '%s'", i+1, header[i], fieldName)
		}
	}

	// Process data rows
	var newDataItems []entities.TrainingDataItem
	for rowIndex, record := range records[1:] {
		if len(record) != len(trainingDataset.FieldNames) {
			return nil, fmt.Errorf("row %d has %d columns but expected %d columns", rowIndex+2, len(record), len(trainingDataset.FieldNames))
		}

		dataItem := entities.TrainingDataItem{
			ID:                    uuid.New(),
			Values:                record,
			CorrectsID:            nil,
			SourceDocument:        nil,
			SourceDocumentStart:   nil,
			SourceDocumentEnd:     nil,
			GenerationTimeSeconds: 0,
			Deleted:               false,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		newDataItems = append(newDataItems, dataItem)
	}

	// Add new items to existing data
	trainingDataset.Data = append(trainingDataset.Data, newDataItems...)

	// Update the GenerateExamplesNumber to reflect the new total
	trainingDataset.GenerateExamplesNumber = len(trainingDataset.Data)

	// Update the training dataset
	err = uc.TrainingDatasetRepository.Update(context.Background(), trainingDataset)
	if err != nil {
		return nil, fmt.Errorf("failed to update training dataset: %w", err)
	}

	return &in.UploadTrainingDatasetResult{
		ItemsAdded: len(newDataItems),
		TotalItems: len(trainingDataset.Data),
	}, nil
}
