package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	portClients "ai-platform/internal/application/port/out/clients"
)

type OllamaLLMClientImpl struct {
	apiKey string
	podID  string
	client *http.Client
}

func NewOllamaLLMClientImpl() (*OllamaLLMClientImpl, error) {
	apiKey := os.Getenv("RUNPOD_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RUNPOD_API_KEY environment variable is required")
	}

	podID := os.Getenv("RUNPOD_POD_ID_OLLAMA")
	if podID == "" {
		return nil, fmt.Errorf("RUNPOD_POD_ID_OLLAMA environment variable is required")
	}

	return &OllamaLLMClientImpl{
		apiKey: apiKey,
		podID:  podID,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

func (c *OllamaLLMClientImpl) GenerateCompletion(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens int, temperature float64, topP float64) (*portClients.OllamaLLMClientResult, error) {
	openaiInput := map[string]interface{}{
		"model":       model,
		"prompt":      prompt,
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"top_p":       topP,
	}

	return c.callRunpodAPI(ctx, finetuneID, "/v1/completions", openaiInput)
}

func (c *OllamaLLMClientImpl) GenerateChatCompletion(ctx context.Context, finetuneID *string, messages []portClients.ChatMessage, model string, maxTokens int, temperature float64, topP float64) (*portClients.OllamaLLMClientResult, error) {
	openaiInput := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"top_p":       topP,
	}

	return c.callRunpodAPI(ctx, finetuneID, "/v1/chat/completions", openaiInput)
}

// callRunpodAPI is a common method to make API calls to Runpod
func (c *OllamaLLMClientImpl) callRunpodAPI(ctx context.Context, finetuneID *string, openaiRoute string, openaiInput map[string]interface{}) (*portClients.OllamaLLMClientResult, error) {
	// Build the request payload
	bucket := os.Getenv("APP_S3_BUCKET")
	appEnv := os.Getenv("APP_ENV")

	inputPayload := map[string]interface{}{
		"s3_bucket":    bucket,
		"app_env":      appEnv,
		"openai_route": openaiRoute,
		"openai_input": openaiInput,
	}

	// Only include finetune_id if it's not nil
	if finetuneID != nil {
		inputPayload["finetune_id"] = *finetuneID
	}

	requestPayload := map[string]interface{}{
		"input": inputPayload,
	}

	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request to JSON: %w", err)
	}

	// Create HTTP request to Runpod API
	url := fmt.Sprintf("https://api.runpod.ai/v2/%s/runsync", c.podID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(requestJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Runpod API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Runpod API returned status code %d", resp.StatusCode)
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var responseData OllamaLLMResponseModel
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if status is COMPLETED
	if responseData.Status != "COMPLETED" {
		return nil, fmt.Errorf("runpod job status is %s, not COMPLETED", responseData.Status)
	}

	// Extract the response text from the completion
	if len(responseData.Output) == 0 {
		return nil, fmt.Errorf("no output in response")
	}

	if len(responseData.Output[0].Choices) == 0 {
		return nil, fmt.Errorf("no completion choices in response")
	}

	// Extract response content - handle both completion (text) and chat completion (message)
	var responseText string
	choice := responseData.Output[0].Choices[0]
	if choice.Message != nil {
		// Chat completion response
		responseText = choice.Message.Content
	} else {
		// Regular completion response
		responseText = choice.Text
	}

	return &portClients.OllamaLLMClientResult{
		Response:      responseText,
		TokensIn:      responseData.Output[0].Usage.PromptTokens,
		TokensOut:     responseData.Output[0].Usage.CompletionTokens,
		DelayTime:     responseData.DelayTime,
		ExecutionTime: responseData.ExecutionTime,
	}, nil
}
