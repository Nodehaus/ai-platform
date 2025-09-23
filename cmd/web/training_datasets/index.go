package training_datasets

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

type TrainingDatasetIndexData struct {
	ProjectName      string
	TrainingDataset  TrainingDatasetData
	TotalDataItems   int
}

type TrainingDatasetData struct {
	Version                int         `json:"version"`
	GeneratePrompt         string      `json:"generate_prompt"`
	InputField             string      `json:"input_field"`
	OutputField            string      `json:"output_field"`
	GenerateExamplesNumber int         `json:"generate_examples_number"`
	CorpusName             string      `json:"corpus_name"`
	LanguageISO            string      `json:"language_iso"`
	Status                 string      `json:"status"`
	FieldNames             []string    `json:"field_names"`
	DataItemsSample        [][]string  `json:"data_items_sample"`
}

func TrainingDatasetIndexHandler(w http.ResponseWriter, r *http.Request) {
	token := web.GetTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	// Extract project ID and training dataset ID from URL path
	// Expected format: /web/projects/{project_id}/training-datasets/{training_dataset_id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 || pathParts[3] == "" || pathParts[5] == "" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	projectIDStr := pathParts[3]
	trainingDatasetIDStr := pathParts[5]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID format", http.StatusBadRequest)
		return
	}

	trainingDatasetID, err := uuid.Parse(trainingDatasetIDStr)
	if err != nil {
		http.Error(w, "Invalid training dataset ID format", http.StatusBadRequest)
		return
	}

	// Fetch training dataset data
	trainingDatasetData, err := fetchTrainingDatasetData(r, token, projectID, trainingDatasetID)
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

	// Calculate total data items (for display purposes, we use the GenerateExamplesNumber)
	totalDataItems := trainingDatasetData.GenerateExamplesNumber

	indexData := TrainingDatasetIndexData{
		ProjectName:     projectName,
		TrainingDataset: *trainingDatasetData,
		TotalDataItems:  totalDataItems,
	}

	templ.Handler(TrainingDatasetIndex(indexData)).ServeHTTP(w, r)
}

func fetchTrainingDatasetData(r *http.Request, token string, projectID uuid.UUID, trainingDatasetID uuid.UUID) (*TrainingDatasetData, error) {
	apiBaseURL := web.GetAPIBaseURL(r)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/%s/training-datasets/%s", apiBaseURL, projectID, trainingDatasetID), nil)
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

	var trainingDataset TrainingDatasetData
	if err := json.NewDecoder(resp.Body).Decode(&trainingDataset); err != nil {
		return nil, err
	}

	return &trainingDataset, nil
}