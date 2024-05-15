package repository

import (
	"github.com/ed16/messenger/pkg/domain"
)

type MockUserRepository struct {
	GetUserByUsernameFunc func(username string) (*domain.User, error)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	return m.GetUserByUsernameFunc(username)
}

func (m *MockUserRepository) InsertUser(user *domain.User) error {
	return m.InsertUser(user)
}
