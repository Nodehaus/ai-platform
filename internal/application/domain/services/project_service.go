package services

import (
	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/out/persistence"
	"context"
	"errors"
	"github.com/google/uuid"
)

type ProjectWithTrainingDataset struct {
	Project           entities.Project
	TrainingDatasetID *uuid.UUID
}

type ProjectService struct {
	ProjectRepository         persistence.ProjectRepository
	TrainingDatasetRepository persistence.TrainingDatasetRepository
}


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

func (s *ProjectService) GetProject(ctx context.Context, projectID uuid.UUID, ownerID uuid.UUID) (*entities.Project, error) {
	project, err := s.ProjectRepository.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, errors.New("project not found")
	}

	if project.OwnerID != ownerID {
		return nil, errors.New("access denied")
	}

	return project, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, ownerID uuid.UUID) ([]ProjectWithTrainingDataset, error) {
	projects, err := s.ProjectRepository.GetActiveByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}

	result := make([]ProjectWithTrainingDataset, len(projects))
	for i, project := range projects {
		var trainingDatasetID *uuid.UUID

		latestTrainingDataset, err := s.TrainingDatasetRepository.GetLatestByProjectID(ctx, project.ID)
		if err == nil && latestTrainingDataset != nil {
			trainingDatasetID = &latestTrainingDataset.ID
		}

		result[i] = ProjectWithTrainingDataset{
			Project:           project,
			TrainingDatasetID: trainingDatasetID,
		}
	}

	return result, nil
}