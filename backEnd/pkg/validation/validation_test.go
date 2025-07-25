package validation_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
)

type TestDTO struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

func TestValidateRequestBody(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedError  *models.CustomError
		expectedResult TestDTO
	}{
		{
			name:          "Valid JSON",
			requestBody:   `{"name": "John Doe", "age": 30}`,
			expectedError: nil,
			expectedResult: TestDTO{
				Name: "John Doe",
				Age:  30,
			},
		},
		{
			name:           "Empty Body",
			requestBody:    "",
			expectedError:  models.NewCustomError(models.ERROR_CODE_JSON_EMPTY, nil, nil, nil),
			expectedResult: TestDTO{},
		},
		{
			name:           "Invalid JSON Syntax",
			requestBody:    `{"name": "John Doe", "age": 30,}`,
			expectedError:  models.NewCustomError(models.ERROR_CODE_JSON_SYNTAX, "JSON syntax error: invalid character '}' looking for beginning of object key string", nil, nil),
			expectedResult: TestDTO{},
		},
		{
			name:           "Type Mismatch",
			requestBody:    `{"name": "John Doe", "age": "thirty"}`,
			expectedError:  models.NewCustomError(models.ERROR_CODE_JSON_TYPE_MISMATCH, "Invalid type for field age: expected int, got string", nil, nil),
			expectedResult: TestDTO{},
		},
		{
			name:           "Unknown Field",
			requestBody:    `{"name": "John Doe", "age": 30, "address": "123 Main St"}`,
			expectedError:  models.NewCustomError(models.ERROR_CODE_JSON_UKNOWN_FIELD, "Unknown field: address", nil, nil),
			expectedResult: TestDTO{},
		},
		{
			name:          "Valid JSON with optional field",
			requestBody:   `{"name": "John Doe", "age": 30, "email": "john@example.com"}`,
			expectedError: nil,
			expectedResult: TestDTO{
				Name:  "John Doe",
				Age:   30,
				Email: "john@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))

			var dto TestDTO
			err := validation.ValidateRequestBody(&dto, req)

			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Code != tt.expectedError.Code {
					t.Errorf("Expected error code %s, got %s", tt.expectedError.Code, err.Code)
				}
			}

			if tt.expectedError == nil && dto != tt.expectedResult {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, dto)
			}
		})
	}
}

func TestValidateRequestBody_ReadError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", errorReader{})

	var dto TestDTO
	err := validation.ValidateRequestBody(&dto, req)

	if err == nil {
		t.Error("Expected error for failed body read, got nil")
	} else if err.Code != models.ERROR_CODE_INTERNAL_SERVER {
		t.Errorf("Expected error code %s, got %s", models.ERROR_CODE_INTERNAL_SERVER, err.Code)
	}
}

func TestValidateRequestBody_InvalidUTF8(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte{0xff, 0xfe, 0xfd}))

	var dto TestDTO
	err := validation.ValidateRequestBody(&dto, req)

	if err == nil {
		t.Error("Expected error for invalid UTF-8, got nil")
	} else if err.Code != models.ERROR_CODE_JSON_SYNTAX {
		t.Errorf("Expected error code %s, got %s", models.ERROR_CODE_JSON_SYNTAX, err.Code)
	}
}

type badType struct{}

func (b *badType) UnmarshalJSON([]byte) error {
	return errors.New("custom unmarshal error")
}

func TestValidateRequestBody_GenericError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"field": "value"}`))

	var dto struct {
		Field badType `json:"field"`
	}
	err := validation.ValidateRequestBody(&dto, req)

	if err == nil {
		t.Error("Expected error for custom unmarshal error, got nil")
	} else if err.Code != models.ERROR_CODE_BAD_REQUEST {
		t.Errorf("Expected error code %s, got %s", models.ERROR_CODE_BAD_REQUEST, err.Code)
	}
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock read error")
}
