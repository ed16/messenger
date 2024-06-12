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
		Status:      domain.MessageStatusSent,
		CreatedAt:   time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO messages`).
		WithArgs(message.SenderId, message.RecipientId, message.Content, message.Status).
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
		Content:     "Hello, World!",
		CreatedAt:   time.Now(),
		Status:      domain.MessageStatusSent,
		MediaId:     &mediaID,
	}

	mock.ExpectQuery(`SELECT 
			message_id, 
			sender_id,
			recipient_id, 
			content,
			created_at, 
			status,
			media_id 
		FROM messages m WHERE m.sender_id = \$1 OR m.recipient_id = \$1 ORDER BY m.message_id ASC LIMIT 100`).
		WithArgs(expectedMessage.SenderId).
		WillReturnRows(sqlmock.NewRows([]string{"message_id", "sender_id", "recipient_id", "content", "created_at", "status", "media_id"}).
			AddRow(expectedMessage.MessageId, expectedMessage.SenderId, expectedMessage.RecipientId, expectedMessage.Content, expectedMessage.CreatedAt, expectedMessage.Status, expectedMessage.MediaId))

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

	mock.ExpectQuery(`SELECT 
			message_id, 
			sender_id,
			recipient_id, 
			content,
			created_at, 
			status,
			media_id
		FROM messages m WHERE m.sender_id = \$1 OR m.recipient_id = \$1 ORDER BY m.message_id ASC LIMIT 100`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"message_id", "sender_id", "recipient_id", "content", "created_at", "status", "media_id"}))

	messages, err := messageRepo.GetMessagesByUserId(context.Background(), int64(1))
	assert.NoError(t, err)
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
		Status:      domain.MessageStatusSent,
		CreatedAt:   time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO messages`).
		WithArgs(message.SenderId, message.RecipientId, message.Content, message.Status).
		WillReturnError(sql.ErrConnDone)

	messageID, err := messageRepo.CreateMessage(context.Background(), message)
	assert.Error(t, err)
	assert.Equal(t, int64(0), messageID)
}
