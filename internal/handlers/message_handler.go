package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/services/message"
)

type CreateMessageRequest struct {
	Content string `json:"content"`
}
type CreateMessageResponse struct {
	MessageId string `json:"message_id"`
}

type MessageService interface {
	CreateNewUser(ctx context.Context, password string) error
	AddContact(ctx context.Context, userID int64, contactUsername string) error
	GetUserContacts(ctx context.Context, userID int64) ([]domain.User, error)
	UpdateUserProfile(ctx context.Context, profile *domain.Profile) error
	GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error)
}

type MessageHandler struct {
	Service MessageService
}

func CreateMessageHandler(service *message.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipient_id_str := r.URL.Query().Get("user_id")
		// Validate request
		if recipient_id_str == "" {
			http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
			return
		}

		var req CreateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		sender_id_str := r.Header.Get("User-Id")
		sender_id, err := strconv.ParseInt(sender_id_str, 10, 64)
		recipient_id, err := strconv.ParseInt(recipient_id_str, 10, 64)

		message_id, err := service.CreateMessage(r.Context(), sender_id, recipient_id, req.Content)
		if err != nil {
			log.Println(err)
			http.Error(w, "Sending message failed", getStatusCode(err))
			return
		}
		message_id_str := strconv.Itoa(int(message_id))
		resp := CreateMessageResponse{MessageId: message_id_str}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func GetMessagesHandler(service *message.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id_str := r.Header.Get("User-Id")
		user_id, err := strconv.ParseInt(user_id_str, 10, 64)

		users, err := service.GetMessagesByUserId(r.Context(), user_id)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve messages", getStatusCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
