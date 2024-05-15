package repository

import (
	"github.com/ed16/messenger/pkg/domain"
)

type MockTokenRepository struct{}

type TokenRepository interface {
	InsertNewToken(token string) error
	GetTokenByValue(token string) (*domain.UserToken, error)
}

func (m *MockTokenRepository) InsertNewToken(token string) error {
	return nil
}

func (m *MockTokenRepository) GetTokenByValue(token string) (*domain.UserToken, error) {
	userToken := domain.NewUserToken()
	return userToken, nil
}
