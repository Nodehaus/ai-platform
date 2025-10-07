package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
)

type mockCreateDeploymentUseCase struct {
	result *in.CreateDeploymentResult
	err    error
}

func (m *mockCreateDeploymentUseCase) CreateDeployment(command in.CreateDeploymentCommand) (*in.CreateDeploymentResult, error) {
	if command.ModelName == "" {
		return nil, errors.New("model_name is required")
	}
	if command.ModelName == "invalid-project" {
		return nil, errors.New("project not found")
	}
	if command.ModelName == "invalid-finetune" {
		return nil, errors.New("finetune not found")
	}
	return m.result, m.err
}

func TestCreateDeploymentController_CreateDeployment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()
	deployment := &entities.Deployment{
		ID:         uuid.New(),
		ModelName:  "test-model",
		APIKey:     "sk-test-key",
		ProjectID:  projectID,
		FinetuneID: nil,
	}

	result := &in.CreateDeploymentResult{
		Deployment: deployment,
	}

	mockUseCase := &mockCreateDeploymentUseCase{
		result: result,
		err:    nil,
	}

	controller := &CreateDeploymentController{CreateDeploymentUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects/:project_id/deployments", controller.CreateDeployment)

	request := CreateDeploymentRequest{
		ModelName:  "test-model",
		FinetuneID: nil,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/deployments", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response CreateDeploymentResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response.ModelName != "test-model" {
		t.Errorf("Expected model name 'test-model', got %s", response.ModelName)
	}

	if response.ID == uuid.Nil {
		t.Error("Expected deployment ID to be set")
	}

	if response.APIKey == "" {
		t.Error("Expected API key to be set")
	}
}

func TestCreateDeploymentController_CreateDeployment_WithFinetuneID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()
	finetuneID := uuid.New()
	deployment := &entities.Deployment{
		ID:         uuid.New(),
		ModelName:  "test-model",
		APIKey:     "sk-test-key",
		ProjectID:  projectID,
		FinetuneID: &finetuneID,
	}

	result := &in.CreateDeploymentResult{
		Deployment: deployment,
	}

	mockUseCase := &mockCreateDeploymentUseCase{
		result: result,
		err:    nil,
	}

	controller := &CreateDeploymentController{CreateDeploymentUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects/:project_id/deployments", controller.CreateDeployment)

	request := CreateDeploymentRequest{
		ModelName:  "test-model",
		FinetuneID: &finetuneID,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/deployments", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response CreateDeploymentResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response.ModelName != "test-model" {
		t.Errorf("Expected model name 'test-model', got %s", response.ModelName)
	}
}

func TestCreateDeploymentController_CreateDeployment_NoUserInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockCreateDeploymentUseCase{}
	controller := &CreateDeploymentController{CreateDeploymentUseCase: mockUseCase}

	router := gin.New()
	router.POST("/projects/:project_id/deployments", controller.CreateDeployment)

	projectID := uuid.New()
	request := CreateDeploymentRequest{
		ModelName: "test-model",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/deployments", bytes.NewBuffer(requestBody))
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

func TestCreateDeploymentController_CreateDeployment_InvalidProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockUseCase := &mockCreateDeploymentUseCase{}
	controller := &CreateDeploymentController{CreateDeploymentUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects/:project_id/deployments", controller.CreateDeployment)

	request := CreateDeploymentRequest{
		ModelName: "test-model",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/invalid-uuid/deployments", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "Invalid project ID" {
		t.Errorf("Expected error 'Invalid project ID', got %s", response["error"])
	}
}

func TestCreateDeploymentController_CreateDeployment_MissingModelName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	projectID := uuid.New()
	mockUseCase := &mockCreateDeploymentUseCase{}
	controller := &CreateDeploymentController{CreateDeploymentUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects/:project_id/deployments", controller.CreateDeployment)

	request := CreateDeploymentRequest{
		ModelName: "",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/deployments", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}
