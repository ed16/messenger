package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ed16/messenger/domain"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateContactRequest struct {
	ContactUsername string `json:"contactUsername"`
}

type UserService interface {
	CreateUser(ctx context.Context, username, password string) error
	AddContact(ctx context.Context, userID int64, contactUsername string) error
	GetUserContacts(ctx context.Context, userID int64) ([]domain.Contact, error)
	GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error)
}

func CreateUserHandler(service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		// Validate request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		err := service.CreateUser(r.Context(), req.Username, req.Password)
		if err != nil {
			log.Println(err)
			http.Error(w, "User registration failed", getStatusCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetUsersHandler(service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		// Validate request
		if username == "" {
			http.Error(w, "Username query parameter is required", http.StatusBadRequest)
			return
		}

		users, err := service.GetUsersByUsername(r.Context(), username)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve users", getStatusCode(err))
			return
		}
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	}
}

func CreateContactHandler(service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateContactRequest
		// Validate request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		userIDStr := r.Header.Get("User-Id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		err = service.AddContact(r.Context(), userID, req.ContactUsername)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to add contact", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetContactsHandler(service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("User-Id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		// Validate request
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		contacts, err := service.GetUserContacts(r.Context(), userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve contacts", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(contacts)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrUserAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
