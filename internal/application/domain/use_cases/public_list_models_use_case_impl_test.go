package use_cases

import (
	"context"
	"testing"
	"time"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"github.com/google/uuid"
)

type mockDeploymentRepository struct {
	deployments []entities.Deployment
	err         error
}

func (m *mockDeploymentRepository) Create(deployment *entities.Deployment) error {
	return m.err
}

func (m *mockDeploymentRepository) GetByID(id uuid.UUID) (*entities.Deployment, error) {
	return nil, m.err
}

func (m *mockDeploymentRepository) GetByProjectID(projectID uuid.UUID) ([]entities.Deployment, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.deployments, nil
}

func (m *mockDeploymentRepository) GetByFinetuneID(finetuneID uuid.UUID) (*entities.Deployment, error) {
	return nil, m.err
}

func (m *mockDeploymentRepository) GetByProjectIDAndModelName(projectID uuid.UUID, modelName string) (*entities.Deployment, error) {
	return nil, m.err
}

func (m *mockDeploymentRepository) GetByAPIKey(apiKey string) (*entities.Deployment, error) {
	return nil, m.err
}

func (m *mockDeploymentRepository) Delete(id uuid.UUID) error {
	return m.err
}

func TestPublicListModelsUseCaseImpl_Success(t *testing.T) {
	projectID := uuid.New()
	deployment1ID := uuid.New()
	deployment2ID := uuid.New()
	finetuneID := uuid.New()

	now := time.Now()

	mockRepo := &mockDeploymentRepository{
		deployments: []entities.Deployment{
			{
				ID:         deployment1ID,
				ModelName:  "model-1",
				APIKey:     "api-key-1",
				ProjectID:  projectID,
				FinetuneID: &finetuneID,
				CreatedAt:  now.Add(-24 * time.Hour),
				UpdatedAt:  now.Add(-24 * time.Hour),
			},
			{
				ID:         deployment2ID,
				ModelName:  "model-2",
				APIKey:     "api-key-2",
				ProjectID:  projectID,
				FinetuneID: nil,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
		err: nil,
	}

	useCase := NewPublicListModelsUseCaseImpl(mockRepo)

	command := in.PublicListModelsCommand{
		ProjectID: projectID,
	}

	result, err := useCase.ListModels(context.Background(), command)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Object != "list" {
		t.Errorf("Expected object 'list', got %s", result.Object)
	}

	if len(result.Data) != 2 {
		t.Fatalf("Expected 2 models, got %d", len(result.Data))
	}

	// Check first model
	model1 := result.Data[0]
	if model1.ID != "model-1" {
		t.Errorf("Expected model ID 'model-1', got %s", model1.ID)
	}
	if model1.Object != "model" {
		t.Errorf("Expected object 'model', got %s", model1.Object)
	}
	if model1.Created != now.Add(-24*time.Hour).Unix() {
		t.Errorf("Expected created timestamp %d, got %d", now.Add(-24*time.Hour).Unix(), model1.Created)
	}
	if model1.OwnedBy != "organization" {
		t.Errorf("Expected owned_by 'organization', got %s", model1.OwnedBy)
	}

	// Check second model
	model2 := result.Data[1]
	if model2.ID != "model-2" {
		t.Errorf("Expected model ID 'model-2', got %s", model2.ID)
	}
	if model2.Object != "model" {
		t.Errorf("Expected object 'model', got %s", model2.Object)
	}
	if model2.Created != now.Unix() {
		t.Errorf("Expected created timestamp %d, got %d", now.Unix(), model2.Created)
	}
	if model2.OwnedBy != "organization" {
		t.Errorf("Expected owned_by 'organization', got %s", model2.OwnedBy)
	}
}

func TestPublicListModelsUseCaseImpl_EmptyList(t *testing.T) {
	projectID := uuid.New()

	mockRepo := &mockDeploymentRepository{
		deployments: []entities.Deployment{},
		err:         nil,
	}

	useCase := NewPublicListModelsUseCaseImpl(mockRepo)

	command := in.PublicListModelsCommand{
		ProjectID: projectID,
	}

	result, err := useCase.ListModels(context.Background(), command)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Object != "list" {
		t.Errorf("Expected object 'list', got %s", result.Object)
	}

	if len(result.Data) != 0 {
		t.Errorf("Expected 0 models, got %d", len(result.Data))
	}
}
