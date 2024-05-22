package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ed16/messenger/internal/handlers"
	"github.com/ed16/messenger/internal/repository"
	"github.com/ed16/messenger/services/auth"
)

func main() {
	dbConn := repository.GetMockDB()

	userRepo := repository.NewUserRepo(dbConn.DB())
	authService := auth.NewAuthService(userRepo)

	router := http.NewServeMux()
	router.HandleFunc("/auth/login/", handlers.LoginHandler(authService))
	router.HandleFunc("/auth/validate-token", handlers.ValidateTokenHandler(authService))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Create a channel to listen for interrupts (SIGINT, SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Start the server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for an interrupt
	<-stop

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	// Server gracefully shutdown
}
