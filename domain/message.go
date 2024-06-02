package domain

import "time"

type Message struct {
	MessageId   int64
	SenderId    int64
	RecipientId int64
	Content     string
	CreatedAt   time.Time
	IsRead      bool
	IsReceived  bool
	MediaId     *int64
}
