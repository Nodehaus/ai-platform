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

type UploadNewTrainingDatasetVersionUseCaseImpl struct {
	TrainingDatasetService    *services.TrainingDatasetService
	TrainingDatasetRepository persistence.TrainingDatasetRepository
}

func (uc *UploadNewTrainingDatasetVersionUseCaseImpl) UploadNewTrainingDatasetVersion(command in.UploadNewTrainingDatasetVersionCommand) (*in.UploadNewTrainingDatasetVersionResult, error) {
	// Get the latest training dataset for this project
	latestDataset, err := uc.TrainingDatasetRepository.GetLatestByProjectID(context.Background(), command.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest training dataset: %w", err)
	}

	if latestDataset == nil {
		return nil, fmt.Errorf("no training dataset found for project")
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

	// Use the field names from the CSV header as the new field names
	fieldNames := header

	if len(fieldNames) == 0 {
		return nil, fmt.Errorf("CSV header cannot be empty")
	}

	// Process data rows
	var dataItems []entities.TrainingDataItem
	for rowIndex, record := range records[1:] {
		if len(record) != len(fieldNames) {
			return nil, fmt.Errorf("row %d has %d columns but expected %d columns", rowIndex+2, len(record), len(fieldNames))
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

		dataItems = append(dataItems, dataItem)
	}

	// Create new training dataset with incremented version
	newVersion := latestDataset.Version + 1
	newTrainingDataset := &entities.TrainingDataset{
		ID:                       uuid.New(),
		ProjectID:                command.ProjectID,
		Version:                  newVersion,
		GenerateModel:            latestDataset.GenerateModel,
		GenerateModelRunner:      latestDataset.GenerateModelRunner,
		GenerateGPUInfoCard:      latestDataset.GenerateGPUInfoCard,
		GenerateGPUInfoTotalGB:   latestDataset.GenerateGPUInfoTotalGB,
		GenerateGPUInfoCudaVersion: latestDataset.GenerateGPUInfoCudaVersion,
		InputField:               latestDataset.InputField,
		OutputField:              latestDataset.OutputField,
		TotalGenerationTimeSeconds: nil,
		GeneratePromptHistoryIDs: latestDataset.GeneratePromptHistoryIDs,
		GeneratePromptID:         latestDataset.GeneratePromptID,
		CorpusID:                 latestDataset.CorpusID,
		LanguageISO:              latestDataset.LanguageISO,
		Status:                   entities.TrainingDatasetStatusDone,
		FieldNames:               fieldNames,
		GenerateExamplesNumber:   len(dataItems),
		Data:                     dataItems,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	// Save the new training dataset
	err = uc.TrainingDatasetRepository.Create(context.Background(), newTrainingDataset)
	if err != nil {
		return nil, fmt.Errorf("failed to create new training dataset: %w", err)
	}

	return &in.UploadNewTrainingDatasetVersionResult{
		TrainingDatasetID: newTrainingDataset.ID,
		Version:           newVersion,
		TotalItems:        len(dataItems),
	}, nil
}
