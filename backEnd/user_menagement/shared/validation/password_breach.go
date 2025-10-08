package validation

import (
	"bufio"
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
)

// Doer is an interface for making HTTP requests.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// PwnedPassword encapsulates the logic for interacting with the Pwned Passwords API.
type PwnedPassword struct {
	client Doer
}

// NewPwnedPassword creates a new instance of PwnedPassword with the provided HTTP client.
func NewPwnedPassword(client Doer) *PwnedPassword {
	return &PwnedPassword{client: client}
}

// CheckPasswordBreach checks if a password has been exposed in a data breach using the Pwned Passwords API.
//
// Functionality:
//  1. Prefix and Suffix Generation: Generates the prefix and suffix of the SHA-1 hash of the password.
//  2. Get Pwned Password Hashes: Retrieves the count of times a password has been exposed in a data breach using the Pwned Passwords API.
//
// Parameters:
//   - password: The password to check.
//
// Returns:
//   - *models.CustomError: An error if the operation fails, otherwise nil.
//
// Error Conditions:
//   - if an error occurs with the Pwned Passwords API (e.g., request timeout).
//
// Example:
//
//	password := "password123"
//	err := pwnedPassword.CheckPasswordBreach(password)
//	if err != nil {
//		return err
//	}
func (p *PwnedPassword) CheckPasswordBreach(password string) *models.CustomError {
	prefix, suffix := p.sha1Cheksum(password)
	hashes, err := p.GetPwnedPasswordHashes(prefix)
	if err != nil {
		return err
	}

	if hashes[suffix] > 0 {
		return models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, "Password has been exposed in a data breach", nil, nil)
	}
	return nil

}

// GetPwnedPasswordHashes retrieves the count of times a password has been exposed in a data breach using the Pwned Passwords API.
//
// Functionality:
//  1. API Request: Makes an HTTP GET request to the Pwned Passwords API to retrieve the count of times a password has been exposed in a data breach.
//  2. Response Parsing: Parses the response from the API and returns a map of password hashes to their exposure count.
//
// Parameters:
//   - identifier: The first 5 characters of the SHA-1 hash of the password.
//
// Returns:
//   - map[string]int: A map of password hashes to their exposure count.
//   - *models.CustomError: An error if the operation fails, otherwise nil.
//
// Error Conditions:
//   - Invalid identifier: Must be exactly 5 characters long (ERROR_CODE_BAD_REQUEST).
//   - External dependency error: An error occurs with the Pwned Passwords API (e.g., request timeout).
//   - Internal server error: An error occurs during a critical operation (e.g., parsing the response from the API).
//
// Example:
//
//	identifier := "CB66F"
//	pawnedHashes, err := pwnedPassword.GetPwnedPasswordHashes(identifier)
//	if err != nil {
//		return err
//	}
func (p *PwnedPassword) GetPwnedPasswordHashes(identifier string) (map[string]int, *models.CustomError) {
	if len(identifier) != 5 {
		return nil, models.NewCustomError(models.ERROR_CODE_BAD_REQUEST, "Invalid identifier: must be 5 characters long", nil, nil)
	}
	const (
		userAgent  = "BlueFox"
		apiTimeout = 10 * time.Second
		apiUrl     = "https://api.pwnedpasswords.com/range/%s"
	)

	url := fmt.Sprintf(apiUrl, identifier)
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to create request for Pwned Password API", err, nil)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := p.client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return nil, models.NewCustomError(models.ERROR_CODE_EXTERNAL_SERVER, "Pwned Password API request timed out", err, nil)
		}
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to make a request to Pwned Password API", err, nil)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return make(map[string]int), nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, models.NewCustomError(models.ERROR_CODE_EXTERNAL_SERVER, "Failed to get response from Pwned Password API", nil, nil)
	}

	scanner := bufio.NewScanner(resp.Body)
	results := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")

		if len(parts) == 2 {
			hashSuffix := parts[0]
			var count int
			count, err := strconv.Atoi(parts[1])
			if err == nil {
				results[hashSuffix] = count
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to scan response from Pwned Password API", err, nil)
	}

	return results, nil
}

// sha1Cheksum generates a SHA-1 hash for the given password and returns the prefix and suffix of the hash.
// The prefix is the first 5 characters of the hash, and the suffix is the remaining characters.
// Functionality:
//  1. hash password: The password is hashed using the SHA-1 algorithm.
//  2. get prefix and suffix: The first 5 characters of the hash are extracted as the prefix and the remaining characters are extracted as the suffix.
//
// Parameters:
// - password: The password to be hashed.
//
// Returns:
// - string: The prefix of the SHA-1 hash.
// - string: The suffix of the SHA-1 hash.
// Example:
//   - password: "password123"
//   - prefix, suffix := PwnedPassword.sha1Cheksum(password)
func (p *PwnedPassword) sha1Cheksum(password string) (string, string) {
	passwordSha1 := sha1.Sum([]byte(password))
	fullHashString := fmt.Sprintf("%X", passwordSha1)
	prefix := fullHashString[:5]
	suffix := fullHashString[5:]
	return prefix, suffix
}
