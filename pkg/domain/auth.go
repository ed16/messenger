package domain

import "time"

// NOT USED at the moment
type UserToken struct {
	TokenId   int64
	UserId    int64
	token     string
	Email     string
	Status    string
	ExpiresAt time.Time // TIMESTAMP WITH TIME ZONE '2004-10-19 10:23:54+02'
	CreatedAt time.Time // TIMESTAMP WITH TIME ZONE '2004-10-19 10:23:54+02'
	revoked   bool
}

// NewUser creates a new User instance
func NewUserToken() *UserToken {
	return &UserToken{}
}
