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

func (m *MockUserRepository) GetUserByID(userID int64) (*domain.User, error) {
	return m.GetUserByID(userID)
}

func (m *MockUserRepository) UpdateUser(user *domain.User) error {
	return m.UpdateUser(user)
}

func (m *MockUserRepository) UpdateUserProfile(profile *domain.Profile) error {
	return m.UpdateUserProfile(profile)
}

func (m *MockUserRepository) GetUsersByUsername(username string) ([]*domain.User, error) {
	return m.GetUsersByUsername(username)
}
