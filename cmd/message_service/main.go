package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ed16/messenger/internal/handlers"
	"github.com/ed16/messenger/internal/middleware"
	"github.com/ed16/messenger/internal/repository"
	"github.com/ed16/messenger/internal/tracing"
	"github.com/ed16/messenger/services/message"
)

const serverAddress = ":8080"

func main() {

	cleanup := tracing.InitTracer()
	defer cleanup()

	dbConn, err := repository.GetPostgresConn()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatalf("failed to close the database connection: %v", err)
		}
	}()

	messageRepo := repository.NewMessageRepo(dbConn)
	messageService := message.NewMessageService(messageRepo)

	router := setupRouter(messageService)
	tracedRouter := middleware.TraceRequests(router)     // Add tracing middleware
	loggedRouter := middleware.LogRequests(tracedRouter) // Add logging middleware

	server := &http.Server{
		Addr:    serverAddress,
		Handler: loggedRouter,
	}

	// Handle OS signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func setupRouter(messageService *message.MessageService) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /messages/user/{user_id}", handlers.CreateMessageHandler(messageService)) // POST: Send a personal message to a user
	router.HandleFunc("GET /messages", handlers.GetMessagesHandler(messageService))                   // GET: Retrieve personal messages

	return router
}
