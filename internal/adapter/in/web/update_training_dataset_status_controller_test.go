package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
)

type mockUpdateTrainingDatasetStatusUseCase struct {
	executeError error
}

func (m *mockUpdateTrainingDatasetStatusUseCase) Execute(ctx context.Context, command in.UpdateTrainingDatasetStatusCommand) error {
	if command.TrainingDatasetID.String() == "00000000-0000-0000-0000-000000000000" {
		return fmt.Errorf("training dataset not found")
	}
	if command.Status == "INVALID" {
		return fmt.Errorf("invalid status transition from PLANNING to INVALID")
	}
	return m.executeError
}

func TestUpdateTrainingDatasetStatusController_UpdateStatus_Success(t *testing.T) {
	mockUseCase := &mockUpdateTrainingDatasetStatusUseCase{}
	controller := &UpdateTrainingDatasetStatusController{
		UpdateTrainingDatasetStatusUseCase: mockUseCase,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/api/external/training-datasets/:training_dataset_id/update-status", controller.UpdateStatus)

	request := UpdateTrainingDatasetStatusRequest{
		Status: entities.TrainingDatasetStatusRunning,
	}
	requestBody, _ := json.Marshal(request)

	trainingDatasetID := uuid.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/external/training-datasets/"+trainingDatasetID.String()+"/update-status", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Training dataset status updated successfully", response["message"])
}

func TestUpdateTrainingDatasetStatusController_UpdateStatus_InvalidID(t *testing.T) {
	mockUseCase := &mockUpdateTrainingDatasetStatusUseCase{}
	controller := &UpdateTrainingDatasetStatusController{
		UpdateTrainingDatasetStatusUseCase: mockUseCase,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/api/external/training-datasets/:training_dataset_id/update-status", controller.UpdateStatus)

	request := UpdateTrainingDatasetStatusRequest{
		Status: entities.TrainingDatasetStatusRunning,
	}
	requestBody, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/external/training-datasets/invalid-id/update-status", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid training dataset ID format", response["error"])
}

func TestUpdateTrainingDatasetStatusController_UpdateStatus_NotFound(t *testing.T) {
	mockUseCase := &mockUpdateTrainingDatasetStatusUseCase{}
	controller := &UpdateTrainingDatasetStatusController{
		UpdateTrainingDatasetStatusUseCase: mockUseCase,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/api/external/training-datasets/:training_dataset_id/update-status", controller.UpdateStatus)

	request := UpdateTrainingDatasetStatusRequest{
		Status: entities.TrainingDatasetStatusRunning,
	}
	requestBody, _ := json.Marshal(request)

	// Use a zero UUID to trigger "not found" error in mock
	trainingDatasetID := uuid.UUID{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/external/training-datasets/"+trainingDatasetID.String()+"/update-status", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Training dataset not found", response["error"])
}

func TestUpdateTrainingDatasetStatusController_UpdateStatus_InvalidRequest(t *testing.T) {
	mockUseCase := &mockUpdateTrainingDatasetStatusUseCase{}
	controller := &UpdateTrainingDatasetStatusController{
		UpdateTrainingDatasetStatusUseCase: mockUseCase,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/api/external/training-datasets/:training_dataset_id/update-status", controller.UpdateStatus)

	// Invalid JSON request
	requestBody := []byte(`{"status": }`)

	trainingDatasetID := uuid.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/external/training-datasets/"+trainingDatasetID.String()+"/update-status", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid request format")
}