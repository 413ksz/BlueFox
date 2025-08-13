package validation

import (
	"bufio"
	"context"
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

// GetPwnedPasswordHashes retrieves the map of hashes and their occurrence counts that match the given SHA-1 hash prefix.
// It queries the Pwned Passwords API for the provided identifier (first 5 characters of a SHA-1 hash)
// and returns a map where the keys are hash suffixes and the values are occurrence counts.
// If the API call fails or encounters any error, it returns a CustomError.
//
// Parameters:
// - identifier: A string representing the first 5 characters of a SHA-1 password hash.
//
// Returns:
// - map[string]int: A map of hash suffixes to their occurrence counts.
// - *models.CustomError: A CustomError if any error occurs during the API call or response processing.
func (p *PwnedPassword) GetPwnedPasswordHashes(identifier string) (map[string]int, *models.CustomError) {
	// ---------- Build the API request ----------
	// Add a check at the beginning of the function
	if len(identifier) != 5 {
		return nil, models.NewCustomError(models.ERROR_CODE_BAD_REQUEST, "Invalid identifier: must be 5 characters long", nil, nil)
	}
	const (
		userAgent  = "BlueFox"
		apiTimeout = 10 * time.Second
		apiUrl     = "https://api.pwnedpasswords.com/range/%s"
	)

	// Build the URL for the API request.
	url := fmt.Sprintf(apiUrl, identifier)

	// Create a context with a timeout for resource cleanup on timeouts.
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	// Create the API request.
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to create request for Pwned Password API", &err, nil)
	}
	req.Header.Set("User-Agent", userAgent)

	// Make the API request.
	resp, err := p.client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return nil, models.NewCustomError(models.ERROR_CODE_EXTERNAL_SERVER, "Pwned Password API request timed out", &err, nil)
		}
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to make a request to Pwned Password API", &err, nil)
	}
	defer resp.Body.Close() // Close the response body when done.

	// ---------- Process the API response ----------

	// The API returns an HTTP 404 if no hashes match the prefix.
	if resp.StatusCode == http.StatusNotFound {
		return make(map[string]int), nil
	}

	// The API returns an HTTP 200 if the request is successful.
	if resp.StatusCode != http.StatusOK {
		return nil, models.NewCustomError(models.ERROR_CODE_EXTERNAL_SERVER, "Failed to get response from Pwned Password API", nil, nil)
	}

	// Read the response body line by line using a scanner.
	scanner := bufio.NewScanner(resp.Body)
	results := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")

		// The line should contain a hash suffix and a count.
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
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to read response from Pwned Password API", &err, nil)
	}

	return results, nil
}
