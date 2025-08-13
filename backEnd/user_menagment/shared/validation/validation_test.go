package validation_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect" // Added for deep comparison of maps
	"testing"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/shared/validation"
)

// MockClient is a mock HTTP client that allows us to control the response for testing.
type MockClient struct {
	// DoFunc is a function that a test can provide to mock the http.Client's Do method.
	DoFunc func(req *http.Request) (*http.Response, error)
}

// httpError is a custom error type that implements the Timeout method
// to simulate timeout errors for testing.
type httpError struct {
	timeout bool
}

func (e *httpError) Error() string {
	return "mock timeout error"
}

func (e *httpError) Timeout() bool {
	return e.timeout
}

func (e *httpError) Temporary() bool {
	return false
}

// Do implements the Doer interface for the MockClient.
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	// Return a default successful response if no mock function is provided.
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}, nil
}

// TestGetPwnedPasswordHashes tests the GetPwnedPasswordHashes function with mocked API responses.
func TestGetPwnedPasswordHashes(t *testing.T) {
	// Define a slice of test cases.
	tests := []struct {
		name         string
		identifier   string
		mockResponse *http.Response
		mockError    error
		expectedMap  map[string]int
		expectedCode models.ErrorCode
		expectErr    bool
	}{
		{
			name:         "Invalid identifier length (too short)",
			identifier:   "1234",
			expectErr:    true,
			expectedCode: models.ERROR_CODE_BAD_REQUEST,
		},
		{
			name:         "Invalid identifier length (too long)",
			identifier:   "123456",
			expectErr:    true,
			expectedCode: models.ERROR_CODE_BAD_REQUEST,
		},
		{
			name:       "Valid identifier (successful response)",
			identifier: "CB66F",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("F08FF:1234\nF94DF:5678")),
			},
			expectedMap: map[string]int{"F08FF": 1234, "F94DF": 5678},
			expectErr:   false,
		},
		{
			name:       "No match found (404 Not Found)",
			identifier: "E8640",
			mockResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			expectedMap: map[string]int{},
			expectErr:   false,
		},
		{
			name:       "Unexpected status code (500 Internal Server Error)",
			identifier: "CB66F",
			mockResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			expectErr:    true,
			expectedCode: models.ERROR_CODE_EXTERNAL_SERVER,
		},
		{
			name:         "Network error (request fails)",
			identifier:   "CB66F",
			mockError:    fmt.Errorf("network error"),
			expectErr:    true,
			expectedCode: models.ERROR_CODE_INTERNAL_SERVER,
		},
		{
			name:       "Empty API response (valid but empty)",
			identifier: "CB66F",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			expectedMap: map[string]int{},
			expectErr:   false,
		},
		{
			name:       "Unexpected status code (429 Too Many Requests)",
			identifier: "CB66F",
			mockResponse: &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			expectErr:    true,
			expectedCode: models.ERROR_CODE_EXTERNAL_SERVER,
		},
		{
			name:       "Case sensitivity in identifier",
			identifier: "cb66f", // Lowercase
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("F08FF:1234")),
			},
			expectedMap: map[string]int{"F08FF": 1234},
			expectErr:   false,
		},
		{
			name:       "Malformed API response (non-integer count)",
			identifier: "CB66F",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("F08FF:1234\nF94DF:ABC")),
			},
			expectedMap: map[string]int{"F08FF": 1234},
			expectErr:   false,
		},
		{
			name:       "Malformed API response (invalid line format)",
			identifier: "CB66F",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("F08FF:1234\nINVALIDLINE\nF94DF:5678")),
			},
			expectedMap: map[string]int{"F08FF": 1234, "F94DF": 5678},
			expectErr:   false,
		},
		{
			name:       "Malformed API response (line with extra parts)",
			identifier: "CB66F",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("F08FF:1234:extra\nF94DF:5678")),
			},
			expectedMap: map[string]int{"F94DF": 5678},
			expectErr:   false,
		},
		{
			name:       "Timeout error",
			identifier: "CB66F",
			mockError: &httpError{
				timeout: true,
			},
			expectErr:    true,
			expectedCode: models.ERROR_CODE_EXTERNAL_SERVER,
		},
	}

	// Iterate through the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock client for this specific test case.
			mockClient := &MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockResponse, nil
				},
			}

			// Instantiate the PwnedPassword service with the mock client.
			pwnedPassword := validation.NewPwnedPassword(mockClient)

			// Call the function under test.
			result, customErr := pwnedPassword.GetPwnedPasswordHashes(tt.identifier)

			// Check for expected error.
			if tt.expectErr {
				if customErr == nil {
					t.Errorf("Expected an error, but got none")
					return
				}
				if customErr.Code != tt.expectedCode {
					t.Errorf("Expected error code %v, got %v", tt.expectedCode, customErr.Code)
				}
			} else {
				if customErr != nil {
					t.Errorf("Expected no error, but got: %v", customErr)
					return
				}
				// Corrected: Use reflect.DeepEqual for robust map comparison.
				if !reflect.DeepEqual(result, tt.expectedMap) {
					t.Errorf("Expected map %v, but got %v", tt.expectedMap, result)
				}
			}
		})
	}
}
