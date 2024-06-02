package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/lib/crypto"
	"github.com/golang-jwt/jwt/v4"
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

type AuthService struct {
	userRepo  UserRepository
	secretKey string
}

// Custom claims structure, add any fields that you might need in your token
type CustomClaims struct {
	jwt.RegisteredClaims
	// We can add custom fields for additional claims, e.g., UserRole
	UserRole string
}

func NewAuthService(ur UserRepository) *AuthService {
	return &AuthService{
		userRepo: ur,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = crypto.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		return "", err
	}

	tokenString, err := s.GetToken(ctx, user.UserId)
	return tokenString, err
}

func (s *AuthService) ValidateToken(tokenString string) (userId int64, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Additional claim checks can be added here, e.g., role verification
		userId, _ = strconv.ParseInt(claims.Subject, 10, 64)
		return userId, nil
	} else {
		return 0, fmt.Errorf("invalid token")
	}
}

func (s *AuthService) GetToken(ctx context.Context, userId int64) (string, error) {
	// Prepare the claims of the token
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userId),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserRole: "admin", // Example role
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
