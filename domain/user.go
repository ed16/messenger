package domain

import "time"

type User struct {
	UserId       int64
	Username     string
	Status       byte
	CreatedAt    time.Time
	PasswordHash string
}
