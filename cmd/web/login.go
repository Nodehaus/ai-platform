package web

import (
	"bytes"
	"encoding/json"
	"net/http"

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

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if token := GetTokenFromCookie(r); token != "" {
		if isValidToken(r, token) {
			http.Redirect(w, r, "/web/home", http.StatusSeeOther)
			return
		}
		ClearTokenCookie(w)
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

	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		templ.Handler(LoginForm("Failed to process login request")).ServeHTTP(w, r)
		return
	}

	apiBaseURL := GetAPIBaseURL(r)
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

	SetTokenCookie(w, loginResp.Token)
	w.Header().Set("HX-Redirect", "/web/home")
}