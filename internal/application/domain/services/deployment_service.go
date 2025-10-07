package services

import (
	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/out/persistence"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
)

type DeploymentService struct {
	DeploymentRepository persistence.DeploymentRepository
	ProjectRepository    persistence.ProjectRepository
	FinetuneRepository   persistence.FinetuneRepository
}

func (s *DeploymentService) CreateDeployment(modelName string, projectID uuid.UUID, finetuneID *uuid.UUID) *entities.Deployment {
	apiKey := s.generateAPIKey()

	return &entities.Deployment{
		ID:         uuid.New(),
		ModelName:  modelName,
		APIKey:     apiKey,
		ProjectID:  projectID,
		FinetuneID: finetuneID,
	}
}

func (s *DeploymentService) ValidateModelName(modelName string) error {
	if modelName == "" {
		return errors.New("model_name is required")
	}
	return nil
}

func (s *DeploymentService) ValidateProjectAccess(projectID uuid.UUID, ownerID uuid.UUID) error {
	project, err := s.ProjectRepository.GetByID(projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}
	if project.OwnerID != ownerID {
		return errors.New("access denied")
	}
	return nil
}

func (s *DeploymentService) ValidateFinetuneExists(ctx context.Context, finetuneID uuid.UUID, projectID uuid.UUID) error {
	finetune, err := s.FinetuneRepository.GetByID(ctx, finetuneID)
	if err != nil {
		return err
	}
	if finetune == nil {
		return errors.New("finetune not found")
	}
	if finetune.ProjectID != projectID {
		return errors.New("finetune does not belong to this project")
	}
	return nil
}

func (s *DeploymentService) ValidateFinetuneNotAlreadyDeployed(finetuneID uuid.UUID) error {
	existingDeployment, err := s.DeploymentRepository.GetByFinetuneID(finetuneID)
	if err != nil {
		return err
	}
	if existingDeployment != nil {
		return errors.New("this model is already deployed")
	}
	return nil
}

func (s *DeploymentService) ValidateModelNameUnique(projectID uuid.UUID, modelName string) error {
	existingDeployment, err := s.DeploymentRepository.GetByProjectIDAndModelName(projectID, modelName)
	if err != nil {
		return err
	}
	if existingDeployment != nil {
		return errors.New("a deployment with this model name already exists in this project")
	}
	return nil
}

func (s *DeploymentService) generateAPIKey() string {
	// Generate a 32-byte random key
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to UUID if random generation fails
		return uuid.New().String()
	}
	// Encode to base64 and prefix with "sk-"
	return base64.URLEncoding.EncodeToString(b)
}
