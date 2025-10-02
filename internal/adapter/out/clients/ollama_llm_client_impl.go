package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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

func (c *OllamaLLMClientImpl) GenerateCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64, topP float64) (string, error) {
	// Build the request payload
	requestPayload := map[string]interface{}{
		"input": map[string]interface{}{
			"openai_route": "/v1/completions",
			"openai_input": map[string]interface{}{
				"model":       "qwen3:30b-a3b-instruct-2507-q4_K_M",
				"prompt":      prompt,
				"max_tokens":  maxTokens,
				"temperature": temperature,
				"top_p":       topP,
			},
		},
	}

	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request to JSON: %w", err)
	}

	// Create HTTP request to Runpod API
	url := fmt.Sprintf("https://api.runpod.ai/v2/%s/run", c.podID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(requestJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Runpod API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Runpod API returned status code %d", resp.StatusCode)
	}

	// Parse response
	var responseData OllamaLLMResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract the response text from the completion
	if len(responseData.Output.Choices) == 0 {
		return "", fmt.Errorf("no completion choices in response")
	}

	return responseData.Output.Choices[0].Text, nil
}
