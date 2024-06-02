package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/repository"
	"github.com/ed16/messenger/services/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	username := "testuser"
	password := "password123"

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	err := service.CreateUser(ctx, username, password)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Error(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	username := "testuser"
	password := "password123"

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*domain.User")).Return(errors.New("create user error"))

	err := service.CreateUser(ctx, username, password)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddContact(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	userID := int64(1)
	contactUsername := "contactuser"
	mockUser := domain.User{UserId: userID, Username: "testuser"}
	mockContactUser := domain.User{UserId: 2, Username: contactUsername}

	mockRepo.On("GetUserByID", ctx, userID).Return(mockUser, nil)
	mockRepo.On("GetUserByUsername", ctx, contactUsername).Return(mockContactUser, nil)
	mockRepo.On("CreateUserContact", ctx, mock.AnythingOfType("*domain.Contact")).Return(nil)

	err := service.AddContact(ctx, userID, contactUsername)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddContact_Error(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	userID := int64(1)
	contactUsername := "contactuser"
	mockUser := domain.User{UserId: userID, Username: "testuser"}

	mockRepo.On("GetUserByID", ctx, userID).Return(mockUser, nil)
	mockRepo.On("GetUserByUsername", ctx, contactUsername).Return(domain.User{}, errors.New("user not found"))

	err := service.AddContact(ctx, userID, contactUsername)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddContact_SelfContact(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	userID := int64(1)
	mockUser := domain.User{UserId: userID, Username: "testuser"}

	mockRepo.On("GetUserByID", ctx, userID).Return(mockUser, nil)
	mockRepo.On("GetUserByUsername", ctx, mockUser.Username).Return(mockUser, nil)

	err := service.AddContact(ctx, userID, mockUser.Username)
	assert.Error(t, err)
	assert.Equal(t, "It is not possible to add youself as a contact", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetUserContacts(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	userID := int64(1)
	expectedContacts := []domain.Contact{
		{UserId: userID, ContactUserId: 2},
	}

	mockRepo.On("GetUserContactsByUserID", ctx, userID).Return(expectedContacts, nil)

	contacts, err := service.GetUserContacts(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedContacts, contacts)
	mockRepo.AssertExpectations(t)
}

func TestGetUserContacts_Error(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	userID := int64(1)
	mockRepo.On("GetUserContactsByUserID", ctx, userID).Return([]domain.Contact(nil), errors.New("get contacts error"))

	contacts, err := service.GetUserContacts(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, contacts)
	mockRepo.AssertExpectations(t)
}

func TestGetUsersByUsername(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	username := "testuser"
	expectedUsers := []domain.User{
		{UserId: 1, Username: username},
	}

	mockRepo.On("GetUsersByUsername", ctx, username).Return(expectedUsers, nil)

	users, err := service.GetUsersByUsername(ctx, username)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestGetUsersByUsername_Error(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	username := "testuser"
	mockRepo.On("GetUsersByUsername", ctx, username).Return([]domain.User(nil), errors.New("get users error"))

	users, err := service.GetUsersByUsername(ctx, username)
	assert.Error(t, err)
	assert.Nil(t, users)
	mockRepo.AssertExpectations(t)
}
