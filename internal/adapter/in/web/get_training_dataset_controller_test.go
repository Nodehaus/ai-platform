package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
)

type MockGetTrainingDatasetUseCase struct {
	mock.Mock
}

func (m *MockGetTrainingDatasetUseCase) GetTrainingDataset(command in.GetTrainingDatasetCommand) (*in.GetTrainingDatasetResult, error) {
	args := m.Called(command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*in.GetTrainingDatasetResult), args.Error(1)
}

func TestGetTrainingDatasetController_GetTrainingDataset_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetTrainingDatasetUseCase)
	controller := &GetTrainingDatasetController{
		GetTrainingDatasetUseCase: mockUseCase,
	}

	projectID := uuid.New()
	trainingDatasetID := uuid.New()
	userID := uuid.New()

	trainingDataset := &entities.TrainingDataset{
		ID:                     trainingDatasetID,
		ProjectID:              projectID,
		Version:                1,
		InputField:             "input",
		OutputField:            "output",
		GenerateExamplesNumber: 10,
		LanguageISO:            "en",
		Status:                 entities.TrainingDatasetStatusDone,
		FieldNames:             []string{"input", "output"},
		Data: []entities.TrainingDataItem{
			{
				ID:         uuid.New(),
				Values:     []string{"test input", "test output"},
				CorrectsID: nil,
				Deleted:    false,
			},
		},
	}

	expectedCommand := in.GetTrainingDatasetCommand{
		ProjectID:         projectID,
		TrainingDatasetID: trainingDatasetID,
		OwnerID:           userID,
	}

	expectedResult := &in.GetTrainingDatasetResult{
		TrainingDataset: trainingDataset,
		GeneratePrompt:  "Test prompt",
		CorpusName:      "Test corpus",
	}

	mockUseCase.On("GetTrainingDataset", expectedCommand).Return(expectedResult, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/projects/"+projectID.String()+"/training-datasets/"+trainingDatasetID.String(), nil)
	c.Params = []gin.Param{
		{Key: "project_id", Value: projectID.String()},
		{Key: "training_dataset_id", Value: trainingDatasetID.String()},
	}
	c.Set("user_id", userID)

	controller.GetTrainingDataset(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetTrainingDatasetResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 1, response.Version)
	assert.Equal(t, "Test prompt", response.GeneratePrompt)
	assert.Equal(t, "input", response.InputField)
	assert.Equal(t, "output", response.OutputField)
	assert.Equal(t, 10, response.GenerateExamplesNumber)
	assert.Equal(t, "Test corpus", response.CorpusName)
	assert.Equal(t, "en", response.LanguageISO)
	assert.Equal(t, "DONE", response.Status)
	assert.Equal(t, []string{"input", "output"}, response.FieldNames)
	assert.Len(t, response.DataItemsSample, 1)
	assert.Equal(t, []string{"test input", "test output"}, response.DataItemsSample[0])

	mockUseCase.AssertExpectations(t)
}

func TestGetTrainingDatasetController_GetTrainingDataset_NotDoneStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetTrainingDatasetUseCase)
	controller := &GetTrainingDatasetController{
		GetTrainingDatasetUseCase: mockUseCase,
	}

	projectID := uuid.New()
	trainingDatasetID := uuid.New()
	userID := uuid.New()

	trainingDataset := &entities.TrainingDataset{
		ID:                     trainingDatasetID,
		ProjectID:              projectID,
		Version:                1,
		InputField:             "input",
		OutputField:            "output",
		GenerateExamplesNumber: 10,
		LanguageISO:            "en",
		Status:                 entities.TrainingDatasetStatusRunning,
		FieldNames:             []string{"input", "output"},
		Data: []entities.TrainingDataItem{
			{
				ID:         uuid.New(),
				Values:     []string{"test input", "test output"},
				CorrectsID: nil,
				Deleted:    false,
			},
		},
	}

	expectedCommand := in.GetTrainingDatasetCommand{
		ProjectID:         projectID,
		TrainingDatasetID: trainingDatasetID,
		OwnerID:           userID,
	}

	expectedResult := &in.GetTrainingDatasetResult{
		TrainingDataset: trainingDataset,
		GeneratePrompt:  "Test prompt",
		CorpusName:      "Test corpus",
	}

	mockUseCase.On("GetTrainingDataset", expectedCommand).Return(expectedResult, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/projects/"+projectID.String()+"/training-datasets/"+trainingDatasetID.String(), nil)
	c.Params = []gin.Param{
		{Key: "project_id", Value: projectID.String()},
		{Key: "training_dataset_id", Value: trainingDatasetID.String()},
	}
	c.Set("user_id", userID)

	controller.GetTrainingDataset(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetTrainingDatasetResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "RUNNING", response.Status)
	assert.Empty(t, response.DataItemsSample) // Should be empty when not DONE

	mockUseCase.AssertExpectations(t)
}

