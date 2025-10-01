package services

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type TrainingDatasetService struct{}


func (s *TrainingDatasetService) CreateTrainingDataset(
	projectID uuid.UUID,
	corpusID uuid.UUID,
	promptID uuid.UUID,
	inputField string,
	outputField string,
	languageISO string,
	fieldNames []string,
	generateExamplesNumber int,
) (*entities.TrainingDataset, error) {
	if err := s.ValidateCreateTrainingDatasetRequest(inputField, outputField, languageISO, fieldNames); err != nil {
		return nil, err
	}

	return &entities.TrainingDataset{
		ID:                       uuid.New(),
		ProjectID:                projectID,
		Version:                  1,
		InputField:               inputField,
		OutputField:              outputField,
		GeneratePromptHistoryIDs: []uuid.UUID{},
		GeneratePromptID:         promptID,
		CorpusID:                 corpusID,
		LanguageISO:              languageISO,
		Status:                   entities.TrainingDatasetStatusPlanning,
		FieldNames:               fieldNames,
		GenerateExamplesNumber:   generateExamplesNumber,
		Data:                     []entities.TrainingDataItem{},
	}, nil
}

func (s *TrainingDatasetService) ValidateCreateTrainingDatasetRequest(
	inputField string,
	outputField string,
	languageISO string,
	fieldNames []string,
) error {
	if inputField == "" {
		return errors.New("input_field is required")
	}
	if outputField == "" {
		return errors.New("output_field is required")
	}
	if languageISO == "" {
		return errors.New("language_iso is required")
	}
	if len(languageISO) != 3 {
		return errors.New("language_iso must be a 3-letter ISO code")
	}
	if len(fieldNames) == 0 {
		return errors.New("field_names cannot be empty")
	}

	// Validate that input_field and output_field are in field_names
	inputFieldFound := false
	outputFieldFound := false
	for _, field := range fieldNames {
		if field == inputField {
			inputFieldFound = true
		}
		if field == outputField {
			outputFieldFound = true
		}
	}
	if !inputFieldFound {
		return errors.New("input_field must be present in field_names")
	}
	if !outputFieldFound {
		return errors.New("output_field must be present in field_names")
	}

	return nil
}

func (s *TrainingDatasetService) ValidateCorpusName(corpusName string) error {
	if corpusName == "" {
		return errors.New("corpus_name is required")
	}
	return nil
}

func (s *TrainingDatasetService) ValidateGeneratePrompt(generatePrompt string) error {
	if generatePrompt == "" {
		return errors.New("generate_prompt is required")
	}
	return nil
}

func (s *TrainingDatasetService) GetNextVersion(projectID uuid.UUID, getLatest func(uuid.UUID) (*entities.TrainingDataset, error)) (int, error) {
	latest, err := getLatest(projectID)
	if err != nil {
		return 0, err
	}
	if latest == nil {
		return 1, nil
	}
	return latest.Version + 1, nil
}

func (s *TrainingDatasetService) SelectTrainingDataSubset(
	trainingData []entities.TrainingDataItem,
	numberExamples *int,
	selectRandom bool,
) []entities.TrainingDataItem {
	// Filter out deleted items
	var availableData []entities.TrainingDataItem
	for _, item := range trainingData {
		if !item.Deleted {
			availableData = append(availableData, item)
		}
	}

	// If no number specified or number is greater than available, return all
	if numberExamples == nil || *numberExamples >= len(availableData) {
		return availableData
	}

	// If number is less than or equal to 0, return empty slice
	if *numberExamples <= 0 {
		return []entities.TrainingDataItem{}
	}

	// Select subset
	if selectRandom {
		// Create a copy to avoid modifying the original slice
		dataCopy := make([]entities.TrainingDataItem, len(availableData))
		copy(dataCopy, availableData)

		// Shuffle the data
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(dataCopy), func(i, j int) {
			dataCopy[i], dataCopy[j] = dataCopy[j], dataCopy[i]
		})

		return dataCopy[:*numberExamples]
	} else {
		// Return first N items in order
		return availableData[:*numberExamples]
	}
}

func (s *TrainingDatasetService) ConvertToFinetuneJobData(
	trainingDataItems []entities.TrainingDataItem,
	fieldNames []string,
	inputField string,
) []map[string]interface{} {
	var jobData []map[string]interface{}

	for _, item := range trainingDataItems {
		dataItem := make(map[string]interface{})

		// Map values to field names
		for i, fieldName := range fieldNames {
			if i < len(item.Values) {
				dataItem[fieldName] = item.Values[i]
			}
		}

		// Add extra fields if input field is "source_text"
		if inputField == "source_text" {
			if item.SourceDocument != nil {
				dataItem["source_document"] = *item.SourceDocument
			}
			if item.SourceDocumentStart != nil {
				dataItem["source_document_start"] = *item.SourceDocumentStart
			}
			if item.SourceDocumentEnd != nil {
				dataItem["source_document_end"] = *item.SourceDocumentEnd
			}
		}

		jobData = append(jobData, dataItem)
	}

	return jobData
}

// GenerateCsvFilename creates a filename in the format: dataset_{project_name}_v{version}.csv
// Project name is converted to lowercase, non-alphanumeric and non-space characters are removed,
// and spaces are replaced with underscores
func (s *TrainingDatasetService) GenerateCsvFilename(projectName string, version int) string {
	// Convert to lowercase
	name := strings.ToLower(projectName)

	// Remove all non-alphanumeric, non-space, and non-underscore characters
	reg := regexp.MustCompile(`[^a-z0-9\s_]`)
	name = reg.ReplaceAllString(name, "")

	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")

	// Remove any multiple consecutive underscores
	reg = regexp.MustCompile(`_+`)
	name = reg.ReplaceAllString(name, "_")

	// Trim underscores from start and end
	name = strings.Trim(name, "_")

	return fmt.Sprintf("dataset_%s_v%d.csv", name, version)
}