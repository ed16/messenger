package auth

import (
	"errors"

	"github.com/ed16/messenger/pkg/domain"
)

type AuthService struct {
	UserRepository UserRepository
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	user := s.UserRepository.FindByUsername(username)
	if user != nil && user.Password == password {
		return "somejwttoken", nil // Normally, you'd use a JWT library here
	}
	return "", errors.New("invalid credentials")
}

type UserRepository interface {
	FindByUsername(username string) *domain.User
}
