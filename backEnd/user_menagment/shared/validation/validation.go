package validation

import (
	"math"
	"time"
	"unicode"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/go-playground/validator/v10"
)

// ValidateUsernameField checks if the provided username string matches the USERNAME_PATTERN.
// It returns true if valid, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the username string to validate.
// Returns:
// - bool: True if the username is valid, false otherwise.
func ValidateUsernameField(fieldLevel validator.FieldLevel) bool {
	return usernameRegex.MatchString(fieldLevel.Field().String())
}

// ValidateUsernameString checks if the provided username string matches the username validation rules.
// It returns a ValidationError if the username is invalid, nil otherwise.
// Parameters:
// - username: The username string to validate.
// Returns:
// - *ValidationError: A ValidationError if the username is invalid, nil otherwise.
func ValidateUsernameString(username string) *models.ValidationError {
	// Check if the username is empty.
	if username == "" {
		validationError := models.NewValidationError("username", username, "required", "Username is required.")
		return validationError
	}
	// Check if the username is less than 3 characters.
	if len(username) < 3 {
		validationError := models.NewValidationError("username", username, "min", "Username cannot be less than 3 characters.")
		return validationError
	}
	// Check if the username exceeds 30 characters.
	if len(username) > 30 {
		validationError := models.NewValidationError("username", username, "max", "Username cannot exceed 30 characters.")
		return validationError
	}
	// Check if the username format is valid.
	if !usernameRegex.MatchString(username) {
		validationError := models.NewValidationError("username", username, "regex", "Invalid username format.")
		return validationError
	}
	return nil
}

// ValidateEmailField checks if the provided email field string matches the EMAIL_PATTERN.
// It uses the validator.FieldLevel to access the email field value.
// Therefore, it should be called on dtos that have validation tags like `validate:"email"`
// It returns true if valid, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the email string to validate.
// Returns:
// - bool: True if the email is valid, false otherwise.
func ValidateEmailField(fieldLevel validator.FieldLevel) bool {
	return emailRegex.MatchString(fieldLevel.Field().String())
}

// ValidateEmailString checks if the provided email string matches the email validation rules.
// It returns a ValidationError if the email is invalid, nil otherwise.
// Parameters:
// - email: The email string to validate.
// Returns:
// - *ValidationError: A ValidationError if the email is invalid, nil otherwise.
func ValidateEmailString(email string) *models.ValidationError {
	// Check if the email is empty.
	if email == "" {
		validationError := models.NewValidationError("email", email, "required", "Email is required.")
		return validationError
	}
	// Check if the email exceeds 254 characters for rfc5322.
	if len(email) > 254 {
		validationError := models.NewValidationError("email", email, "max", "Email cannot exceed 254 characters.")
		return validationError
	}
	// Check if the email format is valid.
	if !emailRegex.MatchString(email) {
		validationError := models.NewValidationError("email", email, "regex", "Invalid email format.")
		return validationError
	}

	return nil
}

// ValidatePasswordField checks if the provided password string meets all defined criteria.
// It uses logical AND (&&) to ensure the password matches all individual regex patterns
// for length, uppercase, lowercase, numbers, and special characters.
// It returns true if all criteria are met, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the password string to validate.
// Returns:
// - bool: True if the password is valid, false otherwise.
func ValidatePasswordField(fieldLevel validator.FieldLevel) bool {
	password := fieldLevel.Field().String()
	return passwordUppercaseRegex.MatchString(password) &&
		passwordLowercaseRegex.MatchString(password) &&
		passwordNumberRegex.MatchString(password) &&
		passwordSpecialRegex.MatchString(password)
}

