package services

import (
	"context"
	"testing"

	"ai-platform/internal/application/port/out/clients"
)

// MockOllamaLLMClient is a mock implementation for testing
type MockOllamaLLMClient struct {
	GenerateCompletionFunc     func(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error)
	GenerateChatCompletionFunc func(ctx context.Context, finetuneID *string, messages []string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error)
}

func (m *MockOllamaLLMClient) GenerateCompletion(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error) {
	if m.GenerateCompletionFunc != nil {
		return m.GenerateCompletionFunc(ctx, finetuneID, prompt, model, maxTokens, temperature, topP)
	}
	return &clients.OllamaLLMClientResult{Response: ""}, nil
}

func (m *MockOllamaLLMClient) GenerateChatCompletion(ctx context.Context, finetuneID *string, messages []string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error) {
	if m.GenerateChatCompletionFunc != nil {
		return m.GenerateChatCompletionFunc(ctx, finetuneID, messages, model, maxTokens, temperature, topP)
	}
	return &clients.OllamaLLMClientResult{Response: ""}, nil
}

func TestPromptAnalysisService_ExtractJSON(t *testing.T) {
	service := &PromptAnalysisService{}

	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name: "Extract JSON with markdown code block",
			response: "```json\n{\"test\": \"value\"}\n```",
			wantErr:  false,
		},
		{
			name:     "Extract JSON with plain code block",
			response: "```\n{\"test\": \"value\"}\n```",
			wantErr:  false,
		},
		{
			name:     "Extract JSON directly",
			response: "Some text {\"test\": \"value\"} more text",
			wantErr:  false,
		},
		{
			name:     "No JSON found",
			response: "No JSON here",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.extractJSON(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == "" {
				t.Errorf("extractJSON() returned empty string")
			}
		})
	}
}

func TestPromptAnalysisService_GetJSONStructure(t *testing.T) {
	mockClient := &MockOllamaLLMClient{
		GenerateCompletionFunc: func(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error) {
			return &clients.OllamaLLMClientResult{
				Response: `{
					"json_object_fields": {
						"question": "The question to ask",
						"answer": "The answer to the question"
					},
					"input_field": "question",
					"output_field": "answer",
					"expected_output_size_chars": 150
				}`,
			}, nil
		},
	}

	service := &PromptAnalysisService{
		OllamaLLMClient: mockClient,
	}

	result, err := service.GetJSONStructure(context.Background(), "Create a Q&A dataset")
	if err != nil {
		t.Fatalf("GetJSONStructure() error = %v", err)
	}

	if result.InputField != "question" {
		t.Errorf("Expected input_field 'question', got '%s'", result.InputField)
	}
	if result.OutputField != "answer" {
		t.Errorf("Expected output_field 'answer', got '%s'", result.OutputField)
	}
	if result.ExpectedOutputSizeChars != 150 {
		t.Errorf("Expected expected_output_size_chars 150, got %d", result.ExpectedOutputSizeChars)
	}
	if len(result.JSONObjectFields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(result.JSONObjectFields))
	}
}

func TestPromptAnalysisService_GetPromptAnalysis(t *testing.T) {
	mockClient := &MockOllamaLLMClient{
		GenerateCompletionFunc: func(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens int, temperature float64, topP float64) (*clients.OllamaLLMClientResult, error) {
			return &clients.OllamaLLMClientResult{
				Response: "This is a good prompt. Consider adding more specific examples.",
			}, nil
		},
	}

	service := &PromptAnalysisService{
		OllamaLLMClient: mockClient,
	}

	result, err := service.GetPromptAnalysis(context.Background(), "Test prompt")
	if err != nil {
		t.Fatalf("GetPromptAnalysis() error = %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty analysis result")
	}
}

func TestJSONStructureResponseType(t *testing.T) {
	// Test that JSONStructureResponse is properly defined
	response := JSONStructureResponse{
		JSONObjectFields: map[string]string{
			"field1": "description1",
		},
		InputField:              "field1",
		OutputField:             "field2",
		ExpectedOutputSizeChars: 100,
	}

	if response.InputField != "field1" {
		t.Error("JSONStructureResponse field assignment failed")
	}
}
