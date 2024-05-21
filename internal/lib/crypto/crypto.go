package crypto

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GetPasswordHash(password string) (string, error) {
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

func CheckPasswordHash(password, hash string) error {
	// Compare the hash with the plain-text password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("invalid credentials") // return an error if the password does not match
	}
	return nil
}
