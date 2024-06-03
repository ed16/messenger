package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ed16/messenger/domain"
)

type CreateMessageRequest struct {
	Content string `json:"content"`
}
type CreateMessageResponse struct {
	MessageId string `json:"message_id"`
}

type MessageService interface {
	CreateMessage(ctx context.Context, sender_id, recipient_id int64, content string) (message_id int64, err error)
	GetMessagesByUserId(ctx context.Context, user_id int64) ([]domain.Message, error)
}

func CreateMessageHandler(service MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipient_id_str := r.PathValue("user_id")
		// Validate request
		if recipient_id_str == "" {
			http.Error(w, "userId query parameter is required", http.StatusBadRequest)
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
		if err != nil {
			log.Println(err)
			http.Error(w, "Sending message failed", http.StatusInternalServerError)
			return
		}
		recipient_id, err := strconv.ParseInt(recipient_id_str, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, "Sending message failed", http.StatusInternalServerError)
			return
		}

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
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println(err)
			http.Error(w, "Sending message failed", http.StatusInternalServerError)
			return
		}

	}
}

func GetMessagesHandler(service MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id_str := r.Header.Get("User-Id")
		user_id, err := strconv.ParseInt(user_id_str, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, "Getting messages failed", http.StatusInternalServerError)
			return
		}

		users, err := service.GetMessagesByUserId(r.Context(), user_id)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve messages", getStatusCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			log.Println(err)
			http.Error(w, "Sending message failed", http.StatusInternalServerError)
			return
		}

	}
}
