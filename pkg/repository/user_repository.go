package repository

import (
	"github.com/ed16/messenger/pkg/domain"
)

type UserRepoImpl struct{}

type UserRepository interface {
	InsertUser(user *domain.User) error
	GetUserByUsername(username string) (*domain.User, error)
}

func (m *UserRepoImpl) InsertUser(user *domain.User) error {
	return nil
}

func (m *UserRepoImpl) GetUserByUsername(username string) (*domain.User, error) {
	user := domain.NewUser()
	//TO BE DELETED: For the testing purpose
	if username == "admin" {
		user.UserId = 1
		user.Username = "admin"
		user.PasswordHash = "$2a$10$p7X62PHGUAGFnhdBDLFjs.ufDZY.59FbWlrBi1PxG4OKlHEb.lTVO"
	}
	return user, nil
}
