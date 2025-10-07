package deployments

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

type DeploymentIndexData struct {
	ProjectID    string
	ProjectName  string
	DeploymentID string
	Deployment   DeploymentData
}

type DeploymentData struct {
	ID         uuid.UUID  `json:"id"`
	ModelName  string     `json:"model_name"`
	APIKey     string     `json:"api_key"`
	ProjectID  uuid.UUID  `json:"project_id"`
	FinetuneID *uuid.UUID `json:"finetune_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func DeploymentIndexHandler(w http.ResponseWriter, r *http.Request) {
	token := web.GetTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	// Extract project ID and deployment ID from URL path
	// Expected format: /web/projects/{project_id}/deployments/{deployment_id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 || pathParts[3] == "" || pathParts[5] == "" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	projectIDStr := pathParts[3]
	deploymentIDStr := pathParts[5]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID format", http.StatusBadRequest)
		return
	}

	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		http.Error(w, "Invalid deployment ID format", http.StatusBadRequest)
		return
	}

	// Fetch deployment data
	deploymentData, err := fetchDeploymentData(r, token, projectID, deploymentID)
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

	indexData := DeploymentIndexData{
		ProjectID:    projectIDStr,
		ProjectName:  projectName,
		DeploymentID: deploymentIDStr,
		Deployment:   *deploymentData,
	}

	templ.Handler(DeploymentIndex(indexData)).ServeHTTP(w, r)
}

func fetchDeploymentData(r *http.Request, token string, projectID uuid.UUID, deploymentID uuid.UUID) (*DeploymentData, error) {
	apiBaseURL := web.GetAPIBaseURL(r)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/%s/deployments/%s", apiBaseURL, projectID, deploymentID), nil)
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

	var deployment DeploymentData
	if err := json.NewDecoder(resp.Body).Decode(&deployment); err != nil {
		return nil, err
	}

	return &deployment, nil
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
