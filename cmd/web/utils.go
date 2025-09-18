package web

import (
	"fmt"
	"net/http"
	"os"
)

func GetTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func SetTokenCookie(w http.ResponseWriter, token string) {
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

func ClearTokenCookie(w http.ResponseWriter) {
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

func GetAPIBaseURL(r *http.Request) string {
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