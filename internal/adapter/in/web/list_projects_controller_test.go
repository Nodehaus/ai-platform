package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
)

type mockListProjectsUseCase struct {
	result *in.ListProjectsResult
	err    error
}

func (m *mockListProjectsUseCase) ListProjects(command in.ListProjectsCommand) (*in.ListProjectsResult, error) {
	return m.result, m.err
}

func TestListProjectsController_ListProjects_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	trainingDatasetID := uuid.New()
	projectsWithTrainingDatasets := []in.ProjectWithTrainingDataset{
		{
			Project: entities.Project{
				ID:        uuid.New(),
				Name:      "Project 1",
				OwnerID:   userID,
				Status:    entities.ProjectStatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			TrainingDatasetID: &trainingDatasetID,
		},
		{
			Project: entities.Project{
				ID:        uuid.New(),
				Name:      "Project 2",
				OwnerID:   userID,
				Status:    entities.ProjectStatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			TrainingDatasetID: nil,
		},
	}

	result := &in.ListProjectsResult{
		Projects: projectsWithTrainingDatasets,
	}

	mockUseCase := &mockListProjectsUseCase{
		result: result,
		err:    nil,
	}

	controller := &ListProjectsController{ListProjectsUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.GET("/projects", controller.ListProjects)

	req, _ := http.NewRequest("GET", "/projects", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response ListProjectsResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if len(response.Projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(response.Projects))
	}

	if response.Projects[0].Name != "Project 1" {
		t.Errorf("Expected first project name 'Project 1', got %s", response.Projects[0].Name)
	}

	if response.Projects[1].Name != "Project 2" {
		t.Errorf("Expected second project name 'Project 2', got %s", response.Projects[1].Name)
	}

	if response.Projects[0].TrainingDatasetID == nil {
		t.Error("Expected first project to have a training dataset ID")
	}

	if response.Projects[1].TrainingDatasetID != nil {
		t.Error("Expected second project to have no training dataset ID")
	}
}

func TestListProjectsController_ListProjects_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	result := &in.ListProjectsResult{
		Projects: []in.ProjectWithTrainingDataset{},
	}

	mockUseCase := &mockListProjectsUseCase{
		result: result,
		err:    nil,
	}

	controller := &ListProjectsController{ListProjectsUseCase: mockUseCase}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.GET("/projects", controller.ListProjects)

	req, _ := http.NewRequest("GET", "/projects", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response ListProjectsResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if len(response.Projects) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(response.Projects))
	}
}

func TestListProjectsController_ListProjects_NoUserInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockListProjectsUseCase{}
	controller := &ListProjectsController{ListProjectsUseCase: mockUseCase}

	router := gin.New()
	router.GET("/projects", controller.ListProjects)

	req, _ := http.NewRequest("GET", "/projects", nil)
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