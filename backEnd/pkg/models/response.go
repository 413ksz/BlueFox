package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   any    `json:"value,omitempty"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

type CustomError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
	Err     error  `json:"-"`
}

type Pagination struct {
	TotalItems   int     `json:"totalItems"`
	ItemsPerPage *int    `json:"itemsPerPage,omitempty"`
	PageIndex    *int    `json:"pageIndex,omitempty"`
	NextLink     *string `json:"nextLink,omitempty"`
	PreviousLink *string `json:"previousLink,omitempty"`
}

type ResponseData[TItemtype any] struct {
	Deleted    *bool       `json:"deleted,omitempty"`
	Updated    *time.Time  `json:"updated,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Items      []TItemtype `json:"items"`
}

type ApiResponse[TItemtype any] struct {
	Params     map[string]interface{}   `json:"params,omitempty"`
	Data       *ResponseData[TItemtype] `json:"data,omitempty"`
	Error      *CustomError             `json:"error,omitempty"`
	Message    string                   `json:"message,omitempty"`
	StatusCode int                      `json:"-"`
}

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

func (apiResponse *ApiResponse[TItemtype]) WithData(Message string, Data *ResponseData[TItemtype], statusCode int) {
	apiResponse.Message = Message
	apiResponse.Data = Data
	apiResponse.StatusCode = statusCode
	apiResponse.Error = nil
}

func (apiResponse *ApiResponse[TItemtype]) WithError(Message string, Error *CustomError, statusCode int) {
	apiResponse.Message = Message
	apiResponse.Error = Error
	apiResponse.StatusCode = statusCode
	apiResponse.Data = nil
}

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
