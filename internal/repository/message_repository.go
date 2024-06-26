package repository

import (
	"context"
	"database/sql"

	"github.com/ed16/messenger/domain"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepository {
	return &MessageRepository{
		DB: db,
	}
}

func (m *MessageRepository) CreateMessage(ctx context.Context, message *domain.Message) (message_id int64, err error) {
	query := `
		INSERT INTO messages (sender_id, recipient_id, content, is_read, is_received)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING message_id
	`

	err = m.DB.QueryRowContext(ctx, query, message.SenderId, message.RecipientId, message.Content, message.IsRead, message.IsReceived).Scan(&message_id)
	if err != nil {
		return 0, err
	}

	return message_id, nil
}

func (m *MessageRepository) GetMessagesByUserId(ctx context.Context, user_id int64) ([]domain.Message, error) {
	query := `
	SELECT 
		message_id, 
		sender_id,
		recipient_id, 
		created_at, 
		content,
		is_read,
		is_received,
		media_id
	FROM messages m
	WHERE m.sender_id = $1 OR m.recipient_id = $1
	ORDER BY m.message_id ASC
	LIMIT 100;`

	rows, err := m.DB.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(&message.MessageId, &message.SenderId, &message.RecipientId, &message.CreatedAt,
			&message.Content, &message.IsRead, &message.IsReceived, &message.MediaId)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return messages, domain.ErrNotFound
	}

	return messages, nil
}
