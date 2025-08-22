package passwordHashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"golang.org/x/crypto/argon2"
)

// KriptoArgon2ID is a struct that holds the parameters for Argon2 ID hashing
// and encapsulates the Argon2 ID hashing logic.
type KriptoArgon2ID struct {
	saltLength           uint8  // length of the salt in bytes
	timeCost             uint32 // number of iterations over the memory cost
	memoryCostKiloBytes  uint32 // memory cost in kilobytes
	threads              uint8  // number of threads used
	outputkeyLengthBytes uint32 // length of the key in bytes
	pepperSecret         []byte // pepper secret for password hashing
}

// NewKriptoArgon2ID creates a new instance of KriptoArgon2Id with the specified parameters.
// The parameters are validated before creating the instance.
//
// Parameters:
// - saltLength: The length of the salt in bytes.
// - iterations: The number of iterations over the memory cost.
// - memoryCostMegaBytes: The memory cost in megabytes.
// - threads: The number of threads used.
// - keyLengthB: The length of the key in bytes.
// - pepperSecret: The pepper secret for password hashing.
//
// Rationale for Argon2id Parameters:
//
// The following parameters are set according to modern cryptographic best practices,
// primarily from OWASP and NIST, to ensure a strong defense against offline password guessing attacks.
//
//   - The memory cost (memoryCostMegaBytes) should be as high as possible without causing a denial-of-service on the system.
//     A recommended starting point is 64 MB as this provides a strong memory-hard defense against GPU-based attacks.
//
//   - The number of iterations (iterations) and threads (threads) should be tuned to keep the hashing time
//     within an acceptable range for a user login (typically between 0.5 and 1.0 seconds).
//     Increasing these values linearly increases the time it takes for an attacker to perform a single guess.
//
//   - The salt length (saltLength) should be sufficient to prevent attacks that rely on pre-computed hashes.
//     OWASP recommends a minimum salt length of 16 bytes. NIST SP 800-63B recommends a minimum	salt length of 4 bytes.
//     There's no need to increase the salt length beyond 32 bytes because with because with a high entropy salt
//     it's unlikely to be pre-computed by an attacker it has 2^256 possible values.
//
//   - The key length (keyLengthB) is the hash output length. OWASP recommends it should be at least 32 bytes to
//     comply with modern standards and provide a secure output for cryptographic uses (e.g., AES-256 keys).
//
//   - The pepper secret (pepperSecret) is a shared secret for all password hashes. According to
//     NIST SP 800-63B, cryptographic secrets should have a strength of at least 128 bits.
//     Therefore, the pepper secret should be at least 16 bytes long.
//     NIST SP 800-63B it mantioned as secret key not papper :).
//     The pepper secret should be stored in a secure location like a secret manager service or a secure environment variable.
//
// Returns:
// - *KriptoArgon2Id: A new instance of KriptoArgon2Id.
// - *models.CustomError: An error if the parameters are invalid.
//
// Error Conditions:
// - if the parameters are invalid, (ERROR_CODE_INTERNAL_SERVER) because it's a setup error
//
// Example Usage:
//
//	argon2Id, err := NewKriptoArgon2ID(64, 5, 64, 1, 32, pepperSecret)
//
//	if err != nil {
//		return err
//	}
func NewKriptoArgon2ID(saltLength uint8, iterations uint32, memoryCostMegaBytes uint32, threads uint8, keyLengthB uint32, pepperSecret []byte) (*KriptoArgon2ID, *models.CustomError) {
	if saltLength < 16 {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Invalid Argon2 ID parameter saltLength it must be at least 16 for security reason", nil, nil)
	}
	if memoryCostMegaBytes < 64 {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Invalid Argon2 ID parameter memoryCostKiloBytes it must be at least 64 megabytes", nil, nil)
	}
	if keyLengthB < 32 {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Invalid Argon2 ID parameter keyLengthB it must be at least 32 for security reason", nil, nil)
	}
	if len(pepperSecret) < 16 {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Invalid Argon2 ID parameter pepperSecret it must be at least 16 for security reason", nil, nil)
	}

	return &KriptoArgon2ID{saltLength: saltLength, timeCost: iterations, memoryCostKiloBytes: memoryCostMegaBytes * 1024, threads: threads, outputkeyLengthBytes: keyLengthB, pepperSecret: pepperSecret}, nil
}

