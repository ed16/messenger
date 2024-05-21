package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/services/user"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddContactRequest struct {
	ContactUsername string `json:"contact_username"`
}

type UpdateProfileRequest struct {
	Description string `json:"description"`
	PhotoURL    string `json:"photo_url"`
}

type UserService interface {
	CreateNewUser(username, password string) error
	AddContact(userID int64, contactUsername string) error
	GetUserContacts(userID int64) ([]domain.User, error)
	UpdateUserProfile(profile *domain.Profile) error
	GetUsersByUsername(username string) ([]domain.User, error)
}

// ArticleHandler  represent the httphandler for article
type UserHandler struct {
	Service UserService
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

func ContactsHandler(service *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req AddContactRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()
			userIDStr := r.Header.Get("User-Id")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
			err = service.AddContact(userID, req.ContactUsername)
			if err != nil {
				http.Error(w, "Failed to add contact", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == http.MethodGet {

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
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
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
		userIDStr := r.Header.Get("User-Id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		updatedProfile := &domain.Profile{
			UserId:      userID,
			Description: req.Description,
			PhotoURL:    req.PhotoURL,
		}

		err = service.UpdateUserProfile(updatedProfile)
		if err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// func getStatusCode(err error) int {
// 	if err == nil {
// 		return http.StatusOK
// 	}

// 	switch err {
// 	case domain.ErrInternalServerError:
// 		return http.StatusInternalServerError
// 	case domain.ErrNotFound:
// 		return http.StatusNotFound
// 	case domain.ErrConflict:
// 		return http.StatusConflict
// 	default:
// 		return http.StatusInternalServerError
// 	}
// }
