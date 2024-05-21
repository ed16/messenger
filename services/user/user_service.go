package user

import (
	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/lib/crypto"
)

type UserRepository interface {
	InsertUser(user *domain.User) error
	GetUserByUsername(username string) (domain.User, error)
	GetUserByID(userID int64) (domain.User, error)
	UpdateUser(user *domain.User) error
	UpdateUserProfile(profile *domain.Profile) error
	GetUsersByUsername(username string) ([]domain.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(ur UserRepository) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

func (s *UserService) CreateNewUser(username, password string) error {
	user := &domain.User{}
	user.Username = username

	passwordHash, err := crypto.GetPasswordHash(password)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash

	return s.userRepo.InsertUser(user)
}

func (s *UserService) AddContact(userID int64, contactUsername string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	contact, err := s.userRepo.GetUserByUsername(contactUsername)
	if err != nil {
		return err
	}
	user.Contacts = append(user.Contacts, contact)
	return s.userRepo.UpdateUser(&user)
}

func (s *UserService) GetUserContacts(userID int64) ([]domain.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user.Contacts, nil
}

func (s *UserService) UpdateUserProfile(profile *domain.Profile) error {
	return nil
}

func (s *UserService) GetUsersByUsername(username string) ([]domain.User, error) {
	return s.userRepo.GetUsersByUsername(username)
}
