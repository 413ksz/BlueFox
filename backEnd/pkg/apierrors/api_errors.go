package apierrors

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/rs/zerolog"
)

type ErrorCode string

// Error implements the error interface
func (e ErrorCode) Error() string {
	// You can return a default message here, or just the string representation of the code
	if msg, ok := ErrorMessages[e]; ok {
		return msg.Message
	}
	return string(e) // Fallback if code isn't in map
}

const (
	ERROR_CODE_GENERIC                       ErrorCode = "GENERIC_ERROR"
	ERROR_CODE_JSON_SYNTAX                   ErrorCode = "JSON_SYNTAX_ERROR"
	ERROR_CODE_JSON_TYPE_MISMATCH            ErrorCode = "JSON_MISMATCH_ERROR"
	ERROR_CODE_JSON_UKNOWN_FIELD             ErrorCode = "JSON_UNKNOWN_FIELD_ERROR"
	ERROR_CODE_BAD_REQUEST                   ErrorCode = "BAD_REQUEST"
	ERROR_CODE_UNPROCESSABLE_ENTITY          ErrorCode = "UNPROCESSABLE_ENTITY"
	ERROR_CODE_NOT_FOUND                     ErrorCode = "NOT_FOUND"
	ERROR_CODE_UNAUTHORIZED                  ErrorCode = "UNAUTHORIZED"
	ERROR_CODE_FORBIDDEN                     ErrorCode = "FORBIDDEN"
	ERROR_CODE_INTERNAL_SERVER               ErrorCode = "INTERNAL_SERVER_ERROR"
	ERROR_CODE_DATABASE_ERROR                ErrorCode = "DATABASE_ERROR"
	ERROR_CODE_SERVICE_UNAVAILABLE           ErrorCode = "SERVICE_UNAVAILABLE"
	ERROR_CODE_DATABASE_INITIALIZE           ErrorCode = "DATABASE_INITIALIZE"
	ERROR_CODE_UNIQUE_KEY_VIOLATION          ErrorCode = "UNIQUE_KEY_VIOLATION"
	ERROR_CODE_ENVIREMENT_VARIABLE_NOT_FOUND ErrorCode = "ENVIREMENT_VARIABLE_NOT_FOUND"
	ERROR_CODE_VALIDATION_REGISTRATION_ERROR ErrorCode = "VALIDATION_REGISTRATION_ERROR"
)

var ErrorMessages = map[ErrorCode]struct {
	Message  string
	HttpCode int
	LogLevel zerolog.Level
}{
	ERROR_CODE_GENERIC:                       {Message: "An unexpected error occurred."},
	ERROR_CODE_BAD_REQUEST:                   {Message: "The request was invalid or malformed.", HttpCode: http.StatusBadRequest, LogLevel: zerolog.WarnLevel},
	ERROR_CODE_UNPROCESSABLE_ENTITY:          {Message: "One or more input values are invalid", HttpCode: http.StatusUnprocessableEntity, LogLevel: zerolog.WarnLevel},
	ERROR_CODE_JSON_SYNTAX:                   {Message: "The request body contains malformed JSON or invalid JSON syntax.", HttpCode: http.StatusBadRequest, LogLevel: zerolog.WarnLevel},
	ERROR_CODE_JSON_TYPE_MISMATCH:            {Message: "The request body contains a field with an unexpected type.", HttpCode: http.StatusBadRequest, LogLevel: zerolog.WarnLevel},
	ERROR_CODE_JSON_UKNOWN_FIELD:             {Message: "The request body contains an unknown field.", HttpCode: http.StatusBadRequest, LogLevel: zerolog.WarnLevel},
	ERROR_CODE_NOT_FOUND:                     {Message: "The requested resource could not be found."},
	ERROR_CODE_UNAUTHORIZED:                  {Message: "Authentication failed."},
	ERROR_CODE_FORBIDDEN:                     {Message: "You do not have permission to perform this action."},
	ERROR_CODE_INTERNAL_SERVER:               {Message: "An internal server error occurred.", HttpCode: http.StatusInternalServerError, LogLevel: zerolog.ErrorLevel},
	ERROR_CODE_DATABASE_ERROR:                {Message: "A database operation failed."},
	ERROR_CODE_SERVICE_UNAVAILABLE:           {Message: "The service is temporarily unavailable."},
	ERROR_CODE_DATABASE_INITIALIZE:           {Message: "Database initialization failed."},
	ERROR_CODE_UNIQUE_KEY_VIOLATION:          {Message: "A unique key violation occurred."},
	ERROR_CODE_ENVIREMENT_VARIABLE_NOT_FOUND: {Message: "Environment variable not found."},
	ERROR_CODE_VALIDATION_REGISTRATION_ERROR: {Message: "The validator domain spesific registration failed.", LogLevel: zerolog.FatalLevel},
}

func (code ErrorCode) NewApiError(details any, err error) *models.CustomError {
	responseData := models.CustomError{}
	responseInfo, ok := ErrorMessages[code]

	// If the code is not found in the map, return a generic error
	if !ok {
		responseData.Message = ErrorMessages[ERROR_CODE_GENERIC].Message
		responseData.Code = string(ERROR_CODE_GENERIC)
		return &responseData
	}

	// If the code is found in the map, return the corresponding error
	return &models.CustomError{
		Code:    string(code),
		Message: responseInfo.Message,
		//Details and Err can be nil because They are omitted by default
		HttpCode: responseInfo.HttpCode,
		LogLevel: responseInfo.LogLevel,
		Details:  details,
		Err:      err,
	}
}
