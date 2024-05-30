package crypto

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// TestGetPasswordHash tests the GetPasswordHash function.
func TestGetPasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"ValidPassword", "validpassword", false},
		{"EmptyPassword", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GetPasswordHash(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check that the hash is valid
				if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password)); err != nil {
					t.Errorf("Generated hash does not match password, error = %v", err)
				}
			}
		})
	}
}

// TestCheckPasswordHash tests the CheckPasswordHash function.
func TestCheckPasswordHash(t *testing.T) {
	hash, err := GetPasswordHash("validpassword")
	if err != nil {
		t.Fatalf("GetPasswordHash() error = %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{"ValidPassword", "validpassword", hash, false},
		{"InvalidPassword", "invalidpassword", hash, true},
		{"EmptyPassword", "", hash, true},
		{"EmptyHash", "validpassword", "", true},
		{"EmptyPasswordAndHash", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
