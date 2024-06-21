package domain

import (
	"time"
)

type MessageStatus string

const (
	MessageStatusSent     MessageStatus = "sent"
	MessageStatusReceived MessageStatus = "received"
	MessageStatusRead     MessageStatus = "read"
)

type Message struct {
	MessageId   int64
	SenderId    int64
	RecipientId int64
	Content     string
	CreatedAt   time.Time
	MediaId     *int64
	Status      MessageStatus
}
