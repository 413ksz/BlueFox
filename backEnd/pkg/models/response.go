package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type CustomError struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	Details        any    `json:"details,omitempty"`
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
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
	Context    string                   `json:"context,omitempty"`
	Method     string                   `json:"method,omitempty"`
	Params     map[string]interface{}   `json:"params,omitempty"`
	Data       *ResponseData[TItemtype] `json:"data,omitempty"`
	Error      *CustomError             `json:"error,omitempty"`
	Message    string                   `json:"message,omitempty"`
	StatusCode int                      `json:"-"`
}

func NewApiResponse[TItemtype any]() *ApiResponse[TItemtype] {
	return &ApiResponse[TItemtype]{}
}

func SendApiResponse[T any](w http.ResponseWriter, apiResponse *ApiResponse[T]) {
	// Write the determined HTTP status code to the response header
	if apiResponse.Error != nil {
		w.WriteHeader(apiResponse.Error.HTTPStatusCode)
	} else {
		w.WriteHeader(apiResponse.StatusCode)
	}
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the apiResponse struct to JSON and write it to the response body
	err := json.NewEncoder(w).Encode(apiResponse)
	if err != nil {
		// Log an error if encoding fails, but don't attempt to write to w again
		// as headers have already been sent.
		log.Printf("Error encoding API response: %v", err)
	}
}
