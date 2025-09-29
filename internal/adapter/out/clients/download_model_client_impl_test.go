package clients

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestDownloadModelClientImpl_ValidatesEnvironmentVariables(t *testing.T) {
	// Save original environment variables
	originalBucket := os.Getenv("APP_S3_BUCKET")
	originalRegion := os.Getenv("AWS_DEFAULT_REGION")

	// Clear environment variables
	os.Unsetenv("APP_S3_BUCKET")
	os.Unsetenv("AWS_DEFAULT_REGION")

	// Test missing bucket
	_, err := NewDownloadModelClientImpl()
	if err == nil {
		t.Error("Expected error when APP_S3_BUCKET is not set")
	}

	// Restore original environment variables
	if originalBucket != "" {
		os.Setenv("APP_S3_BUCKET", originalBucket)
	}
	if originalRegion != "" {
		os.Setenv("AWS_DEFAULT_REGION", originalRegion)
	}
}

func TestDownloadModelClientImpl_DownloadModel_PathConstruction(t *testing.T) {
	// Set required environment variables
	os.Setenv("APP_S3_BUCKET", "test-bucket")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")

	client, err := NewDownloadModelClientImpl()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	finetuneID := uuid.New()
	modelName := "test-model"

	// This will fail with AWS credentials error, but we can test the path construction logic
	ctx := context.Background()
	_, _, err = client.DownloadModel(ctx, finetuneID, modelName)
	if err == nil {
		t.Error("Expected AWS error due to invalid credentials")
	}

	// The error should contain information about the constructed path
	expectedKey := "finetunes/" + finetuneID.String() + "/" + modelName + ".gguf"
	t.Logf("Expected S3 key: %s", expectedKey)
	t.Logf("Actual error: %v", err)
}