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

func (m *MessageRepository) CreateMessage(ctx context.Context, message *domain.Message) (messageId int64, err error) {
	query := `
		INSERT INTO messages (sender_id, recipient_id, content, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING message_id
	`

	err = m.DB.QueryRowContext(ctx, query, message.SenderId, message.RecipientId, message.Content, message.Status).Scan(&messageId)
	if err != nil {
		return 0, err
	}

	return messageId, nil
}

func (m *MessageRepository) GetMessagesByUserId(ctx context.Context, userId int64) ([]domain.Message, error) {
	countQuery := `
	SELECT COUNT(*)
	FROM messages m
	WHERE m.sender_id = $1 OR m.recipient_id = $1;`

	var count int
	err := m.DB.QueryRowContext(ctx, countQuery, userId).Scan(&count)
	if err != nil {
		return nil, err
	}

	query := `
	SELECT 
		message_id, 
		sender_id,
		recipient_id, 
		content,
		created_at, 
		status,
		media_id
	FROM messages m
	WHERE m.sender_id = $1 OR m.recipient_id = $1
	ORDER BY m.message_id ASC
	LIMIT 100;`

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize the slice with the exact capacity
	messages := make([]domain.Message, 0, count)
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(&message.MessageId, &message.SenderId, &message.RecipientId, &message.CreatedAt,
			&message.Content, &message.Status, &message.MediaId)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
