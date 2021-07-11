package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plaintext password and returns a hashed password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hash password: %w", err)
	}
	return string(hash), nil
}

// ValidatePassword takes a plaintext password and checks if it is the same
// as the hashed password
func ValidatePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
