package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const MAXPASSWORDLENGTH = 72

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	if len(password) > MAXPASSWORDLENGTH {
		return "", errors.New("password length exceeds 72 bytes (bcrypt limit)")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