// generateSalt generates a random salt.
// The function uses crypto/rand a CSPRNG to generate a random salt of the specified length.
// The salt is cryptographically secure complying with NIST SP 800-90.
//
// Parameters:
// - saltLength: The length of the salt to generate in bytes.
//
// Returns:
// - []byte: The generated salt.
// - *models.CustomError: An error if the salt generation fails.
//
// Error Conditions:
// - if the salt generation fails,(ERROR_CODE_INTERNAL_SERVER)
// Example Usage:
//
//	salt, err := generateSalt()
//	if err != nil {
//		return err
//	}
func (a *KriptoArgon2ID) generateSalt() ([]byte, *models.CustomError) {

	salt := make([]byte, a.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to generate salt for password", &err, nil)
	}
	return salt, nil
}

// GenerateNew generates an Argon2 ID hash for the given password.
//
// Functionality:
//  1. Generate a random salt.
//  2. Append the pepper secret to the password.
//  3. Hash the password using Argon2 ID.
//  4. Format the Argon2 ID hash.
//
// Parameters:
// - password: The password to hash.
//
// Returns:
// - string: The Argon2 ID hash of the password.
// - *models.CustomError: An error if the hash generation fails.
//
// Error Conditions:
// - if the Argon2 ID parameters are not set
// - if the salt generation fails
// Example Usage:
//
//	argon2IdHash, err := KriptoArgon2ID.GenerateNew(password)
//	if err != nil {
//		return err
//	}
func (a *KriptoArgon2ID) GenerateNew(password string) (string, *models.CustomError) {

	if a == nil {
		return "", models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to get the underlying Argon2 ID parameters because they are nil", nil, nil)
	}
	const (
		argon2Version = argon2.Version
	)
	salt, err := a.generateSalt()
	if err != nil {
		return "", err
	}

	passwordBytes := []byte(password)
	defer func() {
		for i := range passwordBytes {
			passwordBytes[i] = 0
		}
	}()
	pepperedPassword := append(passwordBytes, a.pepperSecret...)
	defer func() {
		for i := range pepperedPassword {
			pepperedPassword[i] = 0
		}
	}()

	argon2IdHash := argon2.IDKey(pepperedPassword, salt, a.timeCost, a.memoryCostKiloBytes, a.threads, a.outputkeyLengthBytes)

	if argon2IdHash == nil {
		return "", models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to generate password hash", nil, nil)
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(argon2IdHash)

	fullHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2Version,
		a.memoryCostKiloBytes,
		a.timeCost,
		a.threads,
		b64Salt,
		b64Hash,
	)
	return fullHash, nil
}

// Verify verifies the given password against the provided Argon2 ID hash.
//
// Functionality:
//  1. Parse the Argon2 ID hash into its components.
//  2. Append the pepper secret to the password.
//  3. Hash the password using Argon2 ID.
//  4. Verify the hash against the Argon2 ID hash.
//
// Parameters:
// - password: The password to verify.
// - fullhash: The Argon2 ID hash to verify against.
//
// Returns:
// - *models.CustomError: An error if the verification fails.
//
// Error Conditions:
// - if the Argon2 ID parameters are not set
// - if the hash is not in the expected format
// - if the password verification fails,(ERROR_CODE_UNAUTHORIZED)
// Example Usage:
//
//	if err := Verify(password, argon2IdHash);err != nil {
//		return err
//	}
func (a *KriptoArgon2ID) Verify(password string, fullhash string) *models.CustomError {
	if a == nil {
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to get the underlying Argon2 ID parameters because they are nil", nil, nil)
	}

	argon2IdHash, err := getArgon2IdHashParts(fullhash)
	if err != nil {
		return err
	}

	pepperedPassword := append([]byte(password), a.pepperSecret...)
	defer func() {
		for i := range pepperedPassword {
			pepperedPassword[i] = 0
		}
	}()
	if argon2IdHash.costFactors["p"] < 255 {
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to parse Argon2 threads from hash", nil, nil)

	}

	providedHash := argon2.IDKey(
		pepperedPassword,
		argon2IdHash.Salt,
		argon2IdHash.costFactors["t"],
		argon2IdHash.costFactors["m"],
		uint8(argon2IdHash.costFactors["p"]),
		uint32(len(argon2IdHash.Hash)),
	)

	if subtle.ConstantTimeCompare(providedHash, argon2IdHash.Hash) != 1 {
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Password verification failed", nil, nil)
	}

	return nil
}

