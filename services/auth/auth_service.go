package auth

import (
	"fmt"
	"time"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/lib/crypto"
	"github.com/golang-jwt/jwt/v4"
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

type AuthService struct {
	UserRepository UserRepository
	SecretKey      string
}

// Custom claims structure, add any fields that you might need in your token
type CustomClaims struct {
	jwt.RegisteredClaims
	// We can add custom fields for additional claims, e.g., UserRole
	UserRole string
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	user, err := s.UserRepository.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	err = crypto.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		return "", err
	}

	tokenString, err := s.GetToken(user.UserId)
	return tokenString, err
}

// Typically we do not need to store JWT tokens on the server side for the primary purpose of session management or authentication.
// However, there are a some scenarios where we might consider storing JWTs or related data on the server:
// 1. Revocation List
// 2. Performance Reasons
// 3. Logging and Auditing
func (s *AuthService) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(s.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Additional claim checks can be added here, e.g., role verification
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func (s *AuthService) GetToken(userId int64) (string, error) {
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
	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
