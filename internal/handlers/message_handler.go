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
	CreateMessage(ctx context.Context, senderId, recipientId int64, content string) (messageId int64, err error)
	GetMessagesByUserId(ctx context.Context, userId int64) ([]domain.Message, error)
}

func CreateMessageHandler(service MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipientIdStr := r.PathValue("user_id")
		// Validate request
		if recipientIdStr == "" {
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
		senderIdStr := r.Header.Get("User-Id")
		senderId, err := strconv.ParseInt(senderIdStr, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, "Sending message failed", http.StatusInternalServerError)
			return
		}
		recipientId, err := strconv.ParseInt(recipientIdStr, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, "Sending message failed", http.StatusInternalServerError)
			return
		}

		messageId, err := service.CreateMessage(r.Context(), senderId, recipientId, req.Content)
		if err != nil {
			log.Println(err)
			http.Error(w, "Sending message failed", getStatusCode(err))
			return
		}
		messageIdStr := strconv.Itoa(int(messageId))
		resp := CreateMessageResponse{MessageId: messageIdStr}

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
		userIdStr := r.Header.Get("User-Id")
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, "Getting messages failed", http.StatusInternalServerError)
			return
		}

		users, err := service.GetMessagesByUserId(r.Context(), userId)
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
