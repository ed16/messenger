package message_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/repository"
	"github.com/ed16/messenger/services/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMessage(t *testing.T) {
	mockRepo := new(repository.MockMessageRepository)
	service := message.NewMessageService(mockRepo)
	ctx := context.Background()

	messageID := int64(1)
	mockMessage := &domain.Message{
		SenderId:    1,
		RecipientId: 2,
		Content:     "Hello, World!",
		Status:      domain.MessageStatusSent,
	}

	mockRepo.On("CreateMessage", ctx, mock.AnythingOfType("*domain.Message")).Return(messageID, nil)

	id, err := service.CreateMessage(ctx, mockMessage.SenderId, mockMessage.RecipientId, mockMessage.Content)
	assert.NoError(t, err)
	assert.Equal(t, messageID, id)
	mockRepo.AssertExpectations(t)
}

func TestCreateMessage_Error(t *testing.T) {
	mockRepo := new(repository.MockMessageRepository)
	service := message.NewMessageService(mockRepo)
	ctx := context.Background()

	mockMessage := &domain.Message{
		SenderId:    1,
		RecipientId: 2,
		Content:     "Hello, World!",
		Status:      domain.MessageStatusSent,
	}

	mockRepo.On("CreateMessage", ctx, mock.AnythingOfType("*domain.Message")).Return(int64(0), errors.New("create message error"))

	id, err := service.CreateMessage(ctx, mockMessage.SenderId, mockMessage.RecipientId, mockMessage.Content)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	mockRepo.AssertExpectations(t)
}

func TestGetMessagesByUserId(t *testing.T) {
	mockRepo := new(repository.MockMessageRepository)
	service := message.NewMessageService(mockRepo)
	ctx := context.Background()

	expectedMessages := []domain.Message{
		{
			MessageId:   1,
			SenderId:    1,
			RecipientId: 2,
			Content:     "Hello, World!",
			Status:      domain.MessageStatusSent,
		},
	}

	mockRepo.On("GetMessagesByUserId", ctx, int64(1)).Return(expectedMessages, nil)

	messages, err := service.GetMessagesByUserId(ctx, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetMessagesByUserId_NotFound(t *testing.T) {
	mockRepo := new(repository.MockMessageRepository)
	service := message.NewMessageService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetMessagesByUserId", ctx, int64(1)).Return([]domain.Message{}, domain.ErrNotFound)

	messages, err := service.GetMessagesByUserId(ctx, int64(1))
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Empty(t, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetMessagesByUserId_Error(t *testing.T) {
	mockRepo := new(repository.MockMessageRepository)
	service := message.NewMessageService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetMessagesByUserId", ctx, int64(1)).Return([]domain.Message(nil), errors.New("get messages error"))

	messages, err := service.GetMessagesByUserId(ctx, int64(1))
	assert.Error(t, err)
	assert.Nil(t, messages)
	mockRepo.AssertExpectations(t)
}
