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

func TrainingDatasetStep4Handler(w http.ResponseWriter, r *http.Request) {
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

	// Get parameters from query
	corpus := r.URL.Query().Get("corpus")
	language := r.URL.Query().Get("language")
	prompt := r.URL.Query().Get("prompt")
	fieldNames := r.URL.Query().Get("field_names")
	inputField := r.URL.Query().Get("input_field")
	outputField := r.URL.Query().Get("output_field")

	if corpus == "" || language == "" || prompt == "" || fieldNames == "" || inputField == "" || outputField == "" {
		http.Redirect(w, r, "/web/projects/"+projectIDStr+"/training-datasets/step1", http.StatusSeeOther)
		return
	}

	// Fetch project details to get the project name
	projectName, err := fetchProjectName(r, token, projectID)
	if err != nil {
		web.ClearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	templ.Handler(TrainingDatasetStep4(projectIDStr, projectName, corpus, language, prompt, fieldNames, inputField, outputField)).ServeHTTP(w, r)
}

func CreateTrainingDatasetHandler(w http.ResponseWriter, r *http.Request) {
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
	corpus := r.FormValue("corpus")
	language := r.FormValue("language")
	prompt := r.FormValue("prompt")
	fieldNames := r.FormValue("field_names")
	inputField := r.FormValue("input_field")
	outputField := r.FormValue("output_field")
	examplesCountStr := r.FormValue("examples_count")

	// Validate required fields
	if corpus == "" || language == "" || prompt == "" || fieldNames == "" || inputField == "" || outputField == "" || examplesCountStr == "" {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">All fields are required</div>`))
		return
	}

	examplesCount, err := strconv.Atoi(examplesCountStr)
	if err != nil || examplesCount < 1 || examplesCount > 1000 {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Please enter a valid number of examples (1-1000)</div>`))
		return
	}

	// Split field names by comma and trim spaces
	fieldNamesSlice := make([]string, 0)
	for _, name := range strings.Split(fieldNames, ",") {
		trimmed := strings.TrimSpace(name)
		if trimmed != "" {
			fieldNamesSlice = append(fieldNamesSlice, trimmed)
		}
	}

	// Create training dataset request
	createReq := struct {
		CorpusName              string   `json:"corpus_name"`
		InputField              string   `json:"input_field"`
		OutputField             string   `json:"output_field"`
		LanguageISO             string   `json:"language_iso"`
		FieldNames              []string `json:"field_names"`
		GeneratePrompt          string   `json:"generate_prompt"`
		GenerateExamplesNumber  int      `json:"generate_examples_number"`
	}{
		CorpusName:             corpus,
		InputField:             inputField,
		OutputField:            outputField,
		LanguageISO:            language,
		FieldNames:             fieldNamesSlice,
		GeneratePrompt:         prompt,
		GenerateExamplesNumber: examplesCount,
	}

	jsonData, err := json.Marshal(createReq)
	if err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to process request</div>`))
		return
	}

	apiBaseURL := web.GetAPIBaseURL(r)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/projects/%s/training-datasets", apiBaseURL, projectID), bytes.NewBuffer(jsonData))
	if err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create request</div>`))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second} // Longer timeout for training dataset creation
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
			w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create training dataset</div>`))
			return
		}

		if err := json.Unmarshal(body, &apiError); err != nil || apiError.Error == "" {
			w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create training dataset</div>`))
			return
		}

		w.Write([]byte(fmt.Sprintf(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">%s</div>`, apiError.Error)))
		return
	}

	// Success response
	w.Write([]byte(`<div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded">Training dataset created successfully! <a href="/web/projects/` + projectIDStr + `" class="underline">Go to project</a></div>`))
}