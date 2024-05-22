package user

import (
	"context"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/lib/crypto"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *domain.User) error
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByID(ctx context.Context, userID int64) (domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	UpdateUserProfile(ctx context.Context, profile *domain.Profile) error
	GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(ur UserRepository) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

func (s *UserService) CreateNewUser(ctx context.Context, username, password string) error {
	user := &domain.User{}
	user.Username = username

	passwordHash, err := crypto.GetPasswordHash(password)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash

	return s.userRepo.InsertUser(ctx, user)
}

func (s *UserService) AddContact(ctx context.Context, userID int64, contactUsername string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	contact, err := s.userRepo.GetUserByUsername(ctx, contactUsername)
	if err != nil {
		return err
	}
	user.Contacts = append(user.Contacts, contact)
	return s.userRepo.UpdateUser(ctx, &user)
}

func (s *UserService) GetUserContacts(ctx context.Context, userID int64) ([]domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.Contacts, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, profile *domain.Profile) error {
	return nil
}

func (s *UserService) GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	return s.userRepo.GetUsersByUsername(ctx, username)
}
