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

	// EMAIL_PATTERN defines the regex for valid email addresses.
	EMAIL_PATTERN = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// ARGON2ID_PATTERN defines the regex for valid Argon2 ID hashes.
	ARGON2ID_PATTERN = `^\$argon2id\$v=\d+\$m=\d+,t=\d+,p=\d+\$[A-Za-z0-9+/=]+\$[A-Za-z0-9+/=]+$`
)

// Regex variables are compiled versions of the patterns.
// They are compiled once when the package is initialized for better performance.
var (
	usernameRegex = regexp.MustCompile(USERNAME_PATTERN)
	nameRegex     = regexp.MustCompile(NAME_PATTERN)
	emailRegex    = regexp.MustCompile(EMAIL_PATTERN)
	argon2IDRegex = regexp.MustCompile(ARGON2ID_PATTERN)
)
