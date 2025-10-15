package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/out/clients"
	"ai-platform/internal/application/port/out/persistence"
)

type FinetuneCompletionService struct {
	finetuneRepository persistence.FinetuneRepository
	projectRepository  persistence.ProjectRepository
	ollamaLLMClient    clients.OllamaLLMClient
}

func NewFinetuneCompletionService(
	finetuneRepository persistence.FinetuneRepository,
	projectRepository persistence.ProjectRepository,
	ollamaLLMClient clients.OllamaLLMClient,
) *FinetuneCompletionService {
	return &FinetuneCompletionService{
		finetuneRepository: finetuneRepository,
		projectRepository:  projectRepository,
		ollamaLLMClient:    ollamaLLMClient,
	}
}

func (s *FinetuneCompletionService) ValidateOwnership(ctx context.Context, projectID, ownerID uuid.UUID) error {
	project, err := s.projectRepository.GetByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	if project.OwnerID != ownerID {
		return errors.New("unauthorized: project does not belong to user")
	}

	return nil
}

func (s *FinetuneCompletionService) GetFinetuneModelName(ctx context.Context, finetuneID uuid.UUID) (string, error) {
	finetune, err := s.finetuneRepository.GetByID(ctx, finetuneID)
	if err != nil {
		return "", errors.New("finetune not found")
	}

	if finetune.Status != entities.FinetuneStatusDone {
		return "", errors.New("finetune is not ready for inference")
	}

	return finetune.ModelName, nil
}

func (s *FinetuneCompletionService) GenerateCompletion(ctx context.Context, finetuneID uuid.UUID, modelName, prompt string, maxTokens int, temperature, topP float64) (string, error) {
	// Set defaults if not provided
	if maxTokens == 0 {
		maxTokens = 512
	}
	if topP == 0 {
		topP = 0.9
	}

	finetuneIDStr := finetuneID.String()
	result, err := s.ollamaLLMClient.GenerateCompletion(ctx, &finetuneIDStr, prompt, modelName, maxTokens, temperature, topP)
	if err != nil {
		return "", err
	}

	return result.Response, nil
}
