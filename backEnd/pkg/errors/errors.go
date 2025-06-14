package errors

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
)

type ErrorCode string

const (
	ERROR_CODE_GENERIC             ErrorCode = "GENERIC_ERROR"
	ERROR_CODE_VALIDATION_FAILED   ErrorCode = "VALIDATION_FAILED"
	ERROR_CODE_NOT_FOUND           ErrorCode = "NOT_FOUND"
	ERROR_CODE_UNAUTHORIZED        ErrorCode = "UNAUTHORIZED"
	ERROR_CODE_FORBIDDEN           ErrorCode = "FORBIDDEN"
	ERROR_CODE_INTERNAL_SERVER     ErrorCode = "INTERNAL_SERVER_ERROR"
	ERROR_CODE_DATABASE_ERROR      ErrorCode = "DATABASE_ERROR"
	ERROR_CODE_INVALID_INPUT       ErrorCode = "INVALID_INPUT"
	ERROR_CODE_SERVICE_UNAVAILABLE ErrorCode = "SERVICE_UNAVAILABLE"
	ERROR_CODE_DATABASE_INITIALIZE ErrorCode = "DATABASE_INITIALIZE"
	ERROR_CODE_ENCODE_ERROR        ErrorCode = "ENCODE_ERROR"
	ERROR_UNIQUE_KEY_VIOLATION     ErrorCode = "UNIQUE_KEY_VIOLATION"
)

var ErrorMessages = map[ErrorCode]struct {
	Message string
	Status  int
}{
	ERROR_CODE_GENERIC:             {Message: "An unexpected error occurred.", Status: http.StatusInternalServerError},
	ERROR_CODE_VALIDATION_FAILED:   {Message: "One or more input values are invalid.", Status: http.StatusBadRequest},
	ERROR_CODE_NOT_FOUND:           {Message: "The requested resource could not be found.", Status: http.StatusNotFound},
	ERROR_CODE_UNAUTHORIZED:        {Message: "Authentication required or invalid credentials.", Status: http.StatusUnauthorized},
	ERROR_CODE_FORBIDDEN:           {Message: "You do not have permission to perform this action.", Status: http.StatusForbidden},
	ERROR_CODE_INTERNAL_SERVER:     {Message: "An internal server error occurred.", Status: http.StatusInternalServerError},
	ERROR_CODE_DATABASE_ERROR:      {Message: "A database operation failed.", Status: http.StatusInternalServerError},
	ERROR_CODE_INVALID_INPUT:       {Message: "The provided input is malformed.", Status: http.StatusBadRequest},
	ERROR_CODE_SERVICE_UNAVAILABLE: {Message: "The service is temporarily unavailable.", Status: http.StatusServiceUnavailable},
	ERROR_CODE_DATABASE_INITIALIZE: {Message: "Database initialization failed.", Status: http.StatusInternalServerError},
	ERROR_CODE_ENCODE_ERROR:        {Message: "Error encoding response.", Status: http.StatusInternalServerError},
	ERROR_UNIQUE_KEY_VIOLATION:     {Message: "A unique key violation occurred.", Status: http.StatusConflict},
}

func (code ErrorCode) ApiErrorResponse(details any, err error) *models.CustomError {
	responseData := models.CustomError{}
	responseInfo, ok := ErrorMessages[code]

	// If the code is not found in the map, return a generic error
	if !ok {
		responseData.Message = ErrorMessages[ERROR_CODE_GENERIC].Message
		responseData.HTTPStatusCode = ErrorMessages[ERROR_CODE_GENERIC].Status
		responseData.Code = string(ERROR_CODE_GENERIC)
		return &responseData
	}

	// If the code is found in the map, return the corresponding error
	return &models.CustomError{
		Code:           string(code),
		Message:        responseInfo.Message,
		HTTPStatusCode: responseInfo.Status,
		//Details and Err can be nil because They are omitted by default
		Details: details,
		Err:     err,
	}

}
