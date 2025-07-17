package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   any    `json:"value,omitempty"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

// ValidationError is a struct that encapsulates a single validation error.
// It is used to provide structured validation failure information.
type ValidationError struct {
	Key     string `json:"key"`
	Value   any    `json:"-"` // Value is typically omitted from JSON output for security reasons
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// NewValidationError creates and returns a new ValidationError pointer
// with the provided key, value, failed validation rule, and a custom message.
//
// Parameters:
// - key: The key or name of the field that failed validation.
// - value: The actual value of the field that failed validation (excluded from JSON).
// - failedValidation: The specific validation rule (e.g., "required", "min") that failed.
// - message: A custom, human-readable error message for the validation error.
//
// Returns:
// - *ValidationError: A pointer to the newly created ValidationError object.
func NewValidationError(key string, value any, failedValidation string, message string) *ValidationError {
	return &ValidationError{Key: key, Value: value, Rule: failedValidation, Message: message}
}

// ErrorCode is a type that represents a standardized custom error code
// for API responses.
type ErrorCode string

// Error codes constants define common API error scenarios.
const (
	ERROR_CODE_BAD_REQUEST          ErrorCode = "BAD_REQUEST"
	ERROR_CODE_JSON_SYNTAX          ErrorCode = "JSON_SYNTAX_ERROR"
	ERROR_CODE_JSON_TYPE_MISMATCH   ErrorCode = "JSON_MISMATCH_ERROR"
	ERROR_CODE_JSON_UKNOWN_FIELD    ErrorCode = "JSON_UNKNOWN_FIELD_ERROR"
	ERROR_CODE_UNPROCESSABLE_ENTITY ErrorCode = "UNPROCESSABLE_ENTITY"
	ERROR_CODE_INTERNAL_SERVER      ErrorCode = "INTERNAL_SERVER_ERROR"
)

// CustomError is a struct that encapsulates detailed custom API error information.
// It is designed to provide comprehensive error context to both clients and logs.
type CustomError struct {
	Code       ErrorCode     `json:"code"`
	Message    string        `json:"message"`
	Details    any           `json:"details,omitempty"`
	HttpCode   int           `json:"-"`
	LogLevel   zerolog.Level `json:"-"`
	Err        *error        `json:"-"`
	StackTrace *string       `json:"-"`
}

// NewCustomError creates and returns a new CustomError pointer.
// It initializes the error with a specific code, optional underlying error,
// custom details, and an optional stack trace. It also sets default HTTP
// status code and log level based on the ErrorCode.
//
// Parameters:
//   - code: The custom ErrorCode identifying the type of error.
//   - originalError: An optional pointer to the underlying Go error that caused this CustomError.
//     Pass nil if no original error is associated.
//   - details: Optional structured details about the error, often a slice of ValidationErrorDetail
//     or a simple string. Pass nil if no specific details are needed.
//   - stackTrace: An optional pointer to a string containing the stack trace where the error occurred.
//     Pass nil if no stack trace is available or desired.
//
// Returns:
// - *CustomError: A pointer to the newly created CustomError object.
func NewCustomError(code ErrorCode, originalError *error, details any, stackTrace *string) *CustomError {
	var customError CustomError

	customError.Code = code
	customError.Err = originalError
	customError.Details = details
	customError.StackTrace = stackTrace

	switch code {
	case ERROR_CODE_BAD_REQUEST:
		customError.HttpCode = http.StatusBadRequest
		customError.LogLevel = zerolog.WarnLevel
		customError.Message = "The request was invalid or malformed."
	case ERROR_CODE_UNPROCESSABLE_ENTITY:
		customError.HttpCode = http.StatusUnprocessableEntity
		customError.LogLevel = zerolog.WarnLevel
		customError.Message = "One or more input values are invalid"
	case ERROR_CODE_JSON_SYNTAX:
		customError.HttpCode = http.StatusBadRequest
		customError.LogLevel = zerolog.WarnLevel
		customError.Message = "The request body contains malformed JSON or invalid JSON syntax."
	case ERROR_CODE_JSON_TYPE_MISMATCH:
		customError.HttpCode = http.StatusBadRequest
		customError.LogLevel = zerolog.WarnLevel
		customError.Message = "The request body contains a field with an unexpected type."
	case ERROR_CODE_JSON_UKNOWN_FIELD:
		customError.HttpCode = http.StatusBadRequest
		customError.LogLevel = zerolog.WarnLevel
		customError.Message = "The request body contains an unknown field."
	case ERROR_CODE_INTERNAL_SERVER:
		customError.HttpCode = http.StatusInternalServerError
		customError.LogLevel = zerolog.ErrorLevel
		customError.Message = "An internal server error occurred."
	default:
		customError.HttpCode = http.StatusInternalServerError
		customError.LogLevel = zerolog.ErrorLevel
		customError.Message = fmt.Sprintf("An unexpected error occurred. Code: %s", code)
	}

	return &customError
}

// Pagination is a struct that encapsulates pagination information for API responses.
type Pagination struct {
	TotalItems   int     `json:"totalItems"`
	ItemsPerPage *int    `json:"itemsPerPage,omitempty"`
	PageIndex    *int    `json:"pageIndex,omitempty"`
	NextLink     *string `json:"nextLink,omitempty"`
	PreviousLink *string `json:"previousLink,omitempty"`
}

// ResponseData is a generic struct that encapsulates the actual data payload
// of an API response, along with optional metadata like deletion status,
// update timestamps, and pagination.
type ResponseData[TItemtype any] struct {
	Deleted    *bool       `json:"deleted,omitempty"`
	Updated    *time.Time  `json:"updated,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Items      []TItemtype `json:"items"`
}

