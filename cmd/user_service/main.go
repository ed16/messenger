package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ed16/messenger/internal/handlers"
	"github.com/ed16/messenger/internal/repository"
	"github.com/ed16/messenger/services/user"
)

func main() {

	dbConn, err := repository.GetPostgresConn()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()
	userRepo := repository.NewUserRepo(dbConn)
	userService := user.NewUserService(userRepo)

	router := http.NewServeMux()

	router.HandleFunc("POST /users", handlers.CreateUserHandler(userService))           // POST: Register a new user
	router.HandleFunc("GET /users", handlers.GetUsersHandler(userService))              // GET: Retrieve users based on filter criteria
	router.HandleFunc("POST /users/contacts", handlers.ContactsHandler(userService))    // POST: Add a new contact by user ID; GET: Retrieve a user's contacts
	router.HandleFunc("GET /users/contacts", handlers.ContactsHandler(userService))     // POST: Add a new contact by user ID; GET: Retrieve a user's contacts
	router.HandleFunc("PUT /users/profile", handlers.UpdateProfileHandler(userService)) // Edit user profile details

	// Wrap the router with a middleware that logs each request
	server := http.Server{
		Addr:    ":8080",
		Handler: logRequests(router),
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

// Middleware to log each request and its response status code
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the request and record the response status code
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		// Log the request and response status code
		log.Printf("%s %s - Response status code: %d", r.Method, r.URL.Path, rw.status)
	})
}

// Custom ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// Override WriteHeader to capture the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
