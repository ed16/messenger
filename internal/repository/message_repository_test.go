package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	messageRepo := repository.NewMessageRepo(db)

	message := &domain.Message{
		SenderId:    1,
		RecipientId: 2,
		Content:     "Hello, World!",
		IsRead:      false,
		IsReceived:  false,
	}

	mock.ExpectQuery(`INSERT INTO messages`).
		WithArgs(message.SenderId, message.RecipientId, message.Content, message.IsRead, message.IsReceived).
		WillReturnRows(sqlmock.NewRows([]string{"message_id"}).AddRow(1))

	messageID, err := messageRepo.CreateMessage(context.Background(), message)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), messageID)
}

func TestGetMessagesByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	messageRepo := repository.NewMessageRepo(db)

	mediaID := int64(1)
	expectedMessage := domain.Message{
		MessageId:   1,
		SenderId:    1,
		RecipientId: 2,
		CreatedAt:   time.Now(),
		Content:     "Hello, World!",
		IsRead:      false,
		IsReceived:  false,
		MediaId:     &mediaID,
	}

	mock.ExpectQuery(`SELECT message_id, sender_id, recipient_id, created_at, content, is_read, is_received, media_id FROM messages m WHERE m.sender_id = \$1 OR m.recipient_id = \$1 ORDER BY m.message_id ASC LIMIT 100`).
		WithArgs(expectedMessage.SenderId).
		WillReturnRows(sqlmock.NewRows([]string{"message_id", "sender_id", "recipient_id", "created_at", "content", "is_read", "is_received", "media_id"}).
			AddRow(expectedMessage.MessageId, expectedMessage.SenderId, expectedMessage.RecipientId, expectedMessage.CreatedAt, expectedMessage.Content, expectedMessage.IsRead, expectedMessage.IsReceived, expectedMessage.MediaId))

	messages, err := messageRepo.GetMessagesByUserId(context.Background(), expectedMessage.SenderId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, expectedMessage, messages[0])
}

func TestGetMessagesByUserId_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	messageRepo := repository.NewMessageRepo(db)

	mock.ExpectQuery(`SELECT message_id, sender_id, recipient_id, created_at, content, is_read, is_received, media_id FROM messages m WHERE m.sender_id = \$1 OR m.recipient_id = \$1 ORDER BY m.message_id ASC LIMIT 100`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"message_id", "sender_id", "recipient_id", "created_at", "content", "is_read", "is_received", "media_id"}))

	messages, err := messageRepo.GetMessagesByUserId(context.Background(), int64(1))
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Empty(t, messages)
}

func TestCreateMessage_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	messageRepo := repository.NewMessageRepo(db)

	message := &domain.Message{
		SenderId:    1,
		RecipientId: 2,
		Content:     "Hello, World!",
		IsRead:      false,
		IsReceived:  false,
	}

	mock.ExpectQuery(`INSERT INTO messages`).
		WithArgs(message.SenderId, message.RecipientId, message.Content, message.IsRead, message.IsReceived).
		WillReturnError(sql.ErrConnDone)

	messageID, err := messageRepo.CreateMessage(context.Background(), message)
	assert.Error(t, err)
	assert.Equal(t, int64(0), messageID)
}
