package main

import (
	"net/http"

	"github.com/ed16/messenger/pkg/handlers"
	"github.com/ed16/messenger/pkg/repository"
	"github.com/ed16/messenger/services/auth"
)

func main() {
	userRepo := &repository.UserRepoImpl{}
	authService := auth.AuthService{UserRepository: userRepo}

	http.HandleFunc("/auth/login", handlers.LoginHandler(&authService))
	http.HandleFunc("/auth/validate-token", handlers.ValidateTokenHandler(&authService))
	http.ListenAndServe(":8080", nil)
}
