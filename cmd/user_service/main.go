package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ed16/messenger/internal/handlers"
	"github.com/ed16/messenger/internal/repository"
	"github.com/ed16/messenger/services/user"
)

func main() {

	dbConn := repository.GetMockDB()

	userRepo := repository.NewUserRepo(dbConn.DB())
	userService := user.NewUserService(userRepo)

	router := http.NewServeMux()

	router.HandleFunc("POST /users", handlers.CreateUserHandler(userService))           // POST: Register a new user
	router.HandleFunc("GET /users", handlers.GetUsersHandler(userService))              // GET: Retrieve users based on filter criteria
	router.HandleFunc("POST /users/contacts", handlers.ContactsHandler(userService))    // POST: Add a new contact by user ID; GET: Retrieve a user's contacts
	router.HandleFunc("GET /users/contacts", handlers.ContactsHandler(userService))     // POST: Add a new contact by user ID; GET: Retrieve a user's contacts
	router.HandleFunc("PUT /users/profile", handlers.UpdateProfileHandler(userService)) // Edit user profile details
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
