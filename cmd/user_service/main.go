package main

import (
	"net/http"

	"github.com/ed16/messenger/internal/handlers"
	"github.com/ed16/messenger/internal/repository"
	"github.com/ed16/messenger/services/user"
)

func main() {
	dbConn, _ := repository.GetPostgresConn()
	userRepo := repository.NewUserRepo(dbConn)
	userService := user.NewUserService(userRepo)

	http.HandleFunc("/users", handlers.UsersHandler(userService))                 // POST: Register a new user; GET: Retrieve users based on filter criteria
	http.HandleFunc("/users/contacts", handlers.ContactsHandler(userService))     // POST: Add a new contact by user ID; GET: Retrieve a user's contacts
	http.HandleFunc("/users/profile", handlers.UpdateProfileHandler(userService)) // Edit user profile details

	http.ListenAndServe(":8080", nil)
}
