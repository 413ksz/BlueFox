package validation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// Validator is a struct that encapsulates a validator instance.
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new instance of the Validator struct with required struct validation enabled.
func NewValidator() *Validator {
	return &Validator{validator: validator.New(validator.WithRequiredStructEnabled())}
}

// RegisterCustomValidation registers a custom validation function for the validator instance.
// Parameters:
// - tag: The tag to associate with the custom validation function.
// - functionToRegister: The custom validation function to register.
// Returns: An error if the registration fails, nil otherwise.
func (v *Validator) RegisterCustomValidation(tag string, functionToRegister validator.Func) error {

	if tag == "" {
		err := errors.New("validation tag cannot be empty")
		log.Error().Err(err).
			Str("component", "validation").
			Str("function", "RegisterCustomValidation").
			Str("event", "custom_validation_failed").
			Msg("Attempted to register custom validation with an empty tag")
		return err
	}

	if functionToRegister == nil {
		err := errors.New("validation function cannot be nil")
		log.Error().Err(err).
			Str("component", "validation").
			Str("function", "RegisterCustomValidation").
			Str("event", "custom_validation_failed").
			Str("tag", tag).
			Msg("Attempted to register a nil custom validation function")
		return err
	}
	err := v.validator.RegisterValidation(tag, functionToRegister)
	if err != nil {
		log.Error().Err(err).
			Str("component", "validation").
			Str("function", "RegisterCustomValidation").
			Str("event", "custom_validation_failed").
			Str("tag", tag).
			Msg("Failed to register custom validation")
		return err
	}
	log.Info().
		Str("component", "validation").
		Str("function", "RegisterCustomValidation").
		Str("event", "custom_validation_success").
		Str("tag", tag).
		Msg("Custom validation registered successfully")
	return nil
}

// GetValidationErrorMessage returns a human-readable error message based on the tag and field of a
// validator.FieldError. It is used to provide user-friendly error messages for validation errors.
// Parameters:
// - fieldError: The validator.FieldError representing the validation error.
// Returns:
// - string: A human-readable error message.
func GetValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required.", fieldError.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long.", fieldError.Field(), fieldError.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long.", fieldError.Field(), fieldError.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address.", fieldError.Field())
	case "username":
		return fmt.Sprintf("%s must be a valid username.", fieldError.Field())
	case "password":
		return fmt.Sprintf("%s must be a valid password.", fieldError.Field())
	case "dateofbirth":
		return fmt.Sprintf("%s must be a valid date of birth.", fieldError.Field())
	case "name":
		return fmt.Sprintf("%s must be a valid name.", fieldError.Field())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s.", fieldError.Field(), fieldError.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s.", fieldError.Field(), fieldError.Param())
	default:
		return fmt.Sprintf("Validation failed for %s with tag %s.", fieldError.Field(), fieldError.Tag())
	}
}

// ValidateStruct validates a struct based on the registered validations.
// It utilizes the underlying validator instance to ensure that the struct
// adheres to all specified validation rules.
// Parameters:
// - obj: The struct to validate.
// Returns:
// - error: An error is returned if validation fails, otherwise nil.
func (v *Validator) ValidateStruct(obj interface{}) error {
	err := v.validator.Struct(obj)
	if err != nil {
		return err
	}
	return nil
}

