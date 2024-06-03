package domain

import "time"

type User struct {
	UserId       int64     `json:"UserId"`
	Username     string    `json:"Username"`
	Status       byte      `json:"Status,omitempty"`
	CreatedAt    time.Time `json:"CreatedAt,omitempty"`
	PasswordHash string    `json:"PasswordHash,omitempty"`
}
