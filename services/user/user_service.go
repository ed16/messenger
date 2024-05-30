package user

import (
	"context"
	"fmt"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/lib/crypto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, userID int64) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error)
	CreateUserContact(ctx context.Context, contact *domain.Contact) error
	GetUserContactsByUserID(ctx context.Context, userID int64) ([]domain.Contact, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(ur UserRepository) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username, password string) error {
	user := &domain.User{}
	user.Username = username

	passwordHash, err := crypto.GetPasswordHash(password)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash
	user.Status = 1 // TODO: Implement status schema

	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) AddContact(ctx context.Context, userID int64, contactUsername string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	contactUser, err := s.userRepo.GetUserByUsername(ctx, contactUsername)
	if err != nil {
		return fmt.Errorf("Error while searching contact user: %v", err)
	}
	if user.UserId == contactUser.UserId {
		return fmt.Errorf("It is not possible to add youself as a contact")
	}

	contact := &domain.Contact{
		UserId:        user.UserId,
		ContactUserId: contactUser.UserId,
	}

	return s.userRepo.CreateUserContact(ctx, contact)
}

func (s *UserService) GetUserContacts(ctx context.Context, userID int64) ([]domain.Contact, error) {
	return s.userRepo.GetUserContactsByUserID(ctx, userID)
}

func (s *UserService) GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	return s.userRepo.GetUsersByUsername(ctx, username)
}
