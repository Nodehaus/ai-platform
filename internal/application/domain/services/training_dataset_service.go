package services

import (
	"errors"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type TrainingDatasetService struct{}

func NewTrainingDatasetService() *TrainingDatasetService {
	return &TrainingDatasetService{}
}

func (s *TrainingDatasetService) CreateTrainingDataset(
	projectID uuid.UUID,
	corpusID uuid.UUID,
	promptID uuid.UUID,
	inputField string,
	outputField string,
	languageISO string,
	fieldNames []string,
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