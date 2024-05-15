package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/ed16/messenger/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository repository.UserRepository
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
	if user == nil || err != nil {
		return "", errors.New("invalid credentials")
	}

	err = s.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		return "", errors.New("invalid credentials")
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

func (s *AuthService) GetPasswordHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("invalid credentials")
	}

	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Println(string(hashedPassword))
	return string(hashedPassword), nil
}

func (s *AuthService) CheckPasswordHash(password, hash string) error {
	// Compare the hash with the plain-text password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("invalid credentials") // return an error if the password does not match
	}
	return nil
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
