package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
)

type mockCreateProjectUseCase struct {
	result *in.CreateProjectResult
	err    error
}

func (m *mockCreateProjectUseCase) CreateProject(command in.CreateProjectCommand) (*in.CreateProjectResult, error) {
	// Return error if name starts with "invalid"
	if strings.HasPrefix(command.Name, "invalid") {
		return nil, errors.New("project name cannot be empty")
	}
	// Return error if name is "duplicate"
	if command.Name == "duplicate" {
		return nil, errors.New("project name already exists")
	}
	return m.result, m.err
}

func TestCreateProjectController_CreateProject_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	project := &entities.Project{
		ID:      uuid.New(),
		Name:    "Test Project",
		OwnerID: userID,
		Status:  entities.ProjectStatusActive,
	}

	result := &in.CreateProjectResult{
		Project: project,
	}

	mockUseCase := &mockCreateProjectUseCase{
		result: result,
		err:    nil,
	}

	controller := &CreateProjectController{CreateProjectUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects", controller.CreateProject)

	request := CreateProjectRequest{
		Name: "Test Project",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response CreateProjectResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response.Project.Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got %s", response.Project.Name)
	}

	if response.Message != "Project created successfully" {
		t.Errorf("Expected message 'Project created successfully', got %s", response.Message)
	}
}

func TestCreateProjectController_CreateProject_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockUseCase := &mockCreateProjectUseCase{}

	controller := &CreateProjectController{CreateProjectUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects", controller.CreateProject)

	request := CreateProjectRequest{
		Name: "invalid-project", // Will trigger the validation error in mock
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "project name cannot be empty" {
		t.Errorf("Expected error 'project name cannot be empty', got %s", response["error"])
	}
}

func TestCreateProjectController_CreateProject_DuplicateName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockUseCase := &mockCreateProjectUseCase{}

	controller := &CreateProjectController{CreateProjectUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/projects", controller.CreateProject)

	request := CreateProjectRequest{
		Name: "duplicate",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "project name already exists" {
		t.Errorf("Expected error 'project name already exists', got %s", response["error"])
	}
}

func TestCreateProjectController_CreateProject_NoUserInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockCreateProjectUseCase{}
	controller := &CreateProjectController{CreateProjectUseCase: mockUseCase}

	router := gin.New()
	router.POST("/projects", controller.CreateProject)

	request := CreateProjectRequest{
		Name: "Test Project",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(requestBody))
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