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

type RunpodClientImpl struct {
	apiKey string
	podID  string
	client *http.Client
}

func NewRunpodClientImpl() (*RunpodClientImpl, error) {
	apiKey := os.Getenv("RUNPOD_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RUNPOD_API_KEY environment variable is required")
	}

	podID := os.Getenv("RUNPOD_POD_ID_FINETUNE")
	if podID == "" {
		return nil, fmt.Errorf("RUNPOD_POD_ID_FINETUNE environment variable is required")
	}

	return &RunpodClientImpl{
		apiKey: apiKey,
		podID:  podID,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *RunpodClientImpl) StartFinetuneJob(ctx context.Context, s3Key string, documentsS3Path string, baseModelName string, modelName string, finetuneID string) error {
	// Create client model with environment configuration
	clientModel := RunpodClientModel{
		S3Bucket:               os.Getenv("APP_S3_BUCKET"),
		TrainingDatasetS3Path:  s3Key,
		DocumentsS3Path:        documentsS3Path,
		BaseModelName:          baseModelName,
		ModelName:              modelName,
		FinetuneID:				finetuneID,
	}

	// Wrap the data in the required "input" field for Runpod API
	requestPayload := map[string]interface{}{
		"input": clientModel,
	}

	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal request to JSON: %w", err)
	}

	// Create HTTP request to Runpod API
	url := fmt.Sprintf("https://api.runpod.ai/v2/%s/run", c.podID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(requestJSON))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Runpod API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Runpod API returned status code %d", resp.StatusCode)
	}

	return nil
}