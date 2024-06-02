package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockAuthService struct{}

func (m *MockAuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	if username == "validuser" && password == "validpassword" {
		return "validtoken", nil
	}
	return "", errors.New("authentication failed")
}

func (m *MockAuthService) ValidateToken(tokenString string) (userId int64, err error) {
	if tokenString == "validtoken" {
		return 12345, nil
	}
	return 0, errors.New("invalid token")
}

func (m *MockAuthService) GetToken(ctx context.Context, userId int64) (string, error) {
	return "", nil
}

func TestLoginHandler_Success(t *testing.T) {
	service := &MockAuthService{}
	handler := LoginHandler(service)

	reqBody, _ := json.Marshal(LoginRequest{
		Username: "validuser",
		Password: "validpassword",
	})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp LoginResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "validtoken", resp.Token)
}

func TestLoginHandler_Failure(t *testing.T) {
	service := &MockAuthService{}
	handler := LoginHandler(service)

	reqBody, _ := json.Marshal(LoginRequest{
		Username: "invaliduser",
		Password: "invalidpassword",
	})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Authentication failed\n", rr.Body.String())
}

func TestValidateTokenHandler_Success(t *testing.T) {
	service := &MockAuthService{}
	handler := ValidateTokenHandler(service)

	req, _ := http.NewRequest("GET", "/validate", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Token is valid", rr.Body.String())
	assert.Equal(t, "12345", rr.Header().Get("User-Id"))
}

func TestValidateTokenHandler_MissingAuthHeader(t *testing.T) {
	service := &MockAuthService{}
	handler := ValidateTokenHandler(service)

	req, _ := http.NewRequest("GET", "/validate", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Authorization header missing\n", rr.Body.String())
}

func TestValidateTokenHandler_InvalidAuthHeaderFormat(t *testing.T) {
	service := &MockAuthService{}
	handler := ValidateTokenHandler(service)

	req, _ := http.NewRequest("GET", "/validate", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid Authorization header format\n", rr.Body.String())
}

func TestValidateTokenHandler_InvalidToken(t *testing.T) {
	service := &MockAuthService{}
	handler := ValidateTokenHandler(service)

	req, _ := http.NewRequest("GET", "/validate", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Authentication failed\n", rr.Body.String())
}
