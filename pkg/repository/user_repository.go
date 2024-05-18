package repository

import (
	"errors"

	"github.com/ed16/messenger/pkg/domain"
)

type UserRepoImpl struct{}

type UserRepository interface {
	InsertUser(user *domain.User) error
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByID(userID int64) (*domain.User, error)
	UpdateUser(user *domain.User) error
	GetUsersByUsername(username string) ([]*domain.User, error) // Added method
}

func (m *UserRepoImpl) InsertUser(user *domain.User) error {
	return nil
}

func (m *UserRepoImpl) GetUserByUsername(username string) (*domain.User, error) {
	user := domain.NewUser()
	if username == "admin" {
		user.UserId = 1
		user.Username = "admin"
		user.PasswordHash = "$2a$10$p7X62PHGUAGFnhdBDLFjs.ufDZY.59FbWlrBi1PxG4OKlHEb.lTVO"
	}
	return user, nil
}

func (m *UserRepoImpl) GetUserByID(userID int64) (*domain.User, error) {
	user := domain.NewUser()
	if userID == 1 {
		user.UserId = 1
		user.Username = "user1"
		user.Contacts = []*domain.User{}
	} else {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *UserRepoImpl) UpdateUser(user *domain.User) error {
	return nil
}

// GetUsersByUsername makes search with the like condition <*username*> Status = 0
func (m *UserRepoImpl) GetUsersByUsername(username string) ([]*domain.User, error) {
	var users []*domain.User

	// Simulate fetching users by username
	if username == "admin" {
		user := domain.NewUser()
		user.UserId = 1
		user.Username = "admin"
		users = append(users, user)
	}
	return users, nil
}
