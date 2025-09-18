package web

import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearTokenCookie(w)
	http.Redirect(w, r, "/web/login", http.StatusSeeOther)
}

func isValidToken(r *http.Request, token string) bool {
	// Make a request to the projects endpoint to validate the token
	apiBaseURL := GetAPIBaseURL(r)

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

