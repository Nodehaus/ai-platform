package training_datasets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"

	"ai-platform/cmd/web"
)

type TrainingDatasetIndexData struct {
	ProjectID           string
	ProjectName         string
	TrainingDatasetID   string
	TrainingDataset     TrainingDatasetData
	TotalDataItems      int
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
		ProjectID:           projectIDStr,
		ProjectName:         projectName,
		TrainingDatasetID:   trainingDatasetIDStr,
		TrainingDataset:     *trainingDatasetData,
		TotalDataItems:      totalDataItems,
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

func CreateFinetuneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := web.GetTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	// Extract project ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 || pathParts[3] == "" {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	projectIDStr := pathParts[3]
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID format", http.StatusBadRequest)
		return
	}

	// Get form values
	baseModel := r.FormValue("base-model")
	examplesCountStr := r.FormValue("examples-count")
	randomSelection := r.FormValue("random-selection") == "on"
	trainingDatasetIDStr := r.FormValue("training-dataset-id")

	// Validate required fields
	if baseModel == "" || examplesCountStr == "" || trainingDatasetIDStr == "" {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">All fields are required</div>`))
		return
	}

	examplesCount, err := strconv.Atoi(examplesCountStr)
	if err != nil || examplesCount < 1 {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Please enter a valid number of examples</div>`))
		return
	}

	trainingDatasetID, err := uuid.Parse(trainingDatasetIDStr)
	if err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Invalid training dataset ID</div>`))
		return
	}

	// Create finetune request
	createReq := struct {
		BaseModelName                    string    `json:"base_model_name"`
		TrainingDatasetID                uuid.UUID `json:"training_dataset_id"`
		TrainingDatasetNumberExamples    int       `json:"training_dataset_number_examples"`
		TrainingDatasetSelectRandom      bool      `json:"training_dataset_select_random"`
	}{
		BaseModelName:                 baseModel,
		TrainingDatasetID:             trainingDatasetID,
		TrainingDatasetNumberExamples: examplesCount,
		TrainingDatasetSelectRandom:   randomSelection,
	}

	jsonData, err := json.Marshal(createReq)
	if err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to process request</div>`))
		return
	}

	apiBaseURL := web.GetAPIBaseURL(r)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/projects/%s/finetunes", apiBaseURL, projectID), bytes.NewBuffer(jsonData))
	if err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create request</div>`))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to connect to API</div>`))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// Extract error message from API response
		var apiError struct {
			Error string `json:"error"`
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create finetune</div>`))
			return
		}

		if err := json.Unmarshal(body, &apiError); err != nil || apiError.Error == "" {
			w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create finetune</div>`))
			return
		}

		w.Write([]byte(fmt.Sprintf(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">%s</div>`, apiError.Error)))
		return
	}

	// Success response
	w.Write([]byte(`<div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded">Fine-tuning started successfully! <a href="/web/home" class="underline">Check status on home page</a></div>`))
}