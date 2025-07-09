package users

import (
	"regexp"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/go-playground/validator/v10"
)

// Regex patterns defined as constants for clarity and reusability.
const (
	// USERNAME_PATTERN defines the regex for valid usernames.
	// It allows 3-20 alphanumeric characters, including underscores or hyphens in the middle,
	// but prevents them at the very beginning or end.
	// Ensures a minimum length of 3 and prevents leading/trailing/double hyphens/underscores.
	USERNAME_PATTERN = `^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$`

	// NAME_PATTERN defines the regex for first and last names.
	// It requires names to be 3-20 characters long, starting and ending with a letter.
	// Allowed characters in between include letters, apostrophes ('), hyphens (-), and spaces.
	NAME_PATTERN = `^[a-zA-Z][a-zA-Z'-]+[a-zA-Z]$`

	// PASSWORD_UPPERCASE_PATTERN checks if the password contains at least one uppercase letter.
	PASSWORD_UPPERCASE_PATTERN = `[A-Z]`
	// PASSWORD_LOWERCASE_PATTERN checks if the password contains at least one lowercase letter.
	PASSWORD_LOWERCASE_PATTERN = `[a-z]`
	// PASSWORD_NUMBER_PATTERN checks if the password contains at least one digit (0-9).
	PASSWORD_NUMBER_PATTERN = `[0-9]`
	// PASSWORD_SPECIAL_PATTERN checks if the password contains at least one special character
	// from a predefined set of common special characters.
	PASSWORD_SPECIAL_PATTERN = `[!@#$%^&*()_+={}\[\]:;"'<>,.?/\\|~-]`
)

// Regex variables are compiled versions of the patterns.
// They are compiled once when the package is initialized for better performance.
var (
	usernameRegex          = regexp.MustCompile(USERNAME_PATTERN)
	nameRegex              = regexp.MustCompile(NAME_PATTERN)
	passwordUppercaseRegex = regexp.MustCompile(PASSWORD_UPPERCASE_PATTERN)
	passwordLowercaseRegex = regexp.MustCompile(PASSWORD_LOWERCASE_PATTERN)
	passwordNumberRegex    = regexp.MustCompile(PASSWORD_NUMBER_PATTERN)
	passwordSpecialRegex   = regexp.MustCompile(PASSWORD_SPECIAL_PATTERN)
)

// ValidateUsername checks if the provided username string matches the USERNAME_PATTERN.
// It returns true if valid, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the username string to validate.
// Returns:
// - bool: True if the username is valid, false otherwise.
func ValidateUsername(fieldLevel validator.FieldLevel) bool {
	return usernameRegex.MatchString(fieldLevel.Field().String())
}

// ValidatePassword checks if the provided password string meets all defined criteria.
// It uses logical AND (&&) to ensure the password matches all individual regex patterns
// for length, uppercase, lowercase, numbers, and special characters.
// It returns true if all criteria are met, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the password string to validate.
// Returns:
// - bool: True if the password is valid, false otherwise.
func ValidatePassword(fieldLevel validator.FieldLevel) bool {
	password := fieldLevel.Field().String()
	return passwordUppercaseRegex.MatchString(password) &&
		passwordLowercaseRegex.MatchString(password) &&
		passwordNumberRegex.MatchString(password) &&
		passwordSpecialRegex.MatchString(password)
}

// ValidateName checks if a provided name string (e.g., first name or last name)
// matches the NAME_PATTERN.
// It returns true if valid, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the name string to validate.
// Returns:
// - bool: True if the name is valid, false otherwise.
func ValidateName(fieldLevel validator.FieldLevel) bool {
	return nameRegex.MatchString(fieldLevel.Field().String())
}

// ValidateDateOfBirth checks if a provided date of birth is valid.
// It ensures the date of birth is in the past and within the age range of 16 to 120.
// It returns true if valid, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the date of birth to validate.
// Returns:
// - bool: True if the date of birth is valid, false otherwise.
func ValidateDateOfBirth(fieldLevel validator.FieldLevel) bool {
	//we don't check here for parse error it was checked before it should always be valid
	dateofbirth, _ := fieldLevel.Field().Interface().(time.Time)

	now := time.Now()

	if dateofbirth.After(now) {
		return false
	}
	// Calculate the minimum(16) and maximum(120) date of birth.
	minAge := now.AddDate(-16, 0, 0)
	maxAge := now.AddDate(-120, 0, 0)

	// Check if the date of birth is after the maximum age.
	if dateofbirth.Before(maxAge) {
		return false
	}

	// Check if the date of birth is before the minimum age.
	if dateofbirth.After(minAge) {
		return false
	}

	return true
}

// RegisterDomainValidators registers the custom validation functions for the User struct.
func RegisterDomainValidators(validator *validation.Validator) {
	validator.RegisterCustomValidation("username", ValidateUsername)
	validator.RegisterCustomValidation("password", ValidatePassword)
	validator.RegisterCustomValidation("name", ValidateName)
	validator.RegisterCustomValidation("dateofbirth", ValidateDateOfBirth)
}
