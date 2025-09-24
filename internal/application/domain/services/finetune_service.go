package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type FinetuneService struct{}

func (s *FinetuneService) ValidateBaseModelName(baseModelName string) error {
	if baseModelName == "" {
		return errors.New("base model name cannot be empty")
	}
	if len(baseModelName) > 100 {
		return errors.New("base model name cannot exceed 100 characters")
	}
	return nil
}

func (s *FinetuneService) GenerateModelName(baseModelName string, projectName string, version int) string {
	projectNameModelString := s.convertProjectNameToModelString(projectName)
	baseModelNameModelString := s.convertBaseModelNameToModelString(baseModelName)
	return fmt.Sprintf("%s_%s_v%d", baseModelNameModelString, projectNameModelString, version)
}

func (s *FinetuneService) convertBaseModelNameToModelString(baseModelName string) string {
	result := strings.ToLower(baseModelName)
	result = regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(result, "_")
	result = regexp.MustCompile(`_+`).ReplaceAllString(result, "_")
	result = strings.Trim(result, "_")

	return result
}

func (s *FinetuneService) convertProjectNameToModelString(projectName string) string {
	result := strings.ToLower(projectName)
	result = regexp.MustCompile(`[^a-z0-9 ]`).ReplaceAllString(result, "")
	result = strings.ReplaceAll(result, " ", "_")
	result = regexp.MustCompile(`_+`).ReplaceAllString(result, "_")
	result = strings.Trim(result, "_")

	if result == "" {
		result = "model"
	}

	return result
}

func (s *FinetuneService) CreateFinetune(projectID, trainingDatasetID uuid.UUID, version int, modelName, baseModelName string, trainingDatasetNumberExamples *int, trainingDatasetSelectRandom bool) *entities.Finetune {
	return &entities.Finetune{
		ID:                               uuid.New(),
		ProjectID:                        projectID,
		Version:                          version,
		ModelName:                        modelName,
		BaseModelName:                    baseModelName,
		TrainingDatasetID:                trainingDatasetID,
		TrainingDatasetNumberExamples:    trainingDatasetNumberExamples,
		TrainingDatasetSelectRandom:      trainingDatasetSelectRandom,
		Status:                           entities.FinetuneStatusPlanning,
		InferenceSamples:                 []entities.InferenceSample{},
	}
}