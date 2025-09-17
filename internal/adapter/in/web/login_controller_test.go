package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/port/in"
)

type mockLoginUseCase struct {
	user *entities.User
	err  error
}

func (m *mockLoginUseCase) Login(command in.LoginCommand) (*entities.User, error) {
	return m.user, m.err
}

func TestLoginController_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &entities.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	mockUseCase := &mockLoginUseCase{
		user: user,
		err:  nil,
	}

	controller := NewLoginController(mockUseCase)

	router := gin.New()
	router.POST("/login", controller.Login)

	loginRequest := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	requestBody, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response LoginResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response.User.Email != "test@example.com" {
		t.Errorf("Expected email %s, got %s", "test@example.com", response.User.Email)
	}

	if response.Message != "Login successful" {
		t.Errorf("Expected message 'Login successful', got %s", response.Message)
	}
}

func TestLoginController_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockLoginUseCase{
		user: nil,
		err:  errors.New("invalid credentials"),
	}

	controller := NewLoginController(mockUseCase)

	router := gin.New()
	router.POST("/login", controller.Login)

	loginRequest := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	requestBody, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "invalid credentials" {
		t.Errorf("Expected error 'invalid credentials', got %s", response["error"])
	}
}

func TestLoginController_Login_InvalidRequestFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockLoginUseCase{}
	controller := NewLoginController(mockUseCase)

	router := gin.New()
	router.POST("/login", controller.Login)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if response["error"] != "Invalid request format" {
		t.Errorf("Expected error 'Invalid request format', got %s", response["error"])
	}
}