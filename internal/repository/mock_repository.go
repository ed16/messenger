package repository

import (
	"context"

	"github.com/ed16/messenger/domain"
)

type MockUserRepository struct {
	GetUserByUsernameFunc func(ctx context.Context, username string) (domain.User, error)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return m.CreateUser(ctx, user)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return m.UpdateUser(ctx, user)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID int64) (domain.User, error) {
	return m.GetUserByID(ctx, userID)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	return m.GetUserByUsernameFunc(ctx, username)
}

func (m *MockUserRepository) GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	return m.GetUsersByUsername(ctx, username)
}

func (m *MockUserRepository) CreateUserContact(ctx context.Context, contact *domain.Contact) error {
	return m.CreateUserContact(ctx, contact)
}

func (m *MockUserRepository) GetUserContactsByUserID(ctx context.Context, userID int64) ([]domain.Contact, error) {
	return m.GetUserContactsByUserID(ctx, userID)
}
