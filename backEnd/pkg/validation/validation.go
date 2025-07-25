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
// Upon encountering an error, it populates a custom error with appropriate error details
// and an HTTP status code, then returns a CustomError.
//
// Parameters:
//   - dto: A pointer to the struct that will receive the unmarshaled JSON data from the request body.
//   - request: The incoming `http.Request` containing the JSON body to be validated and decoded.
//
// Returns:
//   - *models.CustomError: Returns a pointer to a `CustomError` if any validation or decoding
//     error occurs. Returns nil if the request body is successfully read, parsed, and unmarshaled into the 'dto'
//     **without any validation or decoding errors.**
func ValidateRequestBody[T any](dto *T, request *http.Request) *models.CustomError {

	// ----------- Read Request Body -----------
	// Read the entire request body into a byte slice.
	// This is necessary because the request.Body is an io.Reader and can only be read once.
	// If there's an error reading the body (e.g., network issue, client disconnect),
	// create an internal server error.
	body, err := io.ReadAll(request.Body)
	if err != nil {
		apiError := models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "error reading request body", &err, nil)
		return apiError
	}
	// Defer closing the request body. This ensures the body is closed after the function returns.
	defer request.Body.Close()

	// ----------- JSON Decoding -----------
	// Create a new JSON decoder for the read body.
	// Configure the decoder to return an error if the JSON contains fields not present
	// in the `dto` struct.
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	// Attempt to decode the JSON request body into the provided DTO struct.
	// If there's an error, create a custom error and return it.
	// The `dto` parameter is already a pointer, so we pass it directly.
	if err := decoder.Decode(dto); err != nil {

		// Check if the error is a JSON UnmarshalTypeError, which indicates a type mismatch
		// (e.g., expecting an integer but receiving a string).
		if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			apiError := models.NewCustomError(models.ERROR_CODE_JSON_TYPE_MISMATCH,
				fmt.Sprintf("Invalid type for field %s: expected %s, got %s", unmarshalErr.Field, unmarshalErr.Type, unmarshalErr.Value),
				&err, nil)
			return apiError
		}

		// Check if the error is a JSON SyntaxError, which indicates malformed JSON.
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			apiError := models.NewCustomError(models.ERROR_CODE_JSON_SYNTAX,
				fmt.Sprintf("JSON syntax error: %s", syntaxErr.Error()),
				&err, nil)
			return apiError
		}

		// Check if the error message indicates an unknown field.
		// This handles cases where `DisallowUnknownFields()` leads to an error containing this string.
		// For example: "json: unknown field \"field_name\"".
		// Multiple unknown fields are not checked because the standard library's JSON decoder
		// reports only one unknown field at a time in its error string.
		// (e.g., it will not produce an error like "json: unknown field \"field1\", json: unknown field \"field2\"").
		// This is a limitation of the built-in JSON decoder; it can break if the format of the error message changes,
		// making this parsing logic highly fragile.
		if strings.Contains(err.Error(), "json: unknown field") {
			unknownField := ""
			startIndex := strings.Index(err.Error(), "\"")
			endIndex := strings.LastIndex(err.Error(), "\"")
			if startIndex < endIndex {
				unknownField = err.Error()[startIndex+1 : endIndex]
			}

			apiError := models.NewCustomError(models.ERROR_CODE_JSON_UKNOWN_FIELD,
				fmt.Sprintf("Unknown field: %s", unknownField), &err, nil)
			return apiError
		}

		// Check if the error is io.EOF, which indicates an empty request body.
		if err == io.EOF {
			// Create a specific API error for an empty request body.
			apiError := models.NewCustomError(models.ERROR_CODE_JSON_EMPTY, nil, &err, nil)
			return apiError
		}

		// For any other unhandled decoding errors, return a generic bad request error.
		apiError := models.NewCustomError(models.ERROR_CODE_BAD_REQUEST, "Unexpected error decoding request body", &err, nil)
		return apiError
	}

	return nil
}