// ValidatePasswordForEntropy checks if the provided password string meets all defined criteria.
// It uses Shannon entropy to check the relative password strength of 100 bits of entropy.
// It also checks if the password length is between 16 and 72 characters.
// parameters:
// - password: the password string to validate
// returns:
// - *ValidationError: a ValidationError if the password is invalid, nil otherwise
func ValidatePasswordForEntropy(password string) *models.ValidationError {

	// ---------- Password Length Validation ----------
	// Adheres to and exceeds current NIST SP 800-63B guidelines (Section 5.1.1.2) for memorized secrets,
	// which recommend a minimum length of 8 characters. We enforce a minimum of 16 characters for
	// enhanced security against brute-force attacks and future-proofing and a maximum of 72 characters
	// for bcrypt compatibility.

	// Enforce a minimum password length
	if password == "" {
		validationError := models.NewValidationError("password", password, "required", "Password is required.")
		return validationError
	}
	if len(password) < 16 {
		validationError := models.NewValidationError("password", password, "min", "Password cannot be less than 16 characters.")
		return validationError

	}
	// Enforce a maximum password length of 72 characters for bcrypt compatibility.
	if len(password) > 72 {
		validationError := models.NewValidationError("password", password, "max", "Password cannot exceed 72 characters.")
		return validationError
	}

	// --------- Check password entropy ---------
	// Shannon entropy is a measure of the amount of information required to represent a random event.
	// It quantifies the uncertainty of a random variable.
	// A higher entropy value indicates a more complex and unpredictable pattern in the password,
	// making it more secure against brute-force attacks.
	// We enforce a minimum password entropy of 100 bits. While NIST SP 800-63B primarily focuses
	// on length and pattern avoidance, this entropy target aligns with common industry best practices
	// for strong memorized secrets.
	//
	// Note on character sets:
	// We use unicode.Is* functions to detect character categories (lowercase, uppercase, number, symbol) present in the password.
	// However, for the purpose of calculating the base of the logarithm (charCount) for Shannon entropy,
	// we assume fixed character set sizes for common categories (e.g., 26 for Latin lowercase, 10 for Arabic numerals).
	// This provides a practical and commonly understood minimum useful entropy estimation, even though
	// unicode.Is* can technically match broader ranges of characters.
	// therefore, it's not purely accurate because unicode.Is* can encompass a much wider range of characters than the fixed counts assume,
	// potentially leading to an underestimation of the true character set size for highly diverse passwords

	// minimum password entropy in bits
	const minEntropyBits = 100.0
	// variables for password entropy calculation
	var (
		hasLowercase         = false
		hasUppercase         = false
		hasNumber            = false
		hasSymbol            = false
		charCount    float64 = 0
	)

	// The following loop determine which character categories (lowercase, uppercase, number, symbol) are
	// present in the password. These findings are then used to calculate the size of the
	// character pool for entropy calculation
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowercase = true
		} else if unicode.IsUpper(char) {
			hasUppercase = true
		} else if unicode.IsNumber(char) {
			hasNumber = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSymbol = true
		}
	}

	// --------- check password character categories ---------

	if hasLowercase {
		charCount += 26 // Assumed character set size for lowercase Latin letters (a-z)
	}

	if hasUppercase {
		charCount += 26 // Assumed character set size for uppercase Latin letters (A-Z)
	}

	if hasNumber {
		charCount += 10 // Assumed character set size for Arabic numerals (0-9)
	}
	if hasSymbol {
		charCount += 32 // Assumed character set size for common symbols (e.g., ASCII punctuation and symbols)
	}

	// Check if the password contains at least one character category.
	if charCount == 0 {
		validationError := models.NewValidationError("password", password, "charCategories", "Password must contain at least one recognized character category (lowercase, uppercase, number, or symbol).")
		return validationError
	}

	// Calculate password entropy.
	length := float64(len(password))
	entropy := length * math.Log2(charCount)

	// Check if the password has a minimum entropy.
	if entropy < minEntropyBits {
		validationError := models.NewValidationError("password", password, "entropy", "Password must have a minimum entropy of 100 bits.")
		return validationError
	}

	return nil
}

