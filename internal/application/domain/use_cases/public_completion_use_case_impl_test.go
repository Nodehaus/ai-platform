package use_cases

import (
	"context"
	"testing"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
	"github.com/google/uuid"
)

type mockOllamaLLMClient struct {
	result *clients.OllamaLLMClientResult
	err    error
}

func (m *mockOllamaLLMClient) GenerateCompletion(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error) {
	return m.result, m.err
}

func (m *mockOllamaLLMClient) GenerateChatCompletion(ctx context.Context, finetuneID *string, messages []clients.ChatMessage, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error) {
	return m.result, m.err
}

func (m *mockOllamaLLMClient) GenerateChatCompletionStream(ctx context.Context, finetuneID *string, messages []clients.ChatMessage, model string, maxTokens int, temperature float64, topP float64) (<-chan clients.StreamChunk, error) {
	ch := make(chan clients.StreamChunk)
	close(ch)
	return ch, m.err
}

type mockDeploymentLogsRepository struct {
	logs []*entities.DeploymentLogs
	err  error
}

func (m *mockDeploymentLogsRepository) Create(log *entities.DeploymentLogs) error {
	if m.err != nil {
		return m.err
	}
	m.logs = append(m.logs, log)
	return nil
}

func (m *mockDeploymentLogsRepository) GetLatest(deploymentID uuid.UUID, limit int) ([]*entities.DeploymentLogs, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.logs, nil
}

func (m *mockDeploymentLogsRepository) GetAll(deploymentID uuid.UUID) ([]*entities.DeploymentLogs, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.logs, nil
}

func TestPublicCompletionUseCaseImpl_Success(t *testing.T) {
	finetuneID := uuid.New()
	deploymentID := uuid.New()

	mockClient := &mockOllamaLLMClient{
		result: &clients.OllamaLLMClientResult{
			Response:      "This is a test response",
			TokensIn:      10,
			TokensOut:     20,
			DelayTime:     100,
			ExecutionTime: 500,
		},
		err: nil,
	}

	mockLogsRepo := &mockDeploymentLogsRepository{
		logs: []*entities.DeploymentLogs{},
	}

	useCase := &PublicCompletionUseCaseImpl{
		OllamaLLMClient:          mockClient,
		DeploymentLogsRepository: mockLogsRepo,
	}

	command := in.PublicCompletionCommand{
		DeploymentID: deploymentID,
		FinetuneID:   &finetuneID,
		ModelName:    "test-model",
		Prompt:       "Test prompt",
		MaxTokens:    100,
		Temperature:  0.7,
		TopP:         0.9,
	}

	result, err := useCase.GenerateCompletion(context.Background(), command)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Response != "This is a test response" {
		t.Errorf("Expected response 'This is a test response', got %s", result.Response)
	}

	if len(mockLogsRepo.logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(mockLogsRepo.logs))
	}

	log := mockLogsRepo.logs[0]
	if log.DeploymentID != deploymentID {
		t.Errorf("Expected deployment ID %s, got %s", deploymentID, log.DeploymentID)
	}

	if log.TokensIn != 10 {
		t.Errorf("Expected tokens in 10, got %d", log.TokensIn)
	}

	if log.TokensOut != 20 {
		t.Errorf("Expected tokens out 20, got %d", log.TokensOut)
	}

	if log.Input != "Test prompt" {
		t.Errorf("Expected input 'Test prompt', got %s", log.Input)
	}

	if log.Output != "This is a test response" {
		t.Errorf("Expected output 'This is a test response', got %s", log.Output)
	}

	if log.DelayTime != 100 {
		t.Errorf("Expected delay time 100, got %d", log.DelayTime)
	}

	if log.ExecutionTime != 500 {
		t.Errorf("Expected execution time 500, got %d", log.ExecutionTime)
	}
}
