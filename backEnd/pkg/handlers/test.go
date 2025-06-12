package handlers

import (
	"encoding/json"
	"net/http"
)

type TestResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	// Create an instance of our response struct
	response := TestResponse{
		Message: "Hello from the /api/test endpoint! Routing works!",
		Status:  "success",
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code
	w.WriteHeader(http.StatusOK)

	// Encode the response struct to JSON and write it to the response writer
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
