package domain

import "time"

// Contact represents a contact relationship between users
type Contact struct {
	UserId        int64
	ContactUserId int64
	CreatedAt     time.Time
}
