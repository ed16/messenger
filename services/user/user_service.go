package user

import (
	"github.com/ed16/messenger/pkg/domain"
	"github.com/ed16/messenger/pkg/repository"
	"github.com/ed16/messenger/services/auth"
)

type UserService struct {
	UserRepository repository.UserRepository
	AuthService    auth.AuthService
}

func (s *UserService) CreateNewUser(username, password string) error {
	user := domain.NewUser()
	user.Username = username

	passwordHash, err := s.AuthService.GetPasswordHash(password)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash

	return s.UserRepository.InsertUser(user)
}

func (s *UserService) AddContact(userID int64, contactID int64) error {
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}
	contact, err := s.UserRepository.GetUserByID(contactID)
	if err != nil {
		return err
	}
	user.Contacts = append(user.Contacts, contact)
	return s.UserRepository.UpdateUser(user)
}

func (s *UserService) GetUserContacts(userID int64) ([]*domain.User, error) {
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user.Contacts, nil
}

func (s *UserService) UpdateUserProfile(userID int64, updatedUser *domain.User) error {
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.Username = updatedUser.Username
	user.Status = updatedUser.Status
	return s.UserRepository.UpdateUser(user)
}

func (s *UserService) GetUsersByUsername(username string) ([]*domain.User, error) {
	return s.UserRepository.GetUsersByUsername(username)
}
