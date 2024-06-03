package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Custom claims structure, add any fields that you might need in your token
type CustomClaims struct {
	jwt.RegisteredClaims
	// We can add custom fields for additional claims, e.g., UserRole
	UserRole string
}

type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
	ValidateToken(tokenString string) (userId int64, err error)
	GetToken(ctx context.Context, userId int64) (string, error)
}

func LoginHandler(service AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		// Validate request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		token, err := service.Authenticate(r.Context(), req.Username, req.Password)
		if err != nil {
			log.Println(err)
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		resp := LoginResponse{Token: token}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println(err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	}
}

func ValidateTokenHandler(service AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Split the header value to get the token part
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		userId, err := service.ValidateToken(token)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		_, err = w.Write([]byte("Token is valid"))
		if err != nil {
			log.Println(err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("User-Id", fmt.Sprint(userId))
		w.WriteHeader(http.StatusOK)
	}
}
