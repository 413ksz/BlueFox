package passwordHashing

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const MAXPASSWORDLENGTH = 72

// HashPassword generates a bcrypt hash for the given password.
//
// Parameters:
// - password: The password to hash.
//
// Returns:
// - string: The hashed password.
// - *models.CustomError: An error if the operation fails, otherwise nil.
//
// Error Conditions:
// - if the password is too long to hash(ERROR_CODE_UNPROCESSABLE_ENTITY).
// - if an error occurs during the hashing process(ERROR_CODE_INTERNAL_SERVER).
//
// Example:
//
//	password := "password123"
//	hashedPassword, err := passwordhashing.HashPassword(password)
//	if err != nil {
//		return err
//	}
func HashPassword(password string) (string, *models.CustomError) {
	if len(password) > MAXPASSWORDLENGTH {
		return "", models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, "Password is too long to hash", nil, nil)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "error hashing password", &err, nil)
	}
	return string(bytes), nil
}

// VerifyPassword verifies if the given password matches the stored hash.
//
// Parameters:
// - password: The password to verify.
// - hash: The stored hash to compare against.
//
// Returns:
// - *models.CustomError: An error if the operation fails, otherwise nil.
//
// Error Conditions:
// - if the password does not match the stored hash(ERROR_CODE_UNAUTHORIZED).
//
// Example:
// password := "password123"
// hash := "hashedPassword123"
//
//	if err := passwordhashing.VerifyPassword(password, hash); err != nil {
//		return err
//	}
func VerifyPassword(password, hash string) *models.CustomError {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "invalid password", &err, nil)
	}
	return nil
}