// ApiResponse is a struct that encapsulates a response to the client.
type ApiResponse[TItemtype any] struct {
	Params     map[string]interface{}   `json:"params,omitempty"`
	Data       *ResponseData[TItemtype] `json:"data,omitempty"`
	Error      *CustomError             `json:"error,omitempty"`
	Message    string                   `json:"message,omitempty"`
	StatusCode int                      `json:"-"`
}

// NewApiResponse creates and returns a new ApiResponse pointer with
// initial parameters, status code, and a message. It defaults to HTTP 200 OK
// and its corresponding status text if statusCode or message are zero/empty.
//
// Parameters:
//   - params: Optional map of request parameters to include in the response.
//   - statusCode: The HTTP status code for the response (e.g., http.StatusOK, http.StatusBadRequest).
//     Defaults to http.StatusOK if 0.
//   - message: A descriptive message for the API response. Defaults to HTTP status text if empty.
//
// Returns:
// - *ApiResponse: A pointer to the newly created ApiResponse object.
func NewApiResponse[TItemtype any](params map[string]interface{}, statusCode int, message string) *ApiResponse[TItemtype] {
	if statusCode == 0 {
		statusCode = http.StatusOK

	}
	if message == "" {
		message = http.StatusText(statusCode)

	}

	return &ApiResponse[TItemtype]{
		Params:     params,
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
		Error:      nil,
	}
}

// WithData modifies an existing ApiResponse to include a data payload.
// It sets the response message, data, and HTTP status code, and clears any existing error.
//
// Parameters:
// - Message: The message for the successful response.
// - Data: The ResponseData payload to include.
// - statusCode: The HTTP status code for the successful response (e.g., http.StatusOK).
func (apiResponse *ApiResponse[TItemtype]) WithData(Message string, Data *ResponseData[TItemtype], statusCode int) {
	apiResponse.Message = Message
	apiResponse.Data = Data
	apiResponse.StatusCode = statusCode
	apiResponse.Error = nil
}

// WithError modifies an existing ApiResponse to include an error payload.
// It sets the response message, CustomError, and HTTP status code, and clears any existing data.
//
// Parameters:
// - Message: The message for the error response.
// - Error: The CustomError object containing detailed error information.
// - statusCode: The HTTP status code for the error response (e.g., http.StatusBadRequest).
func (apiResponse *ApiResponse[TItemtype]) WithError(Message string, Error *CustomError, statusCode int) {
	apiResponse.Message = Message
	apiResponse.Error = Error
	apiResponse.StatusCode = statusCode
	apiResponse.Data = nil
}

// SendApiResponse sends the given ApiResponse as a JSON response to the client.
// It sets the Content-Type header, writes the HTTP status code, and encodes the
// apiResponse to the response body. Errors during encoding are logged.
//
// Parameters:
// - r: The HTTP request.
// - w: The HTTP response writer.
// - apiResponse: The ApiResponse to be sent.
// - requestId: The unique request ID for logging purposes (e.g., for tracing).
// - componentName: The name of the component or service sending the response for logging.
func SendApiResponse[T any](r *http.Request, w http.ResponseWriter, apiResponse *ApiResponse[T], requestId string, componentName string) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the HTTP status code
	w.WriteHeader(apiResponse.StatusCode)

	// Encode the apiResponse struct to JSON and write it to the response body
	err := json.NewEncoder(w).Encode(apiResponse)
	if err != nil {
		// Log an error if encoding fails, but don't attempt to write to w again
		// as headers have already been sent.
		log.Error().
			Str("request_id", requestId).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("component", componentName).
			Str("status", "failed").
			Str("event", "api_response_encoding_failed").
			Msg("Error occurred encoding API response")
	}
}
