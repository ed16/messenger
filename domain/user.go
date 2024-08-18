package domain

import "time"

type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusDeleted UserStatus = "deleted"
	UserStatusBlocked UserStatus = "blocked"
)

// User represents a user in the system
type User struct {
	UserId       int64      `json:"UserId"`
	Username     string     `json:"Username"`
	Status       UserStatus `json:"Status,omitempty"`
	CreatedAt    time.Time  `json:"CreatedAt,omitempty"`
	PasswordHash string     `json:"PasswordHash,omitempty"`
}
