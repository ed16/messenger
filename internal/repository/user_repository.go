package repository

import (
	"database/sql"
	"errors"

	"github.com/ed16/messenger/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (m *UserRepository) InsertUser(user *domain.User) error {
	return nil
}

func (m *UserRepository) GetUserByUsername(username string) (domain.User, error) {
	user := domain.User{}
	if username == "admin" {
		user.UserId = 1
		user.Username = "admin"
		user.PasswordHash = "$2a$10$p7X62PHGUAGFnhdBDLFjs.ufDZY.59FbWlrBi1PxG4OKlHEb.lTVO"
	}
	return user, nil
}

func (m *UserRepository) GetUserByID(userID int64) (domain.User, error) {
	user := domain.User{}
	if userID == 1 {
		user.UserId = 1
		user.Username = "user1"
		user.Contacts = []domain.User{}
	} else {
		return user, errors.New("user not found")
	}
	return user, nil
}

// Add contact or change password
func (m *UserRepository) UpdateUser(user *domain.User) error {
	return nil
}

func (m *UserRepository) UpdateUserProfile(profile *domain.Profile) error {
	return nil
}

// GetUsersByUsername makes search with the like condition <*username*> Status = 0
func (m *UserRepository) GetUsersByUsername(username string) ([]domain.User, error) {
	users := make([]domain.User, 0)
	// Simulate fetching users by username
	if username == "admin" {
		user := domain.User{}
		user.UserId = 1
		user.Username = "admin"
		users = append(users, user)
	}
	return users, nil
}
