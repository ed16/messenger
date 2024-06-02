package message

import (
	"context"

	"github.com/ed16/messenger/domain"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *domain.Message) (message_id int64, err error)
	GetMessagesByUserId(ctx context.Context, user_id int64) ([]domain.Message, error)
}

type MessageService struct {
	MessageRepo MessageRepository
}

func NewMessageService(ur MessageRepository) *MessageService {
	return &MessageService{
		MessageRepo: ur,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, sender_id, recipient_id int64, content string) (message_id int64, err error) {
	message := &domain.Message{}

	message.SenderId = sender_id
	message.RecipientId = recipient_id
	message.Content = content
	message.IsRead = false
	message.IsReceived = false
	message_id, err = s.MessageRepo.CreateMessage(ctx, message)
	return
}

func (s *MessageService) GetMessagesByUserId(ctx context.Context, user_id int64) ([]domain.Message, error) {
	return s.MessageRepo.GetMessagesByUserId(ctx, user_id)
}
