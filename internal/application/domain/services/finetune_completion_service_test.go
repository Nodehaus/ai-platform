package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/out/clients"
)

// Mock repositories and clients
type MockFinetuneRepository struct {
	mock.Mock
}

func (m *MockFinetuneRepository) Create(ctx context.Context, finetune *entities.Finetune) error {
	args := m.Called(ctx, finetune)
	return args.Error(0)
}

func (m *MockFinetuneRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Finetune, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Finetune), args.Error(1)
}

func (m *MockFinetuneRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Finetune, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]*entities.Finetune), args.Error(1)
}

func (m *MockFinetuneRepository) GetLatestByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.Finetune, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Finetune), args.Error(1)
}

func (m *MockFinetuneRepository) Update(ctx context.Context, finetune *entities.Finetune) error {
	args := m.Called(ctx, finetune)
	return args.Error(0)
}

func (m *MockFinetuneRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.FinetuneStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockFinetuneRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFinetuneRepository) GetNextVersion(ctx context.Context, projectID uuid.UUID) (int, error) {
	args := m.Called(ctx, projectID)
	return args.Int(0), args.Error(1)
}

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) GetByID(id uuid.UUID) (*entities.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByOwnerID(ownerID uuid.UUID) ([]entities.Project, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]entities.Project), args.Error(1)
}

func (m *MockProjectRepository) GetActiveByOwnerID(ownerID uuid.UUID) ([]entities.Project, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]entities.Project), args.Error(1)
}

