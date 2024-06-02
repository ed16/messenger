package repository

import (
	"context"

	"github.com/ed16/messenger/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, usr *domain.User) error {
	args := m.Called(ctx, usr)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, usr *domain.User) error {
	args := m.Called(ctx, usr)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID int64) (domain.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) CreateUserContact(ctx context.Context, contact *domain.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserContactsByUserID(ctx context.Context, userID int64) ([]domain.Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Contact), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockMessageRepository is a mock implementation of the MessageRepository interface
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) CreateMessage(ctx context.Context, msg *domain.Message) (int64, error) {
	args := m.Called(ctx, msg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMessageRepository) GetMessagesByUserId(ctx context.Context, user_id int64) ([]domain.Message, error) {
	args := m.Called(ctx, user_id)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Message), args.Error(1)
	}
	return nil, args.Error(1)
}
