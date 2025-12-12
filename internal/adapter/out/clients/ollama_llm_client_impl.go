package clients

import (
	"bufio"
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

func (c *OllamaLLMClientImpl) GenerateCompletionStream(ctx context.Context, finetuneID *string, prompt string, model string, maxTokens int, temperature float64, topP float64) (<-chan portClients.StreamChunk, error) {
	openaiInput := map[string]interface{}{
		"model":       model,
		"prompt":      prompt,
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"top_p":       topP,
		"stream":      true,
	}

	// Build the request payload
	bucket := os.Getenv("APP_S3_BUCKET")
	appEnv := os.Getenv("APP_ENV")

	inputPayload := map[string]interface{}{
		"s3_bucket":    bucket,
		"app_env":      appEnv,
		"openai_route": "/v1/completions",
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

	// Create HTTP request to Runpod API /run endpoint (for streaming)
	url := fmt.Sprintf("https://api.runpod.ai/v2/%s/run", c.podID)
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

	// Read response body to get run_id
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var runResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(bodyBytes, &runResponse); err != nil {
		return nil, fmt.Errorf("failed to decode run response: %w", err)
	}

	runID := runResponse.ID

	// Create a channel for streaming chunks
	chunkChan := make(chan portClients.StreamChunk)

	// Start goroutine to poll the stream endpoint
	go func() {
		defer close(chunkChan)

		streamURL := fmt.Sprintf("https://api.runpod.ai/v2/%s/stream/%s", c.podID, runID)

		for {
			select {
			case <-ctx.Done():
				chunkChan <- portClients.StreamChunk{
					Error: ctx.Err(),
				}
				return
			default:
			}

			// Poll the stream endpoint
			streamReq, err := http.NewRequestWithContext(ctx, "GET", streamURL, nil)
			if err != nil {
				chunkChan <- portClients.StreamChunk{
					Error: fmt.Errorf("failed to create stream request: %w", err),
				}
				return
			}

			streamReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

			streamResp, err := c.client.Do(streamReq)
			if err != nil {
				chunkChan <- portClients.StreamChunk{
					Error: fmt.Errorf("failed to send stream request: %w", err),
				}
				return
			}

			// Read response line by line
			scanner := bufio.NewScanner(streamResp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}

				var data map[string]interface{}
				if err := json.Unmarshal([]byte(line), &data); err != nil {
					continue
				}

				status, _ := data["status"].(string)

				if status == "IN_PROGRESS" {
					// Extract stream chunks
					if stream, ok := data["stream"].([]interface{}); ok {
						for _, item := range stream {
							chunk, ok := item.(map[string]interface{})
							if !ok {
								continue
							}

							output, ok := chunk["output"].(string)
							if !ok {
								continue
							}

							// Parse the SSE format: "data: {...}"
							if !bytes.HasPrefix([]byte(output), []byte("data: ")) {
								continue
							}

							content := output[6:] // Remove "data: " prefix

							if content == "[DONE]" {
								streamResp.Body.Close()
								return
							}

							// Parse the JSON content
							var chunkData map[string]interface{}
							if err := json.Unmarshal([]byte(content), &chunkData); err != nil {
								continue
							}

							// Extract content from choices[0].delta.content or choices[0].text
							if choices, ok := chunkData["choices"].([]interface{}); ok && len(choices) > 0 {
								if choice, ok := choices[0].(map[string]interface{}); ok {
									// For completions endpoint, check for "text" field
									if text, ok := choice["text"].(string); ok {
										chunkChan <- portClients.StreamChunk{
											Content: text,
										}
									}

									// Check for finish_reason
									if finishReason, ok := choice["finish_reason"].(string); ok && finishReason != "" && finishReason != "null" {
										reason := finishReason
										chunkChan <- portClients.StreamChunk{
											FinishReason: &reason,
										}
									}
								}
							}
						}
					}
				} else if status == "COMPLETED" {
					streamResp.Body.Close()
					return
				} else if status == "FAILED" {
					streamResp.Body.Close()
					chunkChan <- portClients.StreamChunk{
						Error: fmt.Errorf("runpod job failed"),
					}
					return
				}
			}

			streamResp.Body.Close()

			if err := scanner.Err(); err != nil {
				chunkChan <- portClients.StreamChunk{
					Error: fmt.Errorf("error reading stream: %w", err),
				}
				return
			}

			// Small delay before polling again
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return chunkChan, nil
}

func (c *OllamaLLMClientImpl) GenerateChatCompletion(ctx context.Context, finetuneID *string, messages []portClients.ChatMessage, model string, maxTokens int, temperature float64, topP float64) (*portClients.OllamaLLMClientResult, error) {
	openaiInput := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"top_p":       topP,
		"stream":      false,
	}

	return c.callRunpodAPI(ctx, finetuneID, "/v1/chat/completions", openaiInput)
}

func (c *OllamaLLMClientImpl) GenerateChatCompletionStream(ctx context.Context, finetuneID *string, messages []portClients.ChatMessage, model string, maxTokens int, temperature float64, topP float64) (<-chan portClients.StreamChunk, error) {
	openaiInput := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"top_p":       topP,
		"stream":      true,
	}

	// Build the request payload
	bucket := os.Getenv("APP_S3_BUCKET")
	appEnv := os.Getenv("APP_ENV")

	inputPayload := map[string]interface{}{
		"s3_bucket":    bucket,
		"app_env":      appEnv,
		"openai_route": "/v1/chat/completions",
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

	// Create HTTP request to Runpod API /run endpoint (for streaming)
	url := fmt.Sprintf("https://api.runpod.ai/v2/%s/run", c.podID)
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

	// Read response body to get run_id
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var runResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(bodyBytes, &runResponse); err != nil {
		return nil, fmt.Errorf("failed to decode run response: %w", err)
	}

	runID := runResponse.ID

	// Create a channel for streaming chunks
	chunkChan := make(chan portClients.StreamChunk)

	// Start goroutine to poll the stream endpoint
	go func() {
		defer close(chunkChan)

		streamURL := fmt.Sprintf("https://api.runpod.ai/v2/%s/stream/%s", c.podID, runID)

		for {
			select {
			case <-ctx.Done():
				chunkChan <- portClients.StreamChunk{
					Error: ctx.Err(),
				}
				return
			default:
			}

			// Poll the stream endpoint
			streamReq, err := http.NewRequestWithContext(ctx, "GET", streamURL, nil)
			if err != nil {
				chunkChan <- portClients.StreamChunk{
					Error: fmt.Errorf("failed to create stream request: %w", err),
				}
				return
			}

			streamReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

			streamResp, err := c.client.Do(streamReq)
			if err != nil {
				chunkChan <- portClients.StreamChunk{
					Error: fmt.Errorf("failed to send stream request: %w", err),
				}
				return
			}

			// Read response line by line
			scanner := bufio.NewScanner(streamResp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}

				var data map[string]interface{}
				if err := json.Unmarshal([]byte(line), &data); err != nil {
					continue
				}

				status, _ := data["status"].(string)

				if status == "IN_PROGRESS" {
					// Extract stream chunks
					if stream, ok := data["stream"].([]interface{}); ok {
						for _, item := range stream {
							chunk, ok := item.(map[string]interface{})
							if !ok {
								continue
							}

							output, ok := chunk["output"].(string)
							if !ok {
								continue
							}

							// Parse the SSE format: "data: {...}"
							if !bytes.HasPrefix([]byte(output), []byte("data: ")) {
								continue
							}

							content := output[6:] // Remove "data: " prefix

							if content == "[DONE]" {
								streamResp.Body.Close()
								return
							}

							// Parse the JSON content
							var chunkData map[string]interface{}
							if err := json.Unmarshal([]byte(content), &chunkData); err != nil {
								continue
							}

							// Extract content from choices[0].delta.content
							if choices, ok := chunkData["choices"].([]interface{}); ok && len(choices) > 0 {
								if choice, ok := choices[0].(map[string]interface{}); ok {
									if delta, ok := choice["delta"].(map[string]interface{}); ok {
										if deltaContent, ok := delta["content"].(string); ok {
											chunkChan <- portClients.StreamChunk{
												Content: deltaContent,
											}
										}
									}

									// Check for finish_reason
									if finishReason, ok := choice["finish_reason"].(string); ok && finishReason != "" && finishReason != "null" {
										reason := finishReason
										chunkChan <- portClients.StreamChunk{
											FinishReason: &reason,
										}
									}
								}
							}
						}
					}
				} else if status == "COMPLETED" {
					streamResp.Body.Close()
					return
				} else if status == "FAILED" {
					streamResp.Body.Close()
					chunkChan <- portClients.StreamChunk{
						Error: fmt.Errorf("runpod job failed"),
					}
					return
				}
			}

			streamResp.Body.Close()

			if err := scanner.Err(); err != nil {
				chunkChan <- portClients.StreamChunk{
					Error: fmt.Errorf("error reading stream: %w", err),
				}
				return
			}

			// Small delay before polling again
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return chunkChan, nil
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
