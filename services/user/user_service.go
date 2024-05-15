package user

import (
	"github.com/ed16/messenger/pkg/domain"
	"github.com/ed16/messenger/pkg/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func (s *UserService) CreateNewUser() error {
	user := domain.NewUser()
	s.UserRepository.InsertUser(user)
	return nil
}