// ValidatePasswordHash checks if the provided password hash string matches the validation rules.
// It returns a ValidationError if the password hash is invalid, nil otherwise.
// Parameters:
// - passwordHash: The password hash string to validate.
// Returns:
// - *ValidationError: A ValidationError if the password hash is invalid, nil otherwise.
func ValidatePasswordHash(passwordHash string) *models.ValidationError {
	// Check if the password hash is empty.
	if passwordHash == "" {
		validationError := models.NewValidationError("passwordHash", passwordHash, "required", "Password hash is required.")
		return validationError
	}
	// Check if the password hash format is valid.
	if !bcryptRegex.MatchString(passwordHash) {
		validationError := models.NewValidationError("passwordHash", passwordHash, "regex", "Invalid password hash format.")
		return validationError
	}
	return nil
}

// ValidateNameField checks if a provided name string (e.g., first name or last name)
// matches the NAME_PATTERN.
// It returns true if valid, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the name string to validate.
// Returns:
// - bool: True if the name is valid, false otherwise.
func ValidateNameField(fieldLevel validator.FieldLevel) bool {
	return nameRegex.MatchString(fieldLevel.Field().String())
}

// ValidateNameString checks if a provided name string (e.g., first name or last name)
// matches the validation rules.
// It returns a ValidationError if the name is invalid, nil otherwise.
// Parameters:
// - name: The name string to validate.
// Returns:
// - *ValidationError: A ValidationError if the name is invalid, nil otherwise.
func ValidateNameString(name string) *models.ValidationError {
	// Check if the name is empty.
	if name == "" {
		validationError := models.NewValidationError("name", name, "required", "Name is required.")
		return validationError

	}
	// Check if the name exceeds 70 characters.
	if len(name) > 70 {
		validationError := models.NewValidationError("name", name, "max", "Name cannot exceed 70 characters.")
		return validationError
	}
	// Check if the name format is valid.
	if !nameRegex.MatchString(name) {
		validationError := models.NewValidationError("name", name, "regex", "Invalid name format.")
		return validationError
	}
	return nil
}

// ValidateDateOfBirthField checks if a provided date of birth is valid.
// It ensures the date of birth is in the past and within the age range of 16 to 120.
// It returns true if valid, false otherwise.
// Parameters:
// - fieldLevel: The validator.FieldLevel containing the date of birth to validate.
// Returns:
// - bool: True if the date of birth is valid, false otherwise.
func ValidateDateOfBirthField(fieldLevel validator.FieldLevel) bool {
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

// ValidateDateOfBirthTime checks if a provided date of birth is valid.
// It ensures the date of birth is in the past and within the age range of 16 to 120.
// It returns a ValidationError if the date of birth is invalid, nil otherwise.
// Parameters:
// - dateOfBirth: The date of birth to validate.
// Returns:
// - *ValidationError: A ValidationError if the date of birth is invalid, nil otherwise.
func ValidateDateOfBirthTime(dateOfBirth time.Time) *models.ValidationError {
	now := time.Now()

	// Check if the date of birth is in the future.
	if dateOfBirth.After(now) {
		validationError := models.NewValidationError("dateofbirth", dateOfBirth, "past", "Date of birth cannot be in the future.")
		return validationError
	}
	// Calculate the minimum(16) and maximum(120) date of birth.
	minAge := now.AddDate(-16, 0, 0)
	maxAge := now.AddDate(-120, 0, 0)

	// Check if the date of birth is after the maximum age.
	if dateOfBirth.Before(maxAge) {
		validationError := models.NewValidationError("dateofbirth", dateOfBirth, "max", "Date of birth cannot be older than 120 years.")
		return validationError
	}

	// Check if the date of birth is before the minimum age.
	if dateOfBirth.After(minAge) {
		validationError := models.NewValidationError("dateofbirth", dateOfBirth, "min", "Date of birth cannot be younger than 16 years.")
		return validationError
	}

	return nil
}

// RegisterDomainValidators registers the custom validation functions for the User struct.
func RegisterDomainValidators(validator *validation.Validator) {
	validator.RegisterCustomValidation("username", ValidateUsernameField)
	validator.RegisterCustomValidation("password", ValidatePasswordField)
	validator.RegisterCustomValidation("name", ValidateNameField)
	validator.RegisterCustomValidation("dateofbirth", ValidateDateOfBirthField)
	validator.RegisterCustomValidation("email", ValidateEmailField)
}
