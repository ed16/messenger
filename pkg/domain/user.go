package domain

import "time"

type User struct {
	UserId       int64
	Username     string // Email
	Status       string
	CreatedAt    time.Time // TIMESTAMP WITH TIME ZONE '2004-10-19 10:23:54+02'
	PasswordHash string
}

// NewUser creates a new User instance
func NewUser() *User {
	return &User{}
}
