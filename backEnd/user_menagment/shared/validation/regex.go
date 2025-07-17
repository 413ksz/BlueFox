package validation

import "regexp"

// Regex patterns defined as constants for clarity and reusability.
const (
	// USERNAME_PATTERN defines the regex for valid usernames.
	// It allows alphanumeric characters, including underscores or hyphens in the middle,
	// but prevents them at the very beginning or end.
	// Ensures a minimum length of 2 and allows for consecutive hyphens/underscores.
	USERNAME_PATTERN = `^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$`

	// NAME_PATTERN defines the regex for first and last names.
	// It requires names to start and end with a letter.
	// Allowed characters in between include letters, apostrophes ('), and hyphens (-).
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

	// EMAIL_PATTERN defines the regex for valid email addresses.
	EMAIL_PATTERN = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// BCRYPT_PATTERN defines the regex for valid bcrypt hashes.
	BCRYPT_PATTERN = `^\$2[aby]\$\d{2}\$[./0-9A-Za-z]{22}[./0-9A-Za-z]{31}$`
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
	emailRegex             = regexp.MustCompile(EMAIL_PATTERN)
	bcryptRegex            = regexp.MustCompile(BCRYPT_PATTERN)
)
