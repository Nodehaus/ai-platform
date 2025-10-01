package web

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"ai-platform/internal/application/port/in"
)

type MockDownloadTrainingDatasetUseCase struct {
	DownloadTrainingDatasetFunc func(command in.DownloadTrainingDatasetCommand) (*in.DownloadTrainingDatasetResult, error)
}

func (m *MockDownloadTrainingDatasetUseCase) DownloadTrainingDataset(command in.DownloadTrainingDatasetCommand) (*in.DownloadTrainingDatasetResult, error) {
	return m.DownloadTrainingDatasetFunc(command)
}

func TestDownloadTrainingDatasetController_DownloadTrainingDataset_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	projectID := uuid.New()
	trainingDatasetID := uuid.New()
	userID := uuid.New()

	fieldNames := []string{"input", "output", "category"}
	data := [][]string{
		{"What is AI?", "Artificial Intelligence", "Technology"},
		{"What is ML?", "Machine Learning", "Technology"},
	}

	mockUseCase := &MockDownloadTrainingDatasetUseCase{
		DownloadTrainingDatasetFunc: func(command in.DownloadTrainingDatasetCommand) (*in.DownloadTrainingDatasetResult, error) {
			assert.Equal(t, projectID, command.ProjectID)
			assert.Equal(t, trainingDatasetID, command.TrainingDatasetID)
			assert.Equal(t, userID, command.OwnerID)

			return &in.DownloadTrainingDatasetResult{
				FieldNames: fieldNames,
				Data:       data,
				Filename:   "dataset_test_project_v1.csv",
			}, nil
		},
	}

	controller := &DownloadTrainingDatasetController{
		DownloadTrainingDatasetUseCase: mockUseCase,
	}

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set params
	c.Params = gin.Params{
		{Key: "project_id", Value: projectID.String()},
		{Key: "training_dataset_id", Value: trainingDatasetID.String()},
	}

	// Set user ID in context
	c.Set("user_id", userID)

	// Execute
	controller.DownloadTrainingDataset(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Disposition"), "dataset_test_project_v1.csv")

	// Parse CSV response
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)

	// Should have header + 2 data rows
	assert.Equal(t, 3, len(records))

	// Check header
	assert.Equal(t, fieldNames, records[0])

	// Check first data row
	assert.Equal(t, data[0], records[1])

	// Check second data row
	assert.Equal(t, data[1], records[2])
}

func TestDownloadTrainingDatasetController_DownloadTrainingDataset_InvalidProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &MockDownloadTrainingDatasetUseCase{}
	controller := &DownloadTrainingDatasetController{
		DownloadTrainingDatasetUseCase: mockUseCase,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		{Key: "project_id", Value: "invalid-uuid"},
		{Key: "training_dataset_id", Value: uuid.New().String()},
	}
	c.Set("user_id", uuid.New())

	controller.DownloadTrainingDataset(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadTrainingDatasetController_DownloadTrainingDataset_InvalidTrainingDatasetID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &MockDownloadTrainingDatasetUseCase{}
	controller := &DownloadTrainingDatasetController{
		DownloadTrainingDatasetUseCase: mockUseCase,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		{Key: "project_id", Value: uuid.New().String()},
		{Key: "training_dataset_id", Value: "invalid-uuid"},
	}
	c.Set("user_id", uuid.New())

	controller.DownloadTrainingDataset(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadTrainingDatasetController_DownloadTrainingDataset_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &MockDownloadTrainingDatasetUseCase{
		DownloadTrainingDatasetFunc: func(command in.DownloadTrainingDatasetCommand) (*in.DownloadTrainingDatasetResult, error) {
			return nil, nil
		},
	}

	controller := &DownloadTrainingDatasetController{
		DownloadTrainingDatasetUseCase: mockUseCase,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		{Key: "project_id", Value: uuid.New().String()},
		{Key: "training_dataset_id", Value: uuid.New().String()},
	}
	c.Set("user_id", uuid.New())

	controller.DownloadTrainingDataset(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDownloadTrainingDatasetController_DownloadTrainingDataset_EmptyData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	projectID := uuid.New()
	trainingDatasetID := uuid.New()
	userID := uuid.New()

	fieldNames := []string{"input", "output"}

	mockUseCase := &MockDownloadTrainingDatasetUseCase{
		DownloadTrainingDatasetFunc: func(command in.DownloadTrainingDatasetCommand) (*in.DownloadTrainingDatasetResult, error) {
			return &in.DownloadTrainingDatasetResult{
				FieldNames: fieldNames,
				Data:       [][]string{}, // Empty data
				Filename:   "dataset_empty_v1.csv",
			}, nil
		},
	}

	controller := &DownloadTrainingDatasetController{
		DownloadTrainingDatasetUseCase: mockUseCase,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		{Key: "project_id", Value: projectID.String()},
		{Key: "training_dataset_id", Value: trainingDatasetID.String()},
	}
	c.Set("user_id", userID)

	controller.DownloadTrainingDataset(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse CSV response
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)

	// Should only have header row
	assert.Equal(t, 1, len(records))
	assert.Equal(t, fieldNames, records[0])
}
