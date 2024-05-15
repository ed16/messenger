// auth_service_test.go
package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/ed16/messenger/pkg/domain"
	"github.com/ed16/messenger/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthenticate(t *testing.T) {
	mockRepo := &repository.MockUserRepository{
		GetUserByUsernameFunc: func(username string) (*domain.User, error) {
			if username == "validUser" {
				return &domain.User{
					UserId:       1,
					Username:     "validUser",
					PasswordHash: "$2a$10$ewIvBUkJThkiNpNZspZ9COyZCpgBG7WK/9pWWrtLgx4ZJp2RXGvu.", // "password" hashed
				}, nil
			}
			return nil, errors.New("invalid credentials")
		},
	}

	authService := AuthService{
		UserRepository: mockRepo,
		SecretKey:      "testSecretKey",
	}

	t.Run("valid credentials", func(t *testing.T) {
		token, err := authService.Authenticate("validUser", "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		token, err := authService.Authenticate("invalidUser", "password")
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("invalid password", func(t *testing.T) {
		token, err := authService.Authenticate("validUser", "wrongpassword")
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestParseToken(t *testing.T) {
	authService := AuthService{
		SecretKey: "testSecretKey",
	}

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserRole: "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(authService.SecretKey))
	assert.NoError(t, err)

	t.Run("valid token", func(t *testing.T) {
		parsedClaims, err := authService.ParseToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, claims.Subject, parsedClaims.Subject)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := authService.ParseToken("invalidToken")
		assert.Error(t, err)
	})
}

func TestGetPasswordHash(t *testing.T) {
	authService := AuthService{}

	t.Run("valid password", func(t *testing.T) {
		hash, err := authService.GetPasswordHash("password")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("empty password", func(t *testing.T) {
		hash, err := authService.GetPasswordHash("")
		assert.Error(t, err)
		assert.Empty(t, hash)
	})
}

func TestCheckPasswordHash(t *testing.T) {
	authService := AuthService{}
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	t.Run("valid password", func(t *testing.T) {
		err := authService.CheckPasswordHash("password", string(hash))
		assert.NoError(t, err)
	})

	t.Run("invalid password", func(t *testing.T) {
		err := authService.CheckPasswordHash("wrongpassword", string(hash))
		assert.Error(t, err)
	})
}

func TestGetToken(t *testing.T) {
	authService := AuthService{
		SecretKey: "testSecretKey",
	}

	t.Run("generate token", func(t *testing.T) {
		token, err := authService.GetToken(1)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
