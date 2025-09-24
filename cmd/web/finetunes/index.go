package finetunes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"

	"ai-platform/cmd/web"
)

type FinetuneIndexData struct {
	ProjectID   string
	ProjectName string
	FinetuneID  string
	Finetune    FinetuneData
}

type FinetuneData struct {
	ID                               uuid.UUID              `json:"id"`
	Version                          int                    `json:"version"`
	Status                           string                 `json:"status"`
	BaseModelName                    string                 `json:"base_model_name"`
	TrainingDatasetID                uuid.UUID             `json:"training_dataset_id"`
	TrainingDatasetNumberExamples    *int                   `json:"training_dataset_number_examples"`
	TrainingDatasetSelectRandom      bool                   `json:"training_dataset_select_random"`
	ModelSizeGB                      *int                   `json:"model_size_gb"`
	ModelSizeParameter               *int                   `json:"model_size_parameter"`
	ModelDtype                       *string                `json:"model_dtype"`
	ModelQuantization                *string                `json:"model_quantization"`
	InferenceSamples                 []InferenceSample      `json:"inference_samples"`
	TrainingTimeSeconds              *float64               `json:"training_time_seconds"`
}

type InferenceSample struct {
	AtStep int                   `json:"at_step"`
	Items  []InferenceSampleItem `json:"items"`
}

type InferenceSampleItem struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func FinetuneIndexHandler(w http.ResponseWriter, r *http.Request) {
	token := web.GetTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	// Extract project ID and finetune ID from URL path
	// Expected format: /web/projects/{project_id}/finetunes/{finetune_id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 || pathParts[3] == "" || pathParts[5] == "" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	projectIDStr := pathParts[3]
	finetuneIDStr := pathParts[5]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID format", http.StatusBadRequest)
		return
	}

	finetuneID, err := uuid.Parse(finetuneIDStr)
	if err != nil {
		http.Error(w, "Invalid finetune ID format", http.StatusBadRequest)
		return
	}

	// Fetch finetune data
	finetuneData, err := fetchFinetuneData(r, token, projectID, finetuneID)
	if err != nil {
		// If we can't fetch the data, redirect to login (token might be invalid)
		web.ClearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	// Fetch project details to get the project name
	projectName, err := fetchProjectName(r, token, projectID)
	if err != nil {
		web.ClearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	indexData := FinetuneIndexData{
		ProjectID:   projectIDStr,
		ProjectName: projectName,
		FinetuneID:  finetuneIDStr,
		Finetune:    *finetuneData,
	}

	templ.Handler(FinetuneIndex(indexData)).ServeHTTP(w, r)
}

func fetchFinetuneData(r *http.Request, token string, projectID uuid.UUID, finetuneID uuid.UUID) (*FinetuneData, error) {
	apiBaseURL := web.GetAPIBaseURL(r)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/%s/finetunes/%s", apiBaseURL, projectID, finetuneID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var finetune FinetuneData
	if err := json.NewDecoder(resp.Body).Decode(&finetune); err != nil {
		return nil, err
	}

	return &finetune, nil
}

func fetchProjectName(r *http.Request, token string, projectID uuid.UUID) (string, error) {
	apiBaseURL := web.GetAPIBaseURL(r)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/%s", apiBaseURL, projectID), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var project struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return "", err
	}

	return project.Name, nil
}