func TestGetTrainingDatasetController_GetTrainingDataset_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetTrainingDatasetUseCase)
	controller := &GetTrainingDatasetController{
		GetTrainingDatasetUseCase: mockUseCase,
	}

	projectID := uuid.New()
	trainingDatasetID := uuid.New()
	userID := uuid.New()

	expectedCommand := in.GetTrainingDatasetCommand{
		ProjectID:         projectID,
		TrainingDatasetID: trainingDatasetID,
		OwnerID:           userID,
	}

	mockUseCase.On("GetTrainingDataset", expectedCommand).Return(nil, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/projects/"+projectID.String()+"/training-datasets/"+trainingDatasetID.String(), nil)
	c.Params = []gin.Param{
		{Key: "project_id", Value: projectID.String()},
		{Key: "training_dataset_id", Value: trainingDatasetID.String()},
	}
	c.Set("user_id", userID)

	controller.GetTrainingDataset(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Training dataset not found", response["error"])

	mockUseCase.AssertExpectations(t)
}

func TestGetTrainingDatasetController_GetTrainingDataset_InvalidProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetTrainingDatasetUseCase)
	controller := &GetTrainingDatasetController{
		GetTrainingDatasetUseCase: mockUseCase,
	}

	trainingDatasetID := uuid.New()
	userID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/projects/invalid-id/training-datasets/"+trainingDatasetID.String(), nil)
	c.Params = []gin.Param{
		{Key: "project_id", Value: "invalid-id"},
		{Key: "training_dataset_id", Value: trainingDatasetID.String()},
	}
	c.Set("user_id", userID)

	controller.GetTrainingDataset(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid project ID format", response["error"])
}

func TestGetTrainingDatasetController_GetTrainingDataset_InvalidTrainingDatasetID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetTrainingDatasetUseCase)
	controller := &GetTrainingDatasetController{
		GetTrainingDatasetUseCase: mockUseCase,
	}

	projectID := uuid.New()
	userID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/projects/"+projectID.String()+"/training-datasets/invalid-id", nil)
	c.Params = []gin.Param{
		{Key: "project_id", Value: projectID.String()},
		{Key: "training_dataset_id", Value: "invalid-id"},
	}
	c.Set("user_id", userID)

	controller.GetTrainingDataset(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid training dataset ID format", response["error"])
}

func TestGetTrainingDatasetController_GetTrainingDataset_NoUserInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetTrainingDatasetUseCase)
	controller := &GetTrainingDatasetController{
		GetTrainingDatasetUseCase: mockUseCase,
	}

	projectID := uuid.New()
	trainingDatasetID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/projects/"+projectID.String()+"/training-datasets/"+trainingDatasetID.String(), nil)
	c.Params = []gin.Param{
		{Key: "project_id", Value: projectID.String()},
		{Key: "training_dataset_id", Value: trainingDatasetID.String()},
	}
	// Not setting user_id in context

	controller.GetTrainingDataset(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User ID not found in context", response["error"])
}