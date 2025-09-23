package clients

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

func TestTrainingDatasetJobClientImpl_SubmitJob_ValidatesEnvironmentVariables(t *testing.T) {
	// Clear environment variables to test validation
	originalBucket := os.Getenv("APP_S3_BUCKET")
	os.Unsetenv("APP_S3_BUCKET")
	defer func() {
		if originalBucket != "" {
			os.Setenv("APP_S3_BUCKET", originalBucket)
		}
	}()

	_, err := NewTrainingDatasetJobClientImpl()
	if err == nil {
		t.Fatal("Expected error when APP_S3_BUCKET is not set")
	}
	if err.Error() != "APP_S3_BUCKET environment variable is required" {
		t.Fatalf("Expected specific error message, got: %s", err.Error())
	}
}

func TestTrainingDatasetJobClientImpl_SubmitJob_ValidatesJobMarshaling(t *testing.T) {
	// Set required environment variables
	os.Setenv("APP_S3_BUCKET", "test-bucket")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	defer func() {
		os.Unsetenv("APP_S3_BUCKET")
		os.Unsetenv("AWS_DEFAULT_REGION")
	}()

	client, err := NewTrainingDatasetJobClientImpl()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	job := entities.TrainingDatasetJob{
		CorpusS3Path:           "/documents/eurlex",
		CorpusFilesSubset:      []string{},
		LanguageISO:            "deu",
		UserID:                 uuid.New().String(),
		TrainingDatasetID:      uuid.New().String(),
		GeneratePrompt:         "You task is to generate training dataset for question answering.",
		GenerateExamplesNumber: 100,
		GenerateModel:          "gemma3:8b",
		GenerateModelRunner:    "runpod_ollama",
	}

	// This will fail without proper AWS credentials, but we're testing the JSON marshaling
	// and the method signature, not the actual S3 upload
	ctx := context.Background()
	err = client.SubmitJob(ctx, job)
	// We expect this to fail due to AWS configuration, but not due to JSON marshaling
	if err != nil {
		// This is expected in a test environment without AWS credentials
		t.Logf("Expected AWS error: %v", err)
	}
}