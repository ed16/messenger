package main

import (
	"net/http"

	"github.com/ed16/messenger/internal/handlers"
	"github.com/ed16/messenger/internal/repository"
	"github.com/ed16/messenger/services/auth"
)

func main() {

	dbConn, _ := repository.GetPostgresConn()

	userRepo := repository.NewUserRepo(dbConn)
	authService := auth.AuthService{UserRepository: userRepo}

	http.HandleFunc("/auth/login", handlers.LoginHandler(&authService))
	http.HandleFunc("/auth/validate-token", handlers.ValidateTokenHandler(&authService))
	http.ListenAndServe(":8080", nil)
}
