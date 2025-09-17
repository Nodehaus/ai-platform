package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
	} `json:"user"`
	Token   string `json:"token"`
	Message string `json:"message"`
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

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user already has a valid token
	if token := getTokenFromCookie(r); token != "" {
		// Verify token by calling profile endpoint
		if isValidToken(r, token) {
			http.Redirect(w, r, "/web/home", http.StatusSeeOther)
			return
		}
		// Clear invalid token
		clearTokenCookie(w)
	}

	templ.Handler(LoginPage()).ServeHTTP(w, r)
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		templ.Handler(LoginForm("Email and password are required")).ServeHTTP(w, r)
		return
	}

	// Call the API login endpoint
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		templ.Handler(LoginForm("Failed to process login request")).ServeHTTP(w, r)
		return
	}

	// Get the API base URL
	apiBaseURL := getAPIBaseURL(r)
	resp, err := http.Post(apiBaseURL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		templ.Handler(LoginForm("Failed to connect to authentication service")).ServeHTTP(w, r)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorMsg string
		if resp.StatusCode == http.StatusUnauthorized {
			errorMsg = "Invalid email or password"
		} else {
			errorMsg = "Login failed. Please try again."
		}
		templ.Handler(LoginForm(errorMsg)).ServeHTTP(w, r)
		return
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		templ.Handler(LoginForm("Failed to process login response")).ServeHTTP(w, r)
		return
	}

	// Set the JWT token as an HTTP-only cookie
	setTokenCookie(w, loginResp.Token)

	// Redirect to homepage
	w.Header().Set("HX-Redirect", "/web/home")
}

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	// Call the projects API to get user projects
	projectsData, userEmail, err := fetchProjectsData(r, token)
	if err != nil {
		// Token is invalid, clear it and redirect to login
		clearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	templ.Handler(Homepage(*projectsData, userEmail)).ServeHTTP(w, r)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearTokenCookie(w)
	http.Redirect(w, r, "/web/login", http.StatusSeeOther)
}

// Helper functions

func getTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func setTokenCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60, // 24 hours
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func clearTokenCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func getAPIBaseURL(r *http.Request) string {
	// In development, use the same host as the current request
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Check for environment variable first
	if baseURL := os.Getenv("API_BASE_URL"); baseURL != "" {
		return baseURL
	}

	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func isValidToken(r *http.Request, token string) bool {
	// Make a request to the projects endpoint to validate the token
	apiBaseURL := getAPIBaseURL(r)

	req, err := http.NewRequest("GET", apiBaseURL+"/api/projects", nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func fetchProjectsData(r *http.Request, token string) (*ProjectsData, string, error) {
	apiBaseURL := getAPIBaseURL(r)

	req, err := http.NewRequest("GET", apiBaseURL+"/api/projects", nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var projectsResp ProjectsResponse
	if err := json.Unmarshal(body, &projectsResp); err != nil {
		return nil, "", err
	}

	// Convert API response to template data structure
	projectsData := &ProjectsData{
		Projects: make([]ProjectData, len(projectsResp.Projects)),
	}

	for i, project := range projectsResp.Projects {
		projectsData.Projects[i] = ProjectData{
			ID:        project.ID,
			Name:      project.Name,
			Status:    project.Status,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
		}
	}

	// Extract user email from JWT token claims (simplified - just get from login response)
	// For now, we'll extract it from the token or return a placeholder
	userEmail := extractEmailFromToken(token)

	return projectsData, userEmail, nil
}

func extractEmailFromToken(token string) string {
	// Simplified - in a real implementation, you'd decode the JWT
	// For now, return a placeholder since we don't have the email easily accessible
	return "user@example.com"
}

// ProjectData represents a project for template display
type ProjectData struct {
	ID        uuid.UUID
	Name      string
	Status    string
	CreatedAt string
	UpdatedAt string
}

// ProjectsData represents the collection of projects for template display
type ProjectsData struct {
	Projects []ProjectData
}