func (m *MockProjectRepository) ExistsByNameAndOwnerID(name string, ownerID uuid.UUID) (bool, error) {
	args := m.Called(name, ownerID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) Create(project *entities.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Update(project *entities.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockOllamaLLMClient struct {
	mock.Mock
}

func (m *MockOllamaLLMClient) GenerateCompletion(ctx context.Context, finetuneID string, prompt string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error) {
	args := m.Called(ctx, finetuneID, prompt, model, maxTokens, temperature, topP)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clients.OllamaLLMClientResult), args.Error(1)
}

func (m *MockOllamaLLMClient) GenerateChatCompletion(ctx context.Context, finetuneID string, messages []string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error) {
	args := m.Called(ctx, finetuneID, messages, model, maxTokens, temperature, topP)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clients.OllamaLLMClientResult), args.Error(1)
}

func TestValidateOwnership_Success(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	ownerID := uuid.New()

	mockProjectRepo := new(MockProjectRepository)
	mockFinetuneRepo := new(MockFinetuneRepository)
	mockOllamaClient := new(MockOllamaLLMClient)

	project := &entities.Project{
		ID:      projectID,
		OwnerID: ownerID,
		Name:    "Test Project",
	}

	mockProjectRepo.On("GetByID", projectID).Return(project, nil)

	service := services.NewFinetuneCompletionService(mockFinetuneRepo, mockProjectRepo, mockOllamaClient)

	err := service.ValidateOwnership(ctx, projectID, ownerID)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
}

func TestValidateOwnership_ProjectNotFound(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	ownerID := uuid.New()

	mockProjectRepo := new(MockProjectRepository)
	mockFinetuneRepo := new(MockFinetuneRepository)
	mockOllamaClient := new(MockOllamaLLMClient)

	mockProjectRepo.On("GetByID", projectID).Return(nil, errors.New("not found"))

	service := services.NewFinetuneCompletionService(mockFinetuneRepo, mockProjectRepo, mockOllamaClient)

	err := service.ValidateOwnership(ctx, projectID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, "project not found", err.Error())
	mockProjectRepo.AssertExpectations(t)
}

func TestValidateOwnership_Unauthorized(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	ownerID := uuid.New()
	differentOwnerID := uuid.New()

	mockProjectRepo := new(MockProjectRepository)
	mockFinetuneRepo := new(MockFinetuneRepository)
	mockOllamaClient := new(MockOllamaLLMClient)

	project := &entities.Project{
		ID:      projectID,
		OwnerID: differentOwnerID,
		Name:    "Test Project",
	}

	mockProjectRepo.On("GetByID", projectID).Return(project, nil)

	service := services.NewFinetuneCompletionService(mockFinetuneRepo, mockProjectRepo, mockOllamaClient)

	err := service.ValidateOwnership(ctx, projectID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized: project does not belong to user", err.Error())
	mockProjectRepo.AssertExpectations(t)
}

func TestGetFinetuneModelName_Success(t *testing.T) {
	ctx := context.Background()
	finetuneID := uuid.New()

	mockProjectRepo := new(MockProjectRepository)
	mockFinetuneRepo := new(MockFinetuneRepository)
	mockOllamaClient := new(MockOllamaLLMClient)

	finetune := &entities.Finetune{
		ID:        finetuneID,
		ModelName: "test_model_v1",
		Status:    entities.FinetuneStatusDone,
	}

	mockFinetuneRepo.On("GetByID", ctx, finetuneID).Return(finetune, nil)

	service := services.NewFinetuneCompletionService(mockFinetuneRepo, mockProjectRepo, mockOllamaClient)

	modelName, err := service.GetFinetuneModelName(ctx, finetuneID)

	assert.NoError(t, err)
	assert.Equal(t, "test_model_v1", modelName)
	mockFinetuneRepo.AssertExpectations(t)
}

func TestGetFinetuneModelName_FinetuneNotReady(t *testing.T) {
	ctx := context.Background()
	finetuneID := uuid.New()

	mockProjectRepo := new(MockProjectRepository)
	mockFinetuneRepo := new(MockFinetuneRepository)
	mockOllamaClient := new(MockOllamaLLMClient)

	finetune := &entities.Finetune{
		ID:        finetuneID,
		ModelName: "test_model_v1",
		Status:    entities.FinetuneStatusRunning,
	}

	mockFinetuneRepo.On("GetByID", ctx, finetuneID).Return(finetune, nil)

	service := services.NewFinetuneCompletionService(mockFinetuneRepo, mockProjectRepo, mockOllamaClient)

	modelName, err := service.GetFinetuneModelName(ctx, finetuneID)

	assert.Error(t, err)
	assert.Equal(t, "finetune is not ready for inference", err.Error())
	assert.Equal(t, "", modelName)
	mockFinetuneRepo.AssertExpectations(t)
}

func TestGenerateCompletion_Success(t *testing.T) {
	ctx := context.Background()
	finetuneID := uuid.New()

	mockProjectRepo := new(MockProjectRepository)
	mockFinetuneRepo := new(MockFinetuneRepository)
	mockOllamaClient := new(MockOllamaLLMClient)

	mockOllamaClient.On("GenerateCompletion", ctx, finetuneID.String(), "test prompt", "test_model", 512, 0.7, 0.9).Return(&clients.OllamaLLMClientResult{Response: "completion result"}, nil)

	service := services.NewFinetuneCompletionService(mockFinetuneRepo, mockProjectRepo, mockOllamaClient)

	completion, err := service.GenerateCompletion(ctx, finetuneID, "test_model", "test prompt", 0, 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, "completion result", completion)
	mockOllamaClient.AssertExpectations(t)
}

func TestGenerateCompletion_WithCustomParameters(t *testing.T) {
	ctx := context.Background()
	finetuneID := uuid.New()

	mockProjectRepo := new(MockProjectRepository)
	mockFinetuneRepo := new(MockFinetuneRepository)
	mockOllamaClient := new(MockOllamaLLMClient)

	mockOllamaClient.On("GenerateCompletion", ctx, finetuneID.String(), "test prompt", "test_model", 1024, 0.8, 0.95).Return(&clients.OllamaLLMClientResult{Response: "completion result"}, nil)

	service := services.NewFinetuneCompletionService(mockFinetuneRepo, mockProjectRepo, mockOllamaClient)

	completion, err := service.GenerateCompletion(ctx, finetuneID, "test_model", "test prompt", 1024, 0.8, 0.95)

	assert.NoError(t, err)
	assert.Equal(t, "completion result", completion)
	mockOllamaClient.AssertExpectations(t)
}
