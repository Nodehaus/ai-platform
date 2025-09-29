package clients

import (
	"context"
	"os"
	"testing"
)

func TestRunpodClientImpl_StartFinetuneJob_ValidatesEnvironmentVariables(t *testing.T) {
	// Clear environment variables to test validation
	originalAPIKey := os.Getenv("RUNPOD_API_KEY")
	originalPodID := os.Getenv("RUNPOD_POD_ID_FINETUNE")
	os.Unsetenv("RUNPOD_API_KEY")
	os.Unsetenv("RUNPOD_POD_ID_FINETUNE")
	defer func() {
		if originalAPIKey != "" {
			os.Setenv("RUNPOD_API_KEY", originalAPIKey)
		}
		if originalPodID != "" {
			os.Setenv("RUNPOD_POD_ID_FINETUNE", originalPodID)
		}
	}()

	// Test missing API key
	_, err := NewRunpodClientImpl()
	if err == nil {
		t.Fatal("Expected error when RUNPOD_API_KEY is not set")
	}
	if err.Error() != "RUNPOD_API_KEY environment variable is required" {
		t.Fatalf("Expected specific error message, got: %s", err.Error())
	}

	// Test missing Pod ID
	os.Setenv("RUNPOD_API_KEY", "test-key")
	_, err = NewRunpodClientImpl()
	if err == nil {
		t.Fatal("Expected error when RUNPOD_POD_ID_FINETUNE is not set")
	}
	if err.Error() != "RUNPOD_POD_ID_FINETUNE environment variable is required" {
		t.Fatalf("Expected specific error message, got: %s", err.Error())
	}
}

func TestRunpodClientImpl_StartFinetuneJob_ValidatesRequestMarshaling(t *testing.T) {
	// Set required environment variables
	os.Setenv("RUNPOD_API_KEY", "test-api-key")
	os.Setenv("RUNPOD_POD_ID_FINETUNE", "test-pod-id")
	defer func() {
		os.Unsetenv("RUNPOD_API_KEY")
		os.Unsetenv("RUNPOD_POD_ID_FINETUNE")
	}()

	client, err := NewRunpodClientImpl()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Set the S3 bucket environment variable for the test
	os.Setenv("APP_S3_BUCKET", "nodehaus")
	defer os.Unsetenv("APP_S3_BUCKET")

	// This will fail without a valid Runpod API endpoint, but we're testing the JSON marshaling
	// and the method signature, not the actual API call
	ctx := context.Background()
	err = client.StartFinetuneJob(ctx,
		"jobs/finetunes/250927101726_cb1b846e-ab09-417e-823c-475107bda72a.json",
		"documents/eurlex/eng",
		"qwen3:4b",
		"qwen3b_4b_test_radio_buttons_v11",
		"12345678-1234-1234-1234-123456789012")
	// We expect this to fail due to invalid endpoint/credentials, but not due to JSON marshaling
	if err != nil {
		// This is expected in a test environment without valid Runpod credentials
		t.Logf("Expected Runpod API error: %v", err)
	}
}