package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type TrainingDatasetData struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type FinetuneData struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type DeploymentData struct {
	ID        uuid.UUID `json:"id"`
	ModelName string    `json:"model_name"`
}

type ProjectData struct {
	ID              uuid.UUID             `json:"id"`
	Name            string                `json:"name"`
	Status          string                `json:"status"`
	TrainingDataset *TrainingDatasetData  `json:"training_dataset"`
	Finetune        *FinetuneData         `json:"finetune"`
	Deployments     []DeploymentData      `json:"deployments"`
	CreatedAt       string                `json:"created_at"`
	UpdatedAt       string                `json:"updated_at"`
}

type ProjectsData struct {
	Projects []ProjectData
}

type TrainingDatasetResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type FinetuneResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type DeploymentResponse struct {
	ID        uuid.UUID `json:"id"`
	ModelName string    `json:"model_name"`
}

type ProjectResponse struct {
	ID              uuid.UUID                `json:"id"`
	Name            string                   `json:"name"`
	Status          string                   `json:"status"`
	TrainingDataset *TrainingDatasetResponse `json:"training_dataset"`
	Finetune        *FinetuneResponse        `json:"finetune"`
	Deployments     []DeploymentResponse     `json:"deployments"`
	CreatedAt       string                   `json:"created_at"`
	UpdatedAt       string                   `json:"updated_at"`
}

type ProjectsResponse struct {
	Projects []ProjectResponse `json:"projects"`
}

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type CreateProjectResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	token := GetTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	projectsData, err := fetchProjectsData(r, token)
	if err != nil {
		ClearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	templ.Handler(Home(*projectsData)).ServeHTTP(w, r)
}

func fetchProjectsData(r *http.Request, token string) (*ProjectsData, error) {
	apiBaseURL := GetAPIBaseURL(r)

	req, err := http.NewRequest("GET", apiBaseURL+"/api/projects", nil)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var projectsResp ProjectsResponse
	if err := json.Unmarshal(body, &projectsResp); err != nil {
		return nil, err
	}

	projectsData := &ProjectsData{
		Projects: make([]ProjectData, len(projectsResp.Projects)),
	}

	for i, project := range projectsResp.Projects {
		createdAt, err := time.Parse(time.RFC3339, project.CreatedAt)
		var formattedCreatedAt string
		if err != nil {
			formattedCreatedAt = project.CreatedAt
		} else {
			formattedCreatedAt = createdAt.Format("15:04 on 02.01.2006")
		}

		var trainingDataset *TrainingDatasetData
		if project.TrainingDataset != nil {
			trainingDataset = &TrainingDatasetData{
				ID:     project.TrainingDataset.ID,
				Status: project.TrainingDataset.Status,
			}
		}

		var finetune *FinetuneData
		if project.Finetune != nil {
			finetune = &FinetuneData{
				ID:     project.Finetune.ID,
				Status: project.Finetune.Status,
			}
		}

		deployments := make([]DeploymentData, len(project.Deployments))
		for j, deployment := range project.Deployments {
			deployments[j] = DeploymentData{
				ID:        deployment.ID,
				ModelName: deployment.ModelName,
			}
		}

		projectsData.Projects[i] = ProjectData{
			ID:              project.ID,
			Name:            project.Name,
			Status:          project.Status,
			TrainingDataset: trainingDataset,
			Finetune:        finetune,
			Deployments:     deployments,
			CreatedAt:       formattedCreatedAt,
			UpdatedAt:       project.UpdatedAt,
		}
	}

	return projectsData, nil
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := GetTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	projectName := r.FormValue("projectName")
	if projectName == "" {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Project name is required</div>`))
		return
	}

	createReq := CreateProjectRequest{
		Name: projectName,
	}

	jsonData, err := json.Marshal(createReq)
	if err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to process request</div>`))
		return
	}

	apiBaseURL := GetAPIBaseURL(r)
	req, err := http.NewRequest("POST", apiBaseURL+"/api/projects", bytes.NewBuffer(jsonData))
	if err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create request</div>`))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
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
			w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create project</div>`))
			return
		}

		if err := json.Unmarshal(body, &apiError); err != nil || apiError.Error == "" {
			w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to create project</div>`))
			return
		}

		w.Write([]byte(fmt.Sprintf(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">%s</div>`, apiError.Error)))
		return
	}

	var createResp CreateProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		w.Write([]byte(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">Failed to process response</div>`))
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/web/projects/%s/training-datasets/step0", createResp.ID))
}
