package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
)

// ValidateRequestBody validates the request body of an HTTP request and unmarshals it into a struct (dto).
// It handles various JSON parsing errors such as syntax errors, type mismatches, and unknown fields.
// The function also enforces a maximum request body size, defaulting to 1 MB if no size is specified
//
// TODO: This function reads the entire request body into memory, which is a performance concern
// for very large payloads. A more memory-efficient streaming approach could be used, but this
// would require a change to the error handling logic. For now, this implementation is sufficient.
//
// Parameters:
//   - dto: A pointer to the struct that will receive the unmarshaled JSON data.
//   - request: The incoming `http.Request` containing the JSON body.
//   - maxMBBuffer: The maximum allowed size of the request body in megabytes (MB).
//     If 0 or less, a default of 1 MB is used.
//
// Returns:
//   - *models.CustomError: A custom error if any JSON decoding fails.
//     Returns nil on success.
func ValidateRequestBody[T any](dto *T, request *http.Request, maxMBSize int64) *models.CustomError {
	// Ensure the request body is closed after the function returns.
	defer request.Body.Close()

	// ----------- Configure Max Request Body Size -----------
	// Set a default maximum request body size of 1 MB if maxMBBuffer is 0 or less.
	// Otherwise, use the provided value and convert it to bytes.
	const defaultMaxMbSize = 1

	var maxByteSize int64

	// Ensure maxMBBuffer is not less than 0 MB.
	if maxMBSize <= 0 {
		maxByteSize = ConvertMbToBytes(defaultMaxMbSize) // 1 MB default
	} else {
		maxByteSize = ConvertMbToBytes(maxMBSize) // convert MB to bytes
	}

	// ----------- Read Request Body -----------
	// Read the entire request body into a byte slice.
	// This is necessary because the request.Body is an io.Reader and can only be read once.
	// If there's an error reading the body (e.g., network issue, client disconnect),
	// maximum request body size exceeded error is returned.

	// Use io.LimitReader to protect against excessively large requests.
	// The +1 byte allows us to detect if the body size exceeds the limit.
	body, err := io.ReadAll(io.LimitReader(request.Body, maxByteSize+1)) // +1 byte to detect if the body size exceeds the limit
	// Check if there was an error reading the request body.
	if err != nil {
		apiError := models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Error reading request body", err, nil)
		return apiError
	}
	// Check if the request body size exceeds the maximum allowed size.
	if len(body) >= int(maxByteSize)+1 {
		apiError := models.NewCustomError(models.ERROR_CODE_BAD_REQUEST, "Maximum request body size exceeded", nil, nil)
		return apiError
	}

	// ----------- JSON Decoding -----------
	// Check if the Content-Type header is set to "application/json".
	// Create a new JSON decoder for the read body.
	// Configure the decoder to return an error if the JSON contains fields not present
	// in the `dto` struct.

	// Check if the Content-Type header is set to "application/json".
	if !strings.HasPrefix(strings.ToLower(request.Header.Get("Content-Type")), "application/json") {
		apiError := models.NewCustomError(models.ERROR_CODE_BAD_REQUEST, fmt.Sprintf("Invalid Content-Type header: %s", request.Header.Get("Content-Type")), nil, nil)
		return apiError
	}

	// Create a new JSON decoder for the read body with disallowing unknown fields.
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	// Attempt to decode the JSON request body into the provided DTO struct.
	// If there's an error, create a custom error and return it.
	// The `dto` parameter is already a pointer, so we pass it directly.
	if err := decoder.Decode(dto); err != nil {
		apiError := ValidateRequestBodyErrorHelper(err)
		return apiError
	}

	// If the JSON decoding was successful but there was trailing data in the request body, return a bad request error.
	if decoder.More() {
		apiError := models.NewCustomError(models.ERROR_CODE_JSON_SYNTAX, "Trailing data in request body", nil, nil)
		return apiError
	}

	return nil
}

// ValidateRequestBodyErrorHelper is a helper function to create a CustomError based on the type of error encountered during JSON decoding.
// It takes an error parameter and returns a pointer to a CustomError.
// Parameters:
//   - err: The error encountered during JSON decoding.
//
// Returns:
//   - *models.CustomError: A pointer to a CustomError containing details about the error.
func ValidateRequestBodyErrorHelper(err error) *models.CustomError {
	// Giving a more specific error message for each type of error,
	// which is based on the error type switch
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		// Handles type mismatches during JSON unmarshaling (e.g., expecting int, got string).
		apiError := models.NewCustomError(models.ERROR_CODE_JSON_TYPE_MISMATCH,
			fmt.Sprintf("Invalid type for field '%s': expected %s, got %s", e.Field, e.Type, e.Value),
			err, nil)
		return apiError
	case *json.SyntaxError:
		// Handles malformed JSON syntax errors (e.g., unexpected end of JSON input).
		apiError := models.NewCustomError(models.ERROR_CODE_JSON_SYNTAX,
			fmt.Sprintf("JSON syntax error: %s", e.Error()),
			err, nil)
		return apiError

	default:
		// --- Fragile Unknown Field Detection ---
		// This block attempts to detect "unknown field" errors by string matching.
		// WARNING: This approach is fragile as it relies on the exact error message
		// format from the standard library's JSON decoder. It may break if the
		// format changes in future Go versions or if a different JSON library is used.
		// There is no specific error type for unknown fields in the standard 'json' package.
		if strings.Contains(err.Error(), "json: unknown field") {
			unknownField := ""
			// Attempt to extract the field name enclosed in quotes.
			startIndex := strings.Index(err.Error(), "\"")
			endIndex := strings.LastIndex(err.Error(), "\"")
			if startIndex < endIndex && startIndex != -1 { // Ensure both quotes are found not just the endIndex and the startIndex is still -1
				unknownField = err.Error()[startIndex+1 : endIndex]
			}

			apiError := models.NewCustomError(models.ERROR_CODE_JSON_UKNOWN_FIELD,
				fmt.Sprintf("Unknown field in request body: '%s'", unknownField), err, nil)
			return apiError
		}

		// Check if the error is io.EOF, which indicates an empty request body.
		if err == io.EOF {
			apiError := models.NewCustomError(models.ERROR_CODE_JSON_EMPTY, "Request body is empty or malformed", err, nil)
			return apiError
		}

		// For any other unhandled decoding errors, return a generic bad request error.
		apiError := models.NewCustomError(models.ERROR_CODE_BAD_REQUEST, "Unexpected error decoding request body", err, nil)
		return apiError
	}
}

// ConvertMbToBytes converts megabytes to bytes
// parameters:
//   - mb: megabytes
//
// returns:
//   - int64: bytes
func ConvertMbToBytes(mb int64) int64 {
	return mb * 1024 * 1024
}
