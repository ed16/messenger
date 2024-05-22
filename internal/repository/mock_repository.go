package repository

import (
	"context"

	"github.com/ed16/messenger/domain"
)

type MockUserRepository struct {
	GetUserByUsernameFunc func(ctx context.Context, username string) (domain.User, error)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	return m.GetUserByUsernameFunc(ctx, username)
}

func (m *MockUserRepository) InsertUser(ctx context.Context, user *domain.User) error {
	return m.InsertUser(ctx, user)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID int64) (domain.User, error) {
	return m.GetUserByID(ctx, userID)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return m.UpdateUser(ctx, user)
}

func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, profile *domain.Profile) error {
	return m.UpdateUserProfile(ctx, profile)
}

func (m *MockUserRepository) GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	return m.GetUsersByUsername(ctx, username)
}
