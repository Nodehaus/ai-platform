package use_cases

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/clients"
)

func TestPublicChatCompletionUseCaseImpl_Success(t *testing.T) {
	finetuneID := uuid.New()
	deploymentID := uuid.New()

	mockClient := &mockOllamaLLMClient{
		result: &clients.OllamaLLMClientResult{
			Response:      "This is a chat response",
			TokensIn:      15,
			TokensOut:     25,
			DelayTime:     150,
			ExecutionTime: 600,
		},
		err: nil,
	}

	mockLogsRepo := &mockDeploymentLogsRepository{
		logs: []*entities.DeploymentLogs{},
	}

	useCase := &PublicChatCompletionUseCaseImpl{
		OllamaLLMClient:          mockClient,
		DeploymentLogsRepository: mockLogsRepo,
	}

	messages := []in.ChatMessage{
		{Role: "user", Content: "Hello"},
		{Role: "user", Content: "How are you?"},
	}
	command := in.PublicChatCompletionCommand{
		DeploymentID: deploymentID,
		FinetuneID:   &finetuneID,
		ModelName:    "test-model",
		Messages:     messages,
		MaxTokens:    100,
		Temperature:  0.7,
		TopP:         0.9,
	}

	result, err := useCase.GenerateChatCompletion(context.Background(), command)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Response != "This is a chat response" {
		t.Errorf("Expected response 'This is a chat response', got %s", result.Response)
	}

	if len(mockLogsRepo.logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(mockLogsRepo.logs))
	}

	log := mockLogsRepo.logs[0]
	if log.DeploymentID != deploymentID {
		t.Errorf("Expected deployment ID %s, got %s", deploymentID, log.DeploymentID)
	}

	if log.TokensIn != 15 {
		t.Errorf("Expected tokens in 15, got %d", log.TokensIn)
	}

	if log.TokensOut != 25 {
		t.Errorf("Expected tokens out 25, got %d", log.TokensOut)
	}

	expectedInput := `[{"role":"user","content":"Hello"},{"role":"user","content":"How are you?"}]`
	if log.Input != expectedInput {
		t.Errorf("Expected input '%s', got %s", expectedInput, log.Input)
	}

	if log.Output != "This is a chat response" {
		t.Errorf("Expected output 'This is a chat response', got %s", log.Output)
	}

	if log.DelayTime != 150 {
		t.Errorf("Expected delay time 150, got %d", log.DelayTime)
	}

	if log.ExecutionTime != 600 {
		t.Errorf("Expected execution time 600, got %d", log.ExecutionTime)
	}
}
