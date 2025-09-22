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

type mockCreateTrainingDatasetUseCase struct {
	result *entities.TrainingDataset
	err    error
}

func (m *mockCreateTrainingDatasetUseCase) Execute(ctx context.Context, command in.CreateTrainingDatasetCommand) (*entities.TrainingDataset, error) {
	if command.CorpusName == "invalid" {
		return nil, errors.New("corpus_name is required")
	}
	if command.InputField == "nonexistent" {
		return nil, errors.New("input_field must be present in field_names")
	}
	return m.result, m.err
}


func TestCreateTrainingDatasetController_CreateTrainingDataset_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()
	promptID := uuid.New()

	trainingDataset := &entities.TrainingDataset{
		ID:                       uuid.New(),
		ProjectID:                projectID,
		Version:                  1,
		InputField:               "question",
		OutputField:              "answer",
		GeneratePromptHistoryIDs: []uuid.UUID{},
		GeneratePromptID:         promptID,
		CorpusID:                 uuid.New(),
		LanguageISO:              "deu",
		Status:                   entities.TrainingDatasetStatusPlanning,
		FieldNames:               []string{"question", "answer", "complexity"},
		GenerateExamplesNumber:   100,
		Data:                     []entities.TrainingDataItem{},
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	mockUseCase := &mockCreateTrainingDatasetUseCase{
		result: trainingDataset,
		err:    nil,
	}

	controller := &CreateTrainingDatasetController{CreateTrainingDatasetUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects/:project_id/training-datasets", controller.CreateTrainingDataset)

	request := CreateTrainingDatasetRequest{
		CorpusName:            "eurlex",
		InputField:            "question",
		OutputField:           "answer",
		LanguageISO:           "deu",
		FieldNames:            []string{"question", "answer", "complexity"},
		GeneratePrompt:        "Generate Q&A dataset",
		GenerateExamplesNumber: 100,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/training-datasets", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response CreateTrainingDatasetResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response.ID == uuid.Nil {
		t.Error("Expected response to contain a valid ID")
	}

	if response.ProjectID != projectID {
		t.Errorf("Expected project ID '%s', got '%s'", projectID, response.ProjectID)
	}
}

func TestCreateTrainingDatasetController_CreateTrainingDataset_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()
	mockUseCase := &mockCreateTrainingDatasetUseCase{}

	controller := &CreateTrainingDatasetController{CreateTrainingDatasetUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects/:project_id/training-datasets", controller.CreateTrainingDataset)

	request := CreateTrainingDatasetRequest{
		CorpusName:            "invalid", // Will trigger the validation error in mock
		InputField:            "question",
		OutputField:           "answer",
		LanguageISO:           "deu",
		FieldNames:            []string{"question", "answer"},
		GeneratePrompt:        "Generate Q&A dataset",
		GenerateExamplesNumber: 75,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/training-datasets", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "corpus_name is required" {
		t.Errorf("Expected error 'corpus_name is required', got %s", response["error"])
	}
}

func TestCreateTrainingDatasetController_CreateTrainingDataset_FieldValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()
	mockUseCase := &mockCreateTrainingDatasetUseCase{}

	controller := &CreateTrainingDatasetController{CreateTrainingDatasetUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects/:project_id/training-datasets", controller.CreateTrainingDataset)

	request := CreateTrainingDatasetRequest{
		CorpusName:            "eurlex",
		InputField:            "nonexistent", // Will trigger field validation error in mock
		OutputField:           "answer",
		LanguageISO:           "deu",
		FieldNames:            []string{"question", "answer"},
		GeneratePrompt:        "Generate Q&A dataset",
		GenerateExamplesNumber: 50,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/training-datasets", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "input_field must be present in field_names" {
		t.Errorf("Expected error 'input_field must be present in field_names', got %s", response["error"])
	}
}

func TestCreateTrainingDatasetController_CreateTrainingDataset_InvalidProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockUseCase := &mockCreateTrainingDatasetUseCase{}

	controller := &CreateTrainingDatasetController{CreateTrainingDatasetUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects/:project_id/training-datasets", controller.CreateTrainingDataset)

	request := CreateTrainingDatasetRequest{
		CorpusName:            "eurlex",
		InputField:            "question",
		OutputField:           "answer",
		LanguageISO:           "deu",
		FieldNames:            []string{"question", "answer"},
		GeneratePrompt:        "Generate Q&A dataset",
		GenerateExamplesNumber: 25,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/invalid-id/training-datasets", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "Invalid project ID format" {
		t.Errorf("Expected error 'Invalid project ID format', got %s", response["error"])
	}
}

func TestCreateTrainingDatasetController_CreateTrainingDataset_NoUserInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	projectID := uuid.New()
	mockUseCase := &mockCreateTrainingDatasetUseCase{}
	controller := &CreateTrainingDatasetController{CreateTrainingDatasetUseCase: mockUseCase}

	router := gin.New()
	router.POST("/projects/:project_id/training-datasets", controller.CreateTrainingDataset)

	request := CreateTrainingDatasetRequest{
		CorpusName:            "eurlex",
		InputField:            "question",
		OutputField:           "answer",
		LanguageISO:           "deu",
		FieldNames:            []string{"question", "answer"},
		GeneratePrompt:        "Generate Q&A dataset",
		GenerateExamplesNumber: 25,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/training-datasets", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "User ID not found in context" {
		t.Errorf("Expected error 'User ID not found in context', got %s", response["error"])
	}
}