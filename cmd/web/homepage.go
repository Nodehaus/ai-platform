package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type ProjectData struct {
	ID        uuid.UUID
	Name      string
	Status    string
	CreatedAt string
	UpdatedAt string
}

type ProjectsData struct {
	Projects []ProjectData
}

type ProjectResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type ProjectsResponse struct {
	Projects []ProjectResponse `json:"projects"`
}

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	projectsData, err := fetchProjectsData(r, token)
	if err != nil {
		clearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	templ.Handler(Homepage(*projectsData)).ServeHTTP(w, r)
}

func fetchProjectsData(r *http.Request, token string) (*ProjectsData, error) {
	apiBaseURL := getAPIBaseURL(r)

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

		projectsData.Projects[i] = ProjectData{
			ID:        project.ID,
			Name:      project.Name,
			Status:    project.Status,
			CreatedAt: formattedCreatedAt,
			UpdatedAt: project.UpdatedAt,
		}
	}

	return projectsData, nil
}