// argon2IdHash represents an Argon2 ID hash.
type argon2IdHash struct {
	Argon2Version int
	costFactors   map[string]uint32
	Salt          []byte
	Hash          []byte
}

// String returns a string representation of the Argon2 ID hash.
func (a *argon2IdHash) String() string {
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

// newArgon2IdHash creates a new Argon2 ID hash from the provided parameters.
//
// Parameters:
// - version: The Argon2 version.
// - costFactors: A map of cost factors for memory, time, and threads.
// - salt: The salt value.
// - hash: The hash value.
//
// Returns:
// - *Argon2IdHash: A pointer to an Argon2IdHash struct.
func newArgon2IdHash(version int, costFactors map[string]uint32, salt []byte, hash []byte) *argon2IdHash {
	return &argon2IdHash{
		Argon2Version: version,
		costFactors:   costFactors,
		Salt:          salt,
		Hash:          hash,
	}
}

// getArgon2IdHashParts splits the Argon2 ID hash into its components.
//
// Parameters:
// - fullhash: The Argon2 ID hash to split.
//
// Returns:
// - *Argon2IdHash: A pointer to an Argon2IdHash struct containing the split components.
// - *models.CustomError: An error if the hash is not in the expected format.
//
// Error Conditions:
// - if the hash is not in the expected format,(ERROR_CODE_UNPROCESSABLE_ENTITY)
// Example Usage:
//
//	argon2IdHash, err := getArgon2IdHashParts(argon2IdHashString)
//	if err != nil {
//		return err
//	}
func getArgon2IdHashParts(fullhash string) (*argon2IdHash, *models.CustomError) {
	parts := strings.Split(fullhash, "$")
	if len(parts) != 6 {
		return nil, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, "Failed to get hash parts from Argon2 ID hash($)", nil, nil)
	}

	version, err := strconv.Atoi(strings.Split(parts[2], "=")[1])
	if err != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, "Failed to parse Argon2 version", &err, nil)
	}

	costFactors := strings.Split(parts[3], ",")
	if len(costFactors) != 3 {
		return nil, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, "Failes to get cost factors from Argon2 ID hash(,)", nil, nil)
	}
	costFactorsMap := make(map[string]uint32)

	for _, costFactor := range costFactors {
		costFactorParts := strings.Split(costFactor, "=")
		if len(costFactorParts) != 2 {
			return nil, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, "Failed to get cost factor and value from Argon2 ID hash(=)", nil, nil)
		}

		value, err := strconv.ParseUint(costFactorParts[1], 10, 32)
		if err != nil {
			return nil, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, fmt.Sprintf("Failed to parse cost factor value %s", costFactorParts[1]), &err, nil)
		}
		costFactorsMap[costFactorParts[0]] = uint32(value)

	}

	salt, saltErr := base64.RawStdEncoding.DecodeString(parts[4])
	if saltErr != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, "Failed to generate salt for password", &saltErr, nil)
	}

	hash, hashErr := base64.RawStdEncoding.DecodeString(parts[5])
	if hashErr != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, "Failed to generate password hash", &hashErr, nil)
	}

	argon2IdHash := newArgon2IdHash(version, costFactorsMap, salt, hash)

	return argon2IdHash, nil
}
