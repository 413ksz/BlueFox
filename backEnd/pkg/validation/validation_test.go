package validation_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
)

// TestValidateRequestBodyErrorHelper is the table-driven test function for the ValidateRequestBodyErrorHelper function.
func TestValidateRequestBodyErrorHelper(t *testing.T) {
	// Define the test cases in a slice of structs.
	testCases := []struct {
		name                string
		inputErr            error
		expectedErrorCode   models.ErrorCode
		expectedDetailsPart string
		expectedMessage     string
	}{
		{
			name: "UnmarshalTypeError_ValidInput",
			inputErr: &json.UnmarshalTypeError{
				Value:  "string",
				Type:   reflect.TypeOf(0),
				Offset: 1,
				Struct: "ExampleStruct",
				Field:  "age",
			},
			expectedErrorCode:   models.ERROR_CODE_JSON_TYPE_MISMATCH,
			expectedDetailsPart: "Invalid type for field 'age'",
			expectedMessage:     "The request body contains a field with an unexpected type.",
		},
		{
			name: "SyntaxError_ValidInput",
			inputErr: &json.SyntaxError{
				Offset: 1,
			},
			expectedErrorCode:   models.ERROR_CODE_JSON_SYNTAX,
			expectedDetailsPart: "JSON syntax error",
			expectedMessage:     "The request body contains malformed JSON or invalid JSON syntax.",
		},
		{
			name:                "UnknownField_WithQuotes",
			inputErr:            errors.New("json: unknown field \"name\""),
			expectedErrorCode:   models.ERROR_CODE_JSON_UKNOWN_FIELD,
			expectedDetailsPart: "Unknown field in request body: 'name'",
			expectedMessage:     "The request body contains an unknown field.",
		},
		{
			name:                "UnknownField_WithoutQuotes",
			inputErr:            errors.New("json: unknown field"),
			expectedErrorCode:   models.ERROR_CODE_JSON_UKNOWN_FIELD,
			expectedDetailsPart: "Unknown field in request body: ''",
			expectedMessage:     "The request body contains an unknown field.",
		},
		{
			name:                "EOF_Error",
			inputErr:            io.EOF,
			expectedErrorCode:   models.ERROR_CODE_JSON_EMPTY,
			expectedDetailsPart: "Request body is empty or malformed",
			expectedMessage:     "The request body is empty.",
		},
		{
			name:                "Generic_UnexpectedError",
			inputErr:            errors.New("some other random error"),
			expectedErrorCode:   models.ERROR_CODE_BAD_REQUEST,
			expectedDetailsPart: "Unexpected error decoding request body",
			expectedMessage:     "The request was invalid or malformed.",
		},
	}

	// Iterate over the test cases.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function we are testing.
			result := validation.ValidateRequestBodyErrorHelper(tc.inputErr)

			// Assertions
			if result == nil {
				t.Fatalf("Expected a CustomError, but got nil")
			}

			if result.Code != tc.expectedErrorCode {
				t.Errorf("Expected error code '%s', but got '%s'", tc.expectedErrorCode, result.Code)
			}

			if result.Message != tc.expectedMessage {
				t.Errorf("Expected message '%s', but got '%s'", tc.expectedMessage, result.Message)
			}

			if details, ok := result.Details.(string); ok {
				if !strings.Contains(details, tc.expectedDetailsPart) {
					t.Errorf("Expected details to contain '%s', but got '%s'", tc.expectedDetailsPart, details)
				}
			} else {
				t.Errorf("Expected details to be a string, but it was not.")
			}
		})
	}
}

// FuzzTargetStruct is a sample struct used for fuzzing.
type FuzzTargetStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// FuzzValidateRequestBodyErrorHelper is the fuzzer function for ValidateRequestBodyErrorHelper.
func FuzzValidateRequestBodyErrorHelper(f *testing.F) {
	f.Add([]byte(`{"name": "John", "age": 30}`))                     // Valid JSON
	f.Add([]byte(`{"name": "Jane", "age": "twenty-five"}`))          // JSON type mismatch
	f.Add([]byte(`{"name": "Jack", "age": 20, "extraField": true}`)) // Unknown field
	f.Add([]byte(`{"name": "Mike"`))                                 // JSON syntax error
	f.Add([]byte(``))                                                // Empty input (EOF)
	f.Add([]byte(`bad data`))                                        // Generic invalid JSON

	f.Fuzz(func(t *testing.T, data []byte) {
		r := bytes.NewReader(data)
		decoder := json.NewDecoder(r)
		decoder.DisallowUnknownFields()

		var target FuzzTargetStruct
		err := decoder.Decode(&target)

		if err != nil && err != io.EOF {
			result := validation.ValidateRequestBodyErrorHelper(err)
			if result == nil {
				t.Fatalf("Expected a CustomError for error: %v, but got nil", err)
			}
		} else if err == io.EOF {
			result := validation.ValidateRequestBodyErrorHelper(err)
			if result == nil {
				t.Fatalf("Expected a CustomError for io.EOF, but got nil")
			}
		}
	})
}

