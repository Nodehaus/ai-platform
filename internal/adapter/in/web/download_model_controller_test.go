package web

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

// mockDownloadModelUseCase is a mock implementation of DownloadModelUseCase
type mockDownloadModelUseCase struct {
	reader        io.ReadCloser
	contentLength int64
	filename      string
	err           error
}

func (m *mockDownloadModelUseCase) DownloadModel(ctx context.Context, command in.DownloadModelCommand) (io.ReadCloser, int64, string, error) {
	return m.reader, m.contentLength, m.filename, m.err
}

func TestDownloadModelController_DownloadModel_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock reader
	mockContent := "mock file content"
	mockReader := io.NopCloser(strings.NewReader(mockContent))

	mockUseCase := &mockDownloadModelUseCase{
		reader:        mockReader,
		contentLength: int64(len(mockContent)),
		filename:      "test-model.gguf",
		err:           nil,
	}

	controller := &DownloadModelController{
		DownloadModelUseCase: mockUseCase,
	}

	projectID := uuid.New()
	finetuneID := uuid.New()

	router := gin.New()
	router.GET("/api/projects/:project_id/finetunes/:finetune_id/download", controller.DownloadModel)

	url := "/api/projects/" + projectID.String() + "/finetunes/" + finetuneID.String() + "/download"
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedDisposition := "attachment; filename=\"test-model.gguf\""
	if w.Header().Get("Content-Disposition") != expectedDisposition {
		t.Errorf("Expected Content-Disposition %s, got %s", expectedDisposition, w.Header().Get("Content-Disposition"))
	}

	if w.Header().Get("Content-Type") != "application/octet-stream" {
		t.Errorf("Expected Content-Type application/octet-stream, got %s", w.Header().Get("Content-Type"))
	}

	if w.Header().Get("Content-Length") != "17" {
		t.Errorf("Expected Content-Length 17, got %s", w.Header().Get("Content-Length"))
	}

	if w.Body.String() != mockContent {
		t.Errorf("Expected body %s, got %s", mockContent, w.Body.String())
	}
}


func TestDownloadModelController_DownloadModel_InvalidProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := &DownloadModelController{}

	finetuneID := uuid.New()

	router := gin.New()
	router.GET("/api/projects/:project_id/finetunes/:finetune_id/download", controller.DownloadModel)

	url := "/api/projects/invalid-uuid/finetunes/" + finetuneID.String() + "/download"
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Invalid project_id format") {
		t.Errorf("Expected error message about invalid project_id, got %s", w.Body.String())
	}
}

func TestDownloadModelController_DownloadModel_InvalidFinetuneID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := &DownloadModelController{}

	projectID := uuid.New()

	router := gin.New()
	router.GET("/api/projects/:project_id/finetunes/:finetune_id/download", controller.DownloadModel)

	url := "/api/projects/" + projectID.String() + "/finetunes/invalid-uuid/download"
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Invalid finetune_id format") {
		t.Errorf("Expected error message about invalid finetune_id, got %s", w.Body.String())
	}
}

func TestDownloadModelController_DownloadModel_UseCaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockDownloadModelUseCase{
		reader:        nil,
		contentLength: 0,
		filename:      "",
		err:           errors.New("finetune not found"),
	}

	controller := &DownloadModelController{
		DownloadModelUseCase: mockUseCase,
	}

	projectID := uuid.New()
	finetuneID := uuid.New()

	router := gin.New()
	router.GET("/api/projects/:project_id/finetunes/:finetune_id/download", controller.DownloadModel)

	url := "/api/projects/" + projectID.String() + "/finetunes/" + finetuneID.String() + "/download"
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Body.String(), "finetune not found") {
		t.Errorf("Expected error message about finetune not found, got %s", w.Body.String())
	}
}