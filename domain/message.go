package domain

import "time"

type Message struct {
	MessageId   int64
	SenderId    int64
	RecipientId int64
	Content     string
	CreatedAt   time.Time
	IsRead      byte
	IsReceived  byte
	MedeiaId    int64
}
