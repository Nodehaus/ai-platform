package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
)

type mockCreateFinetuneUseCase struct {
	result *entities.Finetune
	err    error
}

func (m *mockCreateFinetuneUseCase) Execute(ctx context.Context, command in.CreateFinetuneCommand) (*entities.Finetune, error) {
	if command.BaseModelName == "invalid" {
		return nil, errors.New("invalid base model name")
	}
	if command.TrainingDatasetID == uuid.Nil {
		return nil, errors.New("training dataset not found")
	}
	return m.result, m.err
}

func TestCreateFinetuneController_CreateFinetune_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()
	trainingDatasetID := uuid.New()

	finetune := &entities.Finetune{
		ID:                               uuid.New(),
		ProjectID:                        projectID,
		Version:                          1,
		ModelName:                        "qwen3b_test_project_v1",
		BaseModelName:                    "qwen3b:4b",
		TrainingDatasetID:                trainingDatasetID,
		TrainingDatasetNumberExamples:    intPtr(1000),
		TrainingDatasetSelectRandom:      false,
		Status:                           entities.FinetuneStatusPlanning,
		InferenceSamples:                 []entities.InferenceSample{},
		CreatedAt:                        time.Now(),
		UpdatedAt:                        time.Now(),
	}

	mockUseCase := &mockCreateFinetuneUseCase{
		result: finetune,
		err:    nil,
	}

	controller := &CreateFinetuneController{
		CreateFinetuneUseCase: mockUseCase,
	}

	request := CreateFinetuneRequest{
		BaseModelName:                    "qwen3b:4b",
		TrainingDatasetID:                trainingDatasetID.String(),
		TrainingDatasetNumberExamples:    intPtr(1000),
		TrainingDatasetSelectRandom:      false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)
	c.Params = []gin.Param{{Key: "project_id", Value: projectID.String()}}
	c.Request = httptest.NewRequest("POST", "/api/projects/"+projectID.String()+"/finetunes", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateFinetune(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var response CreateFinetuneResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != finetune.ID {
		t.Errorf("Expected ID %s, got %s", finetune.ID, response.ID)
	}
	if response.ProjectID != projectID {
		t.Errorf("Expected ProjectID %s, got %s", projectID, response.ProjectID)
	}
}

func TestCreateFinetuneController_CreateFinetune_InvalidProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	trainingDatasetID := uuid.New()

	mockUseCase := &mockCreateFinetuneUseCase{}
	controller := &CreateFinetuneController{
		CreateFinetuneUseCase: mockUseCase,
	}

	request := CreateFinetuneRequest{
		BaseModelName:                    "qwen3b:4b",
		TrainingDatasetID:                trainingDatasetID.String(),
		TrainingDatasetNumberExamples:    intPtr(1000),
		TrainingDatasetSelectRandom:      false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)
	c.Params = []gin.Param{{Key: "project_id", Value: "invalid-uuid"}}
	c.Request = httptest.NewRequest("POST", "/api/projects/invalid-uuid/finetunes", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateFinetune(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateFinetuneController_CreateFinetune_InvalidTrainingDatasetID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()

	mockUseCase := &mockCreateFinetuneUseCase{}
	controller := &CreateFinetuneController{
		CreateFinetuneUseCase: mockUseCase,
	}

	request := CreateFinetuneRequest{
		BaseModelName:                    "qwen3b:4b",
		TrainingDatasetID:                "invalid-uuid",
		TrainingDatasetNumberExamples:    intPtr(1000),
		TrainingDatasetSelectRandom:      false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)
	c.Params = []gin.Param{{Key: "project_id", Value: projectID.String()}}
	c.Request = httptest.NewRequest("POST", "/api/projects/"+projectID.String()+"/finetunes", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateFinetune(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateFinetuneController_CreateFinetune_UseCaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()
	trainingDatasetID := uuid.New()

	mockUseCase := &mockCreateFinetuneUseCase{
		result: nil,
		err:    errors.New("training dataset not found"),
	}

	controller := &CreateFinetuneController{
		CreateFinetuneUseCase: mockUseCase,
	}

	request := CreateFinetuneRequest{
		BaseModelName:                    "qwen3b:4b",
		TrainingDatasetID:                trainingDatasetID.String(),
		TrainingDatasetNumberExamples:    intPtr(1000),
		TrainingDatasetSelectRandom:      false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)
	c.Params = []gin.Param{{Key: "project_id", Value: projectID.String()}}
	c.Request = httptest.NewRequest("POST", "/api/projects/"+projectID.String()+"/finetunes", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateFinetune(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func intPtr(i int) *int {
	return &i
}