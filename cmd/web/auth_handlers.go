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

type ProfileResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Email   string    `json:"email"`
	Message string    `json:"message"`
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user already has a valid token
	if token := getTokenFromCookie(r); token != "" {
		// Verify token by calling profile endpoint
		if isValidToken(token) {
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

	// Call the profile API to get user data
	profileData, err := fetchProfileData(r, token)
	if err != nil {
		// Token is invalid, clear it and redirect to login
		clearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	templ.Handler(Homepage(*profileData)).ServeHTTP(w, r)
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

func isValidToken(token string) bool {
	// Make a request to the profile endpoint to validate the token
	req, err := http.NewRequest("GET", "http://localhost:8080/api/profile", nil)
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

func fetchProfileData(r *http.Request, token string) (*ProfileData, error) {
	apiBaseURL := getAPIBaseURL(r)

	req, err := http.NewRequest("GET", apiBaseURL+"/api/profile", nil)
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

	var profileResp ProfileResponse
	if err := json.Unmarshal(body, &profileResp); err != nil {
		return nil, err
	}

	return &ProfileData{
		UserID:  profileResp.UserID,
		Email:   profileResp.Email,
		Message: profileResp.Message,
	}, nil
}