// ValidateRequestBody validates the request body of an HTTP request and unmarshals it into a struct (dto).
// It handles various JSON parsing errors such as syntax errors, type mismatches, and unknown fields.
// Upon encountering an error, it populates a custom error with appropriate error details
// and an HTTP status code, then returns a CustomError.
//
//  1. **Go Language Limitation:** In Go, methods defined on a *non-generic* struct (like 'Validator' in this case)
//     cannot have their own type parameters. If 'Validator' were 'Validator[T any]', then its methods could use 'T'.
//     However, since 'Validator' itself doesn't need to be generic, the method must accept a concrete type or an interface.
//     Using [T any](dto[U]) would result in a compile-time error.
//
//  2. **Flexibility for Diverse DTOs:** This design allows the 'ValidateDto' method to handle any arbitrary
//     struct type (e.g., UserCreateRequest, ProductUpdateRequest) as a DTO. All concrete types in Go implicitly
//     implement the empty interface 'interface{}' (aliased as 'any').
//
//  3. **Compatibility with Reflection-Based Validators:** Libraries like 'go-playground/validator' are reflection-based. They inspect the structure and tags of the DTO at runtime to perform validation.
//     Passing 'interface{}' enables these libraries to dynamically introspect the underlying type, making them highly flexible.
//     Generics primarily offer compile-time type safety and code reuse; they don't fundamentally change the
//     runtime reflection mechanism required by such validation libraries.
//
// Parameters:
//   - dto: A pointer to the struct that will receive the unmarshaled JSON data from the request body.
//   - request: The incoming `http.Request` containing the JSON body to be validated and decoded.
//
// Returns:
//   - *models.CustomError: Returns a pointer to a `CustomError` if any validation or decoding
//     error occurs. Returns `nil` if the request body is successfully
//     read, parsed, and unmarshaled into the `dto`.
func (v *Validator) ValidateRequestBody(dto interface{}, request *http.Request) (customError *models.CustomError) {

	// Read the entire request body into a byte slice.
	// This is necessary because the request.Body is an io.Reader and can only be read once.
	body, err := io.ReadAll(request.Body)
	if err != nil {
		// If there's an error reading the body (e.g., network issue, client disconnect),
		// create an internal server error.
		apiError := apierrors.ERROR_CODE_INTERNAL_SERVER.NewApiError("Failed to read request body", err)
		// Return the custom error.
		return apiError
	}
	// Defer closing the request body. This ensures the body is closed after the function
	// returns, releasing resources. This is good practice even if the router might close it.
	defer request.Body.Close()

	// Create a new JSON decoder for the read body.
	decoder := json.NewDecoder(bytes.NewReader(body))
	// Configure the decoder to return an error if the JSON contains fields not present
	// in the `dto` struct. This helps enforce strict API contracts.
	decoder.DisallowUnknownFields()

	// Attempt to decode the JSON request body into the provided DTO struct.
	// The `dto` parameter is already a pointer, so we pass it directly.
	if err := decoder.Decode(dto); err != nil {
		// Handle specific JSON decoding errors to provide more precise feedback to the client.

		// Check if the error is a JSON UnmarshalTypeError, which indicates a type mismatch
		// (e.g., expecting an integer but receiving a string).
		if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			// Create a specific API error for type mismatches, including the field and expected type.
			apiError := apierrors.ERROR_CODE_JSON_TYPE_MISMATCH.NewApiError(
				fmt.Sprintf("Invalid type for field '%s'. Expected type %s.",
					unmarshalErr.Field, unmarshalErr.Type.String()),
				err,
			)
			return apiError
		}

		// Check if the error is a JSON SyntaxError, which indicates malformed JSON.
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			// Create a specific API error for JSON syntax issues.
			apiError := apierrors.ERROR_CODE_JSON_SYNTAX.NewApiError(fmt.Sprintf("Invalid JSON syntax: %s", syntaxErr.Error()), err)
			return apiError
		}

		// Check if the error message indicates an unknown field.
		// This handles cases where `DisallowUnknownFields()` leads to an error containing this string.
		if strings.Contains(err.Error(), "json: unknown field") {
			// Extract the unknown field name from the error message and create a specific API error.
			apiError := apierrors.ERROR_CODE_JSON_UKNOWN_FIELD.NewApiError(strings.TrimPrefix(err.Error(), "json: unknown field "), err)
			return apiError
		}

		// Check if the error is io.EOF, which indicates an empty request body.
		if err == io.EOF {
			// Create a specific API error for an empty request body.
			apiError := apierrors.ERROR_CODE_BAD_REQUEST.NewApiError("Request body is empty", err)
			return apiError
		}

		// For any other unhandled decoding errors, return a generic bad request error.
		apiError := apierrors.ERROR_CODE_BAD_REQUEST.NewApiError("Body decoding failed", err)
		return apiError
	}

	// If decoding is successful and no errors occurred, return nil,
	// indicating that the request body was valid and the DTO was populated.
	return nil
}

// ValidateDto validates a Data Transfer Object (DTO).
//
// This method accepts 'dto' as an 'interface{}' type for the following reasons:
//
//  1. **Go Language Limitation:** In Go, methods defined on a *non-generic* struct (like 'Validator' in this case)
//     cannot have their own type parameters. If 'Validator' were 'Validator[T any]', then its methods could use 'T'.
//     However, since 'Validator' itself doesn't need to be generic, the method must accept a concrete type or an interface.
//     Using [T any](dto[T) would result in a compile-time error.
//
//  2. **Flexibility for Diverse DTOs:** This design allows the 'ValidateDto' method to handle any arbitrary
//     struct type (e.g., UserCreateRequest, ProductUpdateRequest) as a DTO. All concrete types in Go implicitly
//     implement the empty interface 'interface{}' (aliased as 'any').
//
//  3. **Compatibility with Reflection-Based Validators:** Libraries like 'go-playground/validator' are reflection-based. They inspect the structure and tags of the DTO at runtime to perform validation.
//     Passing 'interface{}' enables these libraries to dynamically introspect the underlying type, making them highly flexible.
//     Generics primarily offer compile-time type safety and code reuse; they don't fundamentally change the
//     runtime reflection mechanism required by such validation libraries.
//
// paramaters:
//   - dto: The Data Transfer Object (DTO) to be validated.
//
// returns:
//   - *models.CustomError: A CustomError if validation fails, nil otherwise.
func (v *Validator) ValidateDto(dto interface{}) *models.CustomError {
	// Attempt to validate the DTO using the underlying validation library.
	if err := v.validator.Struct(dto); err != nil {
		// Differentiate between validation errors and other unexpected errors.
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// Pre-allocate slice capacity to optimize for common validation scenarios,
			// reducing reallocations as validation error details are appended.
			details := make([]models.ValidationErrorDetail, 0, len(validationErrors))

			// Process each individual validation error to extract and format relevant details
			// for a user-friendly API response.
			for _, validationError := range validationErrors {
				detail := models.ValidationErrorDetail{
					Field:   validationError.Field(),
					Tag:     validationError.Tag(),
					Param:   validationError.Param(),
					Message: GetValidationErrorMessage(validationError), // Custom message for specific validation tags
				}
				details = append(details, detail)
			}

			// Construct an API error for unprocessable entities (HTTP 422),
			// embedding the specific validation failure details.
			apiError := apierrors.ERROR_CODE_UNPROCESSABLE_ENTITY.NewApiError(details, err)
			return apiError
		}

		// Handle any non-validation errors as internal server errors (HTTP 500).
		// These typically indicate unexpected issues within the validation process itself.
		internalError := apierrors.ERROR_CODE_INTERNAL_SERVER.NewApiError(err.Error(), err)
		return internalError
	}

	// No errors, so validation succeeded.
	return nil
}
