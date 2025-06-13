package models

import (
	"time"
)

type CustomError struct {
	Domain  string `json:"domain"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
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
	Context *string                  `json:"context,omitempty"`
	Method  *string                  `json:"method,omitempty"`
	Params  map[string]interface{}   `json:"params,omitempty"`
	Data    *ResponseData[TItemtype] `json:"data,omitempty"`
	Error   *CustomError             `json:"error,omitempty"`
}

func NewApiResponse[TItemtype any]() *ApiResponse[TItemtype] {
	return &ApiResponse[TItemtype]{}
}
