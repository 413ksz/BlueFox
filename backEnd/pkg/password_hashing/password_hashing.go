package passwordHashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"golang.org/x/crypto/argon2"
)

const (
	argon2KeyLengthBytes = 64 // 64 MB key output size
)

// GenerateSalt generates a random salt.
//
// Parameters:
// - None
//
// Returns:
// - []byte: The generated salt.
// - *models.CustomError: An error if the salt generation fails.
//
// Error Conditions:
// - if the salt generation fails,(ERROR_CODE_INTERNAL_SERVER)
// Example Usage:
//
//	salt, err := GenerateSalt()
//	if err != nil {
//		return err
//	}
func GenerateSalt() ([]byte, *models.CustomError) {
	const n = 16

	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to generate salt for password", &err, nil)
	}
	return bytes, nil
}

// GetPepper retrieves the paper secret key from the environment variable.
//
// Parameters:
// - None
//
// Returns:
// - []byte: The paper secret key.
// - *models.CustomError: An error if the paper secret key is not found.
//
// Error Conditions:
// - if the paper secret key is not found,(ERROR_CODE_INTERNAL_SERVER)
// Example Usage:
//
//	papper, err := GetPepper()
//	if err != nil {
//		return err
//	}
func GetPepper() ([]byte, *models.CustomError) {
	pepper := os.Getenv("PEPPER_SECRET_KEY")
	if pepper == "" {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to get paper secret key", nil, nil)
	}
	return []byte(pepper), nil
}

// GenerateNewArgon2IdHash generates an Argon2 ID hash for the given password.
//
// Parameters:
// - password: The password to hash.
//
// Returns:
// - string: The Argon2 ID hash of the password.
// - *models.CustomError: An error if the hash generation fails.
//
// Error Conditions:
// - if the salt generation fails
// - if the paper secret key is not found
// Example Usage:
//
//	argon2IdHash, err := GenerateNewArgon2IdHash(password)
//	if err != nil {
//		return err
//	}
func GenerateNewArgon2IdHash(password string) (string, *models.CustomError) {
	const (
		argon2Version = argon2.Version
		timeCostS     = 5         // 5 seconds iteration
		memoryCostKb  = 64 * 1024 // 64 MB
		threads       = 1         // because of serverless environment
	)
	salt, err := GenerateSalt()
	if err != nil {
		return "", err
	}

	pepper, err := GetPepper()
	if err != nil {
		return "", err
	}

	passwordBytes := []byte(password)
	pepperedPassword := append(passwordBytes, pepper...)

	argon2IdHash := argon2.IDKey(pepperedPassword, salt, timeCostS, memoryCostKb, threads, argon2KeyLengthBytes)

	if argon2IdHash == nil {
		return "", models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to generate password hash", nil, nil)
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(argon2IdHash)

	fullHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2Version,
		memoryCostKb,
		timeCostS,
		threads,
		b64Salt,
		b64Hash,
	)
	return fullHash, nil
}

// VerifyPassword verifies the given password against the provided Argon2 ID hash.
//
// Parameters:
// - password: The password to verify.
// - fullhash: The Argon2 ID hash to verify against.
//
// Returns:
// - *models.CustomError: An error if the verification fails.
//
// Error Conditions:
// - if the hash is not in the expected format,(ERROR_CODE_INTERNAL_SERVER)
// - if the paper secret key is not found
// - if the password verification fails,(ERROR_CODE_UNAUTHORIZED)
// Example Usage:
//
//	if err := VerifyPassword(password, argon2IdHash);err != nil {
//		return err
//	}
func VerifyPassword(password string, fullhash string) *models.CustomError {

	argon2IdHash, err := GetArgon2IdHashParts(fullhash)
	if err != nil {
		return err
	}

	pepper, err := GetPepper()
	if err != nil {
		return err
	}

	pepperedPassword := append([]byte(password), pepper...)

	providedHash := argon2.IDKey(
		[]byte(pepperedPassword),
		argon2IdHash.Salt,
		uint32(argon2IdHash.costFactors["t"]),
		uint32(argon2IdHash.costFactors["m"]),
		uint8(argon2IdHash.costFactors["p"]),
		argon2KeyLengthBytes,
	)

	if subtle.ConstantTimeCompare(providedHash, argon2IdHash.Hash) != 1 {
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Password verification failed", nil, nil)
	}

	return nil
}

