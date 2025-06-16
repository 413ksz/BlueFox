package validation

import (
	"regexp"
)

// Regex patterns defined as constants for clarity and reusability.
const (
	// USERNAME_PATTERN defines the regex for valid usernames.
	// It allows 3-20 alphanumeric characters, including underscores or hyphens in the middle,
	// but prevents them at the very beginning or end.
	// Ensures a minimum length of 3 and prevents leading/trailing/double hyphens/underscores.
	USERNAME_PATTERN = `^[a-zA-Z0-9][a-zA-Z0-9_-]{1,18}[a-zA-Z0-9]$`

	// EMAIL_PATTERN defines the regex for valid email addresses.
	// It's a common pattern that covers most standard email formats, with strict local part rules.
	// Prevents leading, trailing, or consecutive dots in the local part.
	EMAIL_PATTERN = `^[a-zA-Z0-9_+-]+(?:\.[a-zA-Z0-9_+-]+)*@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// NAME_PATTERN defines the regex for first and last names.
	// It requires names to be 3-20 characters long, starting and ending with a letter.
	// Allowed characters in between include letters, apostrophes ('), hyphens (-), and spaces.
	NAME_PATTERN = `^[a-zA-Z][a-zA-Z' -]{1,18}[a-zA-Z]$`

	// PASSWORD_LENGTH_PATTERN checks for a minimum password length of 8 characters.
	PASSWORD_LENGTH_PATTERN = `.{8,}`
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
	emailRegex             = regexp.MustCompile(EMAIL_PATTERN)
	nameRegex              = regexp.MustCompile(NAME_PATTERN)
	passwordLengthRegex    = regexp.MustCompile(PASSWORD_LENGTH_PATTERN)
	passwordUppercaseRegex = regexp.MustCompile(PASSWORD_UPPERCASE_PATTERN)
	passwordLowercaseRegex = regexp.MustCompile(PASSWORD_LOWERCASE_PATTERN)
	passwordNumberRegex    = regexp.MustCompile(PASSWORD_NUMBER_PATTERN)
	passwordSpecialRegex   = regexp.MustCompile(PASSWORD_SPECIAL_PATTERN)
)

// ValidateUsername checks if the provided username string matches the USERNAME_PATTERN.
// It returns true if valid, false otherwise.
// @param username: The username string to validate.
// @return bool: True if the username is valid, false otherwise.
func ValidateUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

// ValidateEmail checks if the provided email string matches the EMAIL_PATTERN.
// It returns true if valid, false otherwise.
// @param email: The email string to validate.
// @return bool: True if the email is valid, false otherwise.
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if the provided password string meets all defined criteria.
// It uses logical AND (&&) to ensure the password matches all individual regex patterns
// for length, uppercase, lowercase, numbers, and special characters.
// It returns true if all criteria are met, false otherwise.
// @param password: The password string to validate.
// @return bool: True if the password is valid, false otherwise.
func ValidatePassword(password string) bool {
	return passwordLengthRegex.MatchString(password) &&
		passwordUppercaseRegex.MatchString(password) &&
		passwordLowercaseRegex.MatchString(password) &&
		passwordNumberRegex.MatchString(password) &&
		passwordSpecialRegex.MatchString(password)
}

// ValidateName checks if a provided name string (e.g., first name or last name)
// matches the NAME_PATTERN.
// It returns true if valid, false otherwise.
// @param name: The name string to validate.
// @return bool: True if the name is valid, false otherwise.
func ValidateName(name string) bool {
	return nameRegex.MatchString(name)
}
