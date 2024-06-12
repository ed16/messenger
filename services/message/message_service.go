package message

import (
	"context"

	"github.com/ed16/messenger/domain"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *domain.Message) (messageId int64, err error)
	GetMessagesByUserId(ctx context.Context, userId int64) ([]domain.Message, error)
}

type MessageService struct {
	MessageRepo MessageRepository
}

func NewMessageService(ur MessageRepository) *MessageService {
	return &MessageService{
		MessageRepo: ur,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, senderId, recipientId int64, content string) (messageId int64, err error) {
	message := &domain.Message{}

	message.SenderId = senderId
	message.RecipientId = recipientId
	message.Content = content
	message.IsRead = false
	message.IsReceived = false
	messageId, err = s.MessageRepo.CreateMessage(ctx, message)
	return
}

func (s *MessageService) GetMessagesByUserId(ctx context.Context, userId int64) ([]domain.Message, error) {
	return s.MessageRepo.GetMessagesByUserId(ctx, userId)
}
