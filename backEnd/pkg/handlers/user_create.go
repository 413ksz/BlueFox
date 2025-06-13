package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
)

const (
	PASSWORD_REPLACEMENT = "Nice try"
)

var (
	METHOD  = "PUT"
	CONTEXT = "/api/user"
)

// UserCreateHandler handles HTTP PUT requests for creating a new user.
// It expects a JSON request body containing user data, saves it to the database,
// and returns the created user object with a 201 Created status.
func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	apiResponse := &models.ApiResponse[models.User]{}

	apiResponse.Method = &METHOD
	apiResponse.Context = &CONTEXT
	// Get the GORM database instance from your database package.
	db := database.DB

	// Check if the database connection is initialized.
	if db == nil {
		log.Println("Database connection is not initialized in UserCreateHandler.")
		http.Error(w, "Internal server error: Database not ready", http.StatusInternalServerError)
		return
	}

	var newUser models.User

	// Decode the JSON request body into the newUser struct.
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Printf("Error decoding request body in UserCreateHandler: %v", err)
		http.Error(w, "Bad request: Invalid JSON data", http.StatusBadRequest)
		return
	}

	apiResponse.Params = map[string]interface{}{
		"userName":    newUser.Username,
		"email":       newUser.Email,
		"password":    PASSWORD_REPLACEMENT,
		"lastName":    newUser.LastName,
		"firstName":   newUser.FirstName,
		"dateOfBirth": newUser.DateOfBirth,
	}

	// Perform validation checks on the user data.
	if newUser.Email == "" || newUser.Username == "" || newUser.Password == "" || newUser.FirstName == "" || newUser.LastName == "" || newUser.DateOfBirth.IsZero() {
		log.Println("Validation error: required fields are missing.")
		http.Error(w, "Bad request: Missing required fields", http.StatusBadRequest)
		return
	}

	if !validation.ValidateEmail(newUser.Email) {
		log.Println("Validation error: invalid email.")
		http.Error(w, "Bad request: Invalid email", http.StatusBadRequest)
		return
	}

	if !validation.ValidatePassword(newUser.Password) {
		log.Println("Validation error: invalid password.")
		http.Error(w, "Bad request: Invalid password", http.StatusBadRequest)
		return
	}

	if !validation.ValidateUsername(newUser.Username) {
		log.Println("Validation error: invalid username.")
		http.Error(w, "Bad request: Invalid username", http.StatusBadRequest)
		return
	}

	if !validation.ValidateName(newUser.FirstName) {
		log.Println("Validation error: invalid first name.")
		http.Error(w, "Bad request: Invalid first name", http.StatusBadRequest)
		return
	}

	if !validation.ValidateName(newUser.LastName) {
		log.Println("Validation error: invalid last name.")
		http.Error(w, "Bad request: Invalid last name", http.StatusBadRequest)
		return
	}

	passwordHash, err := passwordHashing.HashPassword(newUser.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error: Failed to hash password", http.StatusInternalServerError)
		return
	}
	newUser.Password = passwordHash

	// Attempt to create (insert) the new user record into the database using GORM.
	result := db.Create(&newUser)

	// Check for errors returned by the database operation.
	if result.Error != nil {
		// Log the detailed database error.
		log.Printf("Database error creating user: %v", result.Error)

		http.Error(w, "Internal server error: Failed to create user", http.StatusInternalServerError)
		return
	}

	newUser.Password = PASSWORD_REPLACEMENT
	// Add the newly created user to the response data.
	apiResponse.Data = &models.ResponseData[models.User]{}
	apiResponse.Data.Items = make([]models.User, 0)
	apiResponse.Data.Items = append(apiResponse.Data.Items, newUser)

	// Set the Content-Type header to indicate that the response body is JSON.
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code to 201 Created, which is standard for successful resource creation.
	w.WriteHeader(http.StatusCreated)

	// Encode the newly created user object back into JSON and write it to the response body.
	err = json.NewEncoder(w).Encode(apiResponse)
	if err != nil {
		log.Printf("Error encoding JSON response for created user: %v", err)
	}
}