// Argon2IdHash represents an Argon2 ID hash.
type Argon2IdHash struct {
	Argon2Version int
	costFactors   map[string]int
	Salt          []byte
	Hash          []byte
}

// String returns a string representation of the Argon2 ID hash.
func (a *Argon2IdHash) String() string {
	if a == nil {
		return ""
	}

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		a.Argon2Version,
		a.costFactors["m"],
		a.costFactors["t"],
		a.costFactors["p"],
		base64.RawStdEncoding.EncodeToString(a.Salt),
		base64.RawStdEncoding.EncodeToString(a.Hash),
	)
}

// NewArgon2IdHash creates a new Argon2 ID hash from the provided parameters.
//
// Parameters:
// - version: The Argon2 version.
// - costFactors: A map of cost factors for memory, time, and threads.
// - salt: The salt value.
// - hash: The hash value.
//
// Returns:
// - *Argon2IdHash: A pointer to an Argon2IdHash struct.
func NewArgon2IdHash(version int, costFactors map[string]int, salt []byte, hash []byte) *Argon2IdHash {
	return &Argon2IdHash{
		Argon2Version: version,
		costFactors:   costFactors,
		Salt:          salt,
		Hash:          hash,
	}
}

// GetArgon2IdHashParts splits the Argon2 ID hash into its components.
//
// Parameters:
// - fullhash: The Argon2 ID hash to split.
//
// Returns:
// - *Argon2IdHash: A pointer to an Argon2IdHash struct containing the split components.
// - *models.CustomError: An error if the hash is not in the expected format.
//
// Error Conditions:
// - if the hash is not in the expected format,(ERROR_CODE_INTERNAL_SERVER)
// Example Usage:
//
//	argon2IdHash, err := GetArgon2IdHashParts(argon2IdHashString)
//	if err != nil {
//		return err
//	}
func GetArgon2IdHashParts(fullhash string) (*Argon2IdHash, *models.CustomError) {
	parts := strings.Split(fullhash, "$")
	if len(parts) != 6 {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to get hash parts from Argon2 ID hash($)", nil, nil)
	}

	version, err := StringToInt(parts[2])
	if err != nil {
		return nil, err
	}

	costFactors := strings.Split(parts[3], ",")
	if len(costFactors) != 3 {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failes to get cost factors from Argon2 ID hash(,)", nil, nil)
	}
	costFactorsMap := make(map[string]int)

	for _, costFactor := range costFactors {
		costFactorParts := strings.Split(costFactor, "=")
		if len(costFactorParts) != 2 {
			return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to get cost factor and value from Argon2 ID hash(=)", nil, nil)
		}

		value, err := StringToInt(costFactorParts[1])
		if err != nil {
			return nil, err
		}
		costFactorsMap[costFactorParts[0]] = value

	}

	salt, saltErr := base64.RawStdEncoding.DecodeString(parts[4])
	if saltErr != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to generate salt for password", &saltErr, nil)
	}

	hash, hashErr := base64.RawStdEncoding.DecodeString(parts[5])
	if hashErr != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to generate password hash", &hashErr, nil)
	}

	argon2IdHash := NewArgon2IdHash(version, costFactorsMap, salt, hash)

	return argon2IdHash, nil
}

// StringToInt converts a string to an integer.
//
// Parameters:
// - s: The string to convert.
//
// Returns:
// - int: The integer value of the string.
// - *models.CustomError: An error if the conversion fails.
//
// Error Conditions:
// - if the conversion fails,(ERROR_CODE_INTERNAL_SERVER)
// Example Usage:
//
//	i, err := StringToInt("123")
//	if err != nil {
//		return err
//	}
func StringToInt(s string) (int, *models.CustomError) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, fmt.Sprintf("Failed to convert %s to int", s), &err, nil)
	}
	return i, nil
}
