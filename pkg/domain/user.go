package domain

// User defines the structure for user credentials
type User struct {
	Username string
	Password string
}

// NewUser creates a new User instance
func NewUser(username, password string) *User {
	return &User{Username: username, Password: password}
}
