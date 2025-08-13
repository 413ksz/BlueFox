package validation_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
)

// TestDTO must match the structure expected by your validation.ValidateRequestBody function
type TestDTO struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

// errorReader is a mock io.Reader that always returns an error
type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock read error from errorReader")
}

func TestValidateRequestBody(t *testing.T) {
	// Constants for clarity
	const (
		DEFAULT_MAX_MB_BUFFER = 1
		BYTES_PER_MEGABYTE    = 1024 * 1024
	)

	tests := []struct {
		name         string
		requestBody  string
		contentType  string           // To test Content-Type header
		maxMB        int64            // maxMB is in MB as per the function's parameter
		expectedCode models.ErrorCode // Only checking code, not message
		expectedDto  TestDTO
	}{
		// --- Valid JSON Scenarios ---
		{
			name:         "Valid JSON - Basic",
			requestBody:  `{"name": "John Doe", "age": 30}`,
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: "John Doe", Age: 30},
		},
		{
			name:         "Valid JSON - With Optional Field",
			requestBody:  `{"name": "Jane Smith", "age": 25, "email": "jane@example.com"}`,
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: "Jane Smith", Age: 25, Email: "jane@example.com"},
		},
		{
			name:         "Valid JSON - Default Max MB Buffer (0 as input)",
			requestBody:  `{"name": "Short", "age": 10}`,
			contentType:  "application/json",
			maxMB:        0,  // Should use default 1MB
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: "Short", Age: 10},
		},
		{
			name:         "Valid JSON - Negative Max MB Buffer (uses default)",
			requestBody:  `{"name": "Another", "age": 40}`,
			contentType:  "application/json",
			maxMB:        -5, // Should use default 1MB
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: "Another", Age: 40},
		},
		{
			name:         "Valid JSON - With Charset in Content-Type",
			requestBody:  `{"name": "Charset Test", "age": 35}`,
			contentType:  "application/json; charset=utf-8",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: "Charset Test", Age: 35},
		},
		{
			name:         "Valid JSON - Case Insensitive Content-Type",
			requestBody:  `{"name": "Case Test", "age": 45}`,
			contentType:  "APPLICATION/JSON",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: "Case Test", Age: 45},
		},

		// --- Content-Type Validation ---
		{
			name:         "Missing Content-Type Header",
			requestBody:  `{"name": "test", "age": 1}`,
			contentType:  "", // Missing
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: models.ERROR_CODE_BAD_REQUEST,
			expectedDto:  TestDTO{},
		},
		{
			name:         "Incorrect Content-Type Header",
			requestBody:  `{"name": "test", "age": 1}`,
			contentType:  "text/plain", // Incorrect
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: models.ERROR_CODE_BAD_REQUEST,
			expectedDto:  TestDTO{},
		},

		// --- Body Size Limit Validation ---
		{
			name:         "Body size well within limit (1KB JSON)",
			requestBody:  `{"name": "` + strings.Repeat("x", 900) + `", "age": 1}`, // Valid JSON, 900+ bytes
			contentType:  "application/json",
			maxMB:        1,  // 1MB limit. 1KB is well within.
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: strings.Repeat("x", 900), Age: 1},
		},
		{
			name:         "Body size exactly at limit",
			requestBody:  `{"name": "` + strings.Repeat("x", int(DEFAULT_MAX_MB_BUFFER*BYTES_PER_MEGABYTE)-30) + `", "age": 1}`,
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: strings.Repeat("x", int(DEFAULT_MAX_MB_BUFFER*BYTES_PER_MEGABYTE)-30), Age: 1},
		},
		{
			name:         "Body size slightly over default 1MB limit",
			requestBody:  strings.Repeat("a", int(DEFAULT_MAX_MB_BUFFER*BYTES_PER_MEGABYTE)+1),
			contentType:  "application/json",
			maxMB:        0, // Uses default 1MB
			expectedCode: models.ERROR_CODE_BAD_REQUEST,
			expectedDto:  TestDTO{},
		},
		{
			name:         "Body size slightly over explicit 1MB limit",
			requestBody:  strings.Repeat("a", int(1*BYTES_PER_MEGABYTE)+1),
			contentType:  "application/json",
			maxMB:        1,
			expectedCode: models.ERROR_CODE_BAD_REQUEST,
			expectedDto:  TestDTO{},
		},
		{
			name:         "Large payload with default limit (exceeds)",
			requestBody:  `{"data": "` + strings.Repeat("x", DEFAULT_MAX_MB_BUFFER*BYTES_PER_MEGABYTE) + `"}` + "extra", // Valid JSON string content + extra byte to exceed
			contentType:  "application/json",
			maxMB:        0, // Should use default 1MB limit
			expectedCode: models.ERROR_CODE_BAD_REQUEST,
			expectedDto:  TestDTO{},
		},
		{
			name:         "Custom large limit (5MB)",
			requestBody:  `{"name": "` + strings.Repeat("x", 4*BYTES_PER_MEGABYTE) + `", "age": 1}`,
			contentType:  "application/json",
			maxMB:        5,
			expectedCode: "", // No error
			expectedDto:  TestDTO{Name: strings.Repeat("x", 4*BYTES_PER_MEGABYTE), Age: 1},
		},
		{
			name:         "Empty Body",
			requestBody:  "",
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: models.ERROR_CODE_JSON_EMPTY,
			expectedDto:  TestDTO{},
		},

		// --- JSON Decoding Error Scenarios ---
		{
			name:         "Invalid JSON Syntax",
			requestBody:  `{"name": "John Doe", "age": 30,,}`, // Extra comma
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: models.ERROR_CODE_JSON_SYNTAX,
			expectedDto:  TestDTO{},
		},
		{
			name:         "Type Mismatch",
			requestBody:  `{"name": "John Doe", "age": "thirty"}`,
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: models.ERROR_CODE_JSON_TYPE_MISMATCH,
			expectedDto:  TestDTO{},
		},
		{
			name:         "Unknown Field",
			requestBody:  `{"name": "John Doe", "age": 30, "address": "123 Main St"}`,
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: models.ERROR_CODE_JSON_UKNOWN_FIELD,
			expectedDto:  TestDTO{},
		},
		{
			name:         "JSON with Extra Data After Valid Object",
			requestBody:  `{"name": "John", "age": 30}extra`,
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: models.ERROR_CODE_JSON_SYNTAX,
			expectedDto:  TestDTO{},
		},
		{
			name:         "Invalid UTF-8 in JSON",
			requestBody:  string([]byte{0xed, 0xa0, 0x80}), // Incomplete UTF-8 sequence for a high surrogate
			contentType:  "application/json",
			maxMB:        DEFAULT_MAX_MB_BUFFER,
			expectedCode: models.ERROR_CODE_JSON_SYNTAX,
			expectedDto:  TestDTO{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			var dto TestDTO
			err := validation.ValidateRequestBody(&dto, req, tt.maxMB)

			if tt.expectedCode == "" { // Expect no error
				if err != nil {
					t.Errorf("Test %q: Expected no error, but got %s (code: %s, msg: '%s')", tt.name, err, err.Code, err.Message)
				}
				if tt.expectedDto != dto {
					t.Errorf("Test %q: Expected DTO %+v, but got %+v", tt.name, tt.expectedDto, dto)
				}
			} else { // Expect an error
				if err == nil {
					t.Errorf("Test %q: Expected error with code %s, but got nil", tt.name, tt.expectedCode)
				} else {
					if err.Code != tt.expectedCode {
						t.Errorf("Test %q: Expected error code %s, but got %s. Message: '%s'", tt.name, tt.expectedCode, err.Code, err.Message)
					}
				}
			}
		})
	}
}

// Test for specific read errors that aren't about size limits
func TestValidateRequestBody_IOReadError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", errorReader{})
	req.Header.Set("Content-Type", "application/json") // Must set Content-Type

	var dto TestDTO
	err := validation.ValidateRequestBody(&dto, req, 0) // Use default limit

	if err == nil {
		t.Error("Expected an error due to mock read error, but got nil")
	} else if err.Code != models.ERROR_CODE_INTERNAL_SERVER {
		t.Errorf("Expected error code %s for read error, but got %s", models.ERROR_CODE_INTERNAL_SERVER, err.Code)
	}
}

// Test for nil request body
func TestValidateRequestBody_NilBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")

	var dto TestDTO
	err := validation.ValidateRequestBody(&dto, req, 1)

	if err == nil {
		t.Error("Expected an error for nil request body, but got nil")
	} else if err.Code != models.ERROR_CODE_JSON_EMPTY {
		t.Errorf("Expected error code %s for nil body, but got %s", models.ERROR_CODE_JSON_EMPTY, err.Code)
	}
}
