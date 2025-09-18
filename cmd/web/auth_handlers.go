package web

import (
	"fmt"
	"net/http"
	"os"
	"time"
)





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