// TestConvertMbToBytes tests the ConvertMbToBytes function with table-driven tests.
func TestConvertMbToBytes(t *testing.T) {
	testCases := []struct {
		name            string
		input           int64
		expected        int64
		expectedErrCode models.ErrorCode
	}{
		{
			name:            "Zero megabytes",
			input:           0,
			expected:        0,
			expectedErrCode: "", // No error expected
		},
		{
			name:            "One megabyte",
			input:           1,
			expected:        1048576,
			expectedErrCode: "",
		},
		{
			name:            "Large number within bounds",
			input:           math.MaxInt64 / (1024 * 1024),
			expected:        (math.MaxInt64 / (1024 * 1024)) * 1024 * 1024,
			expectedErrCode: "",
		},
		{
			name:            "Negative input",
			input:           -10,
			expected:        0,
			expectedErrCode: models.ERROR_CODE_INTERNAL_SERVER,
		},
		{
			name:            "Overflow input",
			input:           (math.MaxInt64 / (1024 * 1024)) + 1,
			expected:        0,
			expectedErrCode: models.ERROR_CODE_INTERNAL_SERVER,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := validation.ConvertMbToBytes(tc.input)

			if tc.expectedErrCode != "" {
				if err == nil {
					t.Errorf("Expected an error with code %s but got nil", tc.expectedErrCode)
				} else if err.Code != tc.expectedErrCode {
					t.Errorf("Expected error code %s, but got %s", tc.expectedErrCode, err.Code)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got %v", err)
				}

				if got != tc.expected {
					t.Errorf("For input %d, expected %d, but got %d", tc.input, tc.expected, got)
				}
			}
		})
	}
}

// FuzzConvertMbToBytes tests the ConvertMbToBytes function with a fuzzer.
func FuzzConvertMbToBytes(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(1))
	f.Add(int64(1024))
	f.Add(int64(math.MaxInt64 / (1024 * 1024)))
	f.Add(int64(-1)) // Negative input test case
	f.Fuzz(func(t *testing.T, megaByte int64) {
		bytes, err := validation.ConvertMbToBytes(megaByte)
		if megaByte < 0 {
			if err == nil {
				t.Fatalf("Expected an error for negative input %d, but got nil", megaByte)
			}
			return
		}
		if megaByte > math.MaxInt64/(1024*1024) {
			if err == nil {
				t.Fatalf("Expected an error for overflow input %d, but got nil", megaByte)
			}
			return
		}
		if err != nil {
			t.Fatalf("Did not expect an error for input %d, but got %v", megaByte, err)
		}
		expectedBytes := megaByte * 1024 * 1024
		if bytes != expectedBytes {
			t.Errorf("For input %d, expected %d bytes, but got %d", megaByte, expectedBytes, bytes)
		}
	})
}

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
		// --- Convert Body Size ---
		{
			name:         "Large Max MB Buffer Overflow",
			requestBody:  `{"name": "overflow", "age": 10}`,
			contentType:  "application/json",
			maxMB:        math.MaxInt64, // A value that will cause an overflow in ConvertMbToBytes
			expectedCode: models.ERROR_CODE_INTERNAL_SERVER,
			expectedDto:  TestDTO{},
		},
		{
			name:         "Negative Max MB Buffer (uses default)",
			requestBody:  `{"name": "overflow", "age": 10}`,
			contentType:  "application/json",
			maxMB:        -1, // Should use the default 1MB
			expectedCode: "", // Expect no error
			expectedDto:  TestDTO{Name: "overflow", Age: 10},
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

// FuzzValidateRequestBody is a fuzzer for the ValidateRequestBody function.
func FuzzValidateRequestBody(f *testing.F) {
	f.Add([]byte(`{"name": "John", "age": 30}`), int64(1))                                         // Valid JSON, 1MB limit
	f.Add([]byte(`{"name": "Jane", "age": "twenty-five"}`), int64(1))                              // JSON type mismatch
	f.Add([]byte(`{"name": "Jack", "age": 20, "extraField": true}`), int64(1))                     // Unknown field
	f.Add([]byte(`{"name": "Mike"`), int64(1))                                                     // JSON syntax error
	f.Add([]byte(`{"name": "Test", "age": 1, "email": "test@example.com"}`), int64(0))             // Valid JSON, test default 1MB limit
	f.Add([]byte(``), int64(1))                                                                    // Empty input (EOF)
	f.Add([]byte(`{"data":"`+strings.Repeat("a", 2*1024*1024)+`"}`), int64(1))                     // Oversized valid JSON body (2MB)
	f.Add([]byte(`{"name": "Test", "age": 1, "email": "test@example.com"}trailingdata`), int64(1)) // Valid JSON with trailing data

	f.Fuzz(func(t *testing.T, data []byte, maxMB int64) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		var dto TestDTO
		err := validation.ValidateRequestBody(&dto, req, maxMB)
		if err != nil {
			_ = err
		}
	})
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
