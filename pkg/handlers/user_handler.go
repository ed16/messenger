package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ed16/messenger/pkg/domain"
	"github.com/ed16/messenger/services/user"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UsersHandler(service *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req RegisterRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			err := service.CreateNewUser(req.Username, req.Password)
			if err != nil {
				http.Error(w, "User registration failed", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		} else if r.Method == http.MethodGet {
			username := r.URL.Query().Get("username")
			if username == "" {
				http.Error(w, "Username query parameter is required", http.StatusBadRequest)
				return
			}

			users, err := service.GetUsersByUsername(username)
			if err != nil {
				http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

type AddContactRequest struct {
	UserID    int64 `json:"user_id"`
	ContactID int64 `json:"contact_id"`
}

type UpdateProfileRequest struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Status   byte   `json:"status"`
}

func AddContactHandler(service *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AddContactRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err := service.AddContact(req.UserID, req.ContactID)
		if err != nil {
			http.Error(w, "Failed to add contact", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetContactsHandler(service *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		contacts, err := service.GetUserContacts(userID)
		if err != nil {
			http.Error(w, "Failed to retrieve contacts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contacts)
	}
}

func UpdateProfileHandler(service *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req UpdateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		updatedUser := &domain.User{
			Username: req.Username,
			Status:   req.Status,
		}

		err := service.UpdateUserProfile(req.UserID, updatedUser)
		if err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
