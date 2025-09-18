package services

import (
	"ai-platform/internal/application/domain/entities"
	"errors"
	"github.com/google/uuid"
)

type ProjectService struct{}


func (s *ProjectService) CreateProject(name string, ownerID uuid.UUID) *entities.Project {
	return &entities.Project{
		ID:              uuid.New(),
		Name:            name,
		OwnerID:         ownerID,
		TrainingDataset: nil,
		Finetune:        nil,
		Status:          entities.ProjectStatusActive,
	}
}

func (s *ProjectService) ValidateProjectName(name string) error {
	if name == "" {
		return errors.New("project name cannot be empty")
	}
	if len(name) > 100 {
		return errors.New("project name cannot exceed 100 characters")
	}
	return nil
}

func (s *ProjectService) ValidateProjectNameUniqueness(name string, ownerID uuid.UUID, exists func(string, uuid.UUID) (bool, error)) error {
	nameExists, err := exists(name, ownerID)
	if err != nil {
		return err
	}
	if nameExists {
		return errors.New("project name already exists")
	}
	return nil
}