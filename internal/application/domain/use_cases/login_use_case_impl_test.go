package use_cases

import (
	"testing"
	"time"

	"ai-platform/internal/application/domain/entities"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	users map[string]*entities.User
}

func (m *mockUserRepository) FindByEmail(email string) (*entities.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func TestLoginUseCaseImpl_Login_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo := &mockUserRepository{
		users: map[string]*entities.User{
			"test@example.com": {
				ID:        uuid.New(),
				Email:     "test@example.com",
				Password:  string(hashedPassword),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	userService := &services.UserService{}
	jwtService := &services.JWTService{SecretKey: []byte("test-secret-key")}
	useCase := &LoginUseCaseImpl{
		UserRepository: mockRepo,
		UserService:    userService,
		JwtService:     jwtService,
	}

	command := in.LoginCommand{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResult, err := useCase.Login(command)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if loginResult == nil {
		t.Error("Expected login result, got nil")
	}

	if loginResult.User == nil {
		t.Error("Expected user, got nil")
	}

	if loginResult.User.Email != "test@example.com" {
		t.Errorf("Expected email %s, got %s", "test@example.com", loginResult.User.Email)
	}

	if loginResult.Token == "" {
		t.Error("Expected token, got empty string")
	}
}

func TestLoginUseCaseImpl_Login_InvalidCredentials(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo := &mockUserRepository{
		users: map[string]*entities.User{
			"test@example.com": {
				ID:        uuid.New(),
				Email:     "test@example.com",
				Password:  string(hashedPassword),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	userService := &services.UserService{}
	jwtService := &services.JWTService{SecretKey: []byte("test-secret-key")}
	useCase := &LoginUseCaseImpl{
		UserRepository: mockRepo,
		UserService:    userService,
		JwtService:     jwtService,
	}

	command := in.LoginCommand{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	loginResult, err := useCase.Login(command)

	if err == nil {
		t.Error("Expected error for invalid credentials")
	}

	if loginResult != nil {
		t.Error("Expected nil login result for invalid credentials")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got %s", err.Error())
	}
}

func TestLoginUseCaseImpl_Login_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		users: make(map[string]*entities.User),
	}

	userService := &services.UserService{}
	jwtService := &services.JWTService{SecretKey: []byte("test-secret-key")}
	useCase := &LoginUseCaseImpl{
		UserRepository: mockRepo,
		UserService:    userService,
		JwtService:     jwtService,
	}

	command := in.LoginCommand{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	user, err := useCase.Login(command)

	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if user != nil {
		t.Error("Expected nil user for non-existent user")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got %s", err.Error())
	}
}

func TestLoginUseCaseImpl_Login_EmptyEmail(t *testing.T) {
	mockRepo := &mockUserRepository{
		users: make(map[string]*entities.User),
	}

	userService := &services.UserService{}
	jwtService := &services.JWTService{SecretKey: []byte("test-secret-key")}
	useCase := &LoginUseCaseImpl{
		UserRepository: mockRepo,
		UserService:    userService,
		JwtService:     jwtService,
	}

	command := in.LoginCommand{
		Email:    "",
		Password: "password123",
	}

	user, err := useCase.Login(command)

	if err == nil {
		t.Error("Expected error for empty email")
	}

	if user != nil {
		t.Error("Expected nil user for empty email")
	}

	if err.Error() != "email is required" {
		t.Errorf("Expected 'email is required' error, got %s", err.Error())
	}
}

func TestLoginUseCaseImpl_Login_EmptyPassword(t *testing.T) {
	mockRepo := &mockUserRepository{
		users: make(map[string]*entities.User),
	}

	userService := &services.UserService{}
	jwtService := &services.JWTService{SecretKey: []byte("test-secret-key")}
	useCase := &LoginUseCaseImpl{
		UserRepository: mockRepo,
		UserService:    userService,
		JwtService:     jwtService,
	}

	command := in.LoginCommand{
		Email:    "test@example.com",
		Password: "",
	}

	user, err := useCase.Login(command)

	if err == nil {
		t.Error("Expected error for empty password")
	}

	if user != nil {
		t.Error("Expected nil user for empty password")
	}

	if err.Error() != "password is required" {
		t.Errorf("Expected 'password is required' error, got %s", err.Error())
	}
}