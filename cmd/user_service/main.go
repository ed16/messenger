package main

import (
	"net/http"

	"github.com/ed16/messenger/pkg/handlers"
	"github.com/ed16/messenger/pkg/repository"
	"github.com/ed16/messenger/services/auth"
	"github.com/ed16/messenger/services/user"
)

func main() {
	userRepo := &repository.UserRepoImpl{}
	authService := auth.AuthService{SecretKey: ""}
	userService := user.UserService{
		UserRepository: userRepo,
		AuthService:    authService,
	}

	http.HandleFunc("/users", handlers.UsersHandler(&userService))                  // POST: Register a new user; GET: Retrieve users based on filter criteria
	http.HandleFunc("/users/contacts/id", handlers.AddContactHandler(&userService)) // Add a new contact by user ID
	http.HandleFunc("/users/contacts", handlers.GetContactsHandler(&userService))   // Retrieve a user's contacts
	http.HandleFunc("/users/profile", handlers.UpdateProfileHandler(&userService))  // Edit user profile details

	http.ListenAndServe(":8080", nil)
}
