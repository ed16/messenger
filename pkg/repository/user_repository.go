package repository

import (
	"github.com/ed16/messenger/pkg/domain"
)

// MockUserRepository is a mock repository for demonstration
type MockUserRepository struct{}

func (m *MockUserRepository) FindByUsername(username string) *domain.User {
	if username == "admin" {
		return domain.NewUser(username, "admin")
	}
	return nil
}
