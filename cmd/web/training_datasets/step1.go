package training_datasets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/google/uuid"

	"ai-platform/cmd/web"
)

func TrainingDatasetStep1Handler(w http.ResponseWriter, r *http.Request) {
	token := web.GetTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	// Extract project ID from URL path
	// Expected format: /web/projects/{project_id}/training-datasets/step1
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

	// Fetch project details to get the project name
	projectName, err := fetchProjectName(r, token, projectID)
	if err != nil {
		// If we can't fetch the project, redirect to login (token might be invalid)
		web.ClearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	templ.Handler(TrainingDatasetStep1(projectIDStr, projectName)).ServeHTTP(w, r)
}

func fetchProjectName(r *http.Request, token string, projectID uuid.UUID) (string, error) {
	apiBaseURL := web.GetAPIBaseURL(r)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/%s", apiBaseURL, projectID), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var project web.ProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return "", err
	}

	return project.Name, nil
}