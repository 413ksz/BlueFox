// Test file for user validation package.
// To run: cd backEnd/pkg/validation && go test -v.
// This file is part of the 'validation_test' package,
// which is separate from the 'validation' package.
package user_validation_test

import (
	"testing"

	validation "github.com/413ksz/BlueFox/backEnd/internal/util/validation/user_validation"
)

// TestValidateUsername tests the ValidateUsername function.
func TestValidateUsername(t *testing.T) {
	// Define test cases.
	// Each test case contains a name, username string, and expected result.
	// The expected result is true if the username is valid, false otherwise.
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			name:     "Valid: Alphanumeric",
			username: "user123",
			want:     true,
		},
		{
			name:     "Valid: With underscore",
			username: "my_user_name",
			want:     true,
		},
		{
			name:     "Valid: With hyphen",
			username: "another-user",
			want:     true,
		},
		{
			name:     "Valid: Minimum length (3 chars)",
			username: "abc",
			want:     true,
		},
		{
			name:     "Valid: Maximum length (20 chars)",
			username: "abcdefghijklmnopqrsu", // 20 chars
			want:     true,
		},
		{
			name:     "Invalid: Too short (2 chars)",
			username: "ab",
			want:     false,
		},
		{
			name:     "Invalid: Too long (21 chars)",
			username: "abcdefghijklmnopqrsut", // 21 chars
			want:     false,
		},
		{
			name:     "Invalid: Starts with underscore",
			username: "_username",
			want:     false,
		},
		{
			name:     "Invalid: Starts with hyphen",
			username: "-username",
			want:     false,
		},
		{
			name:     "Invalid: Ends with underscore",
			username: "username_",
			want:     false,
		},
		{
			name:     "Invalid: Ends with hyphen",
			username: "username-",
			want:     false,
		},
		{
			name:     "Invalid: Empty string",
			username: "",
			want:     false,
		},
		{
			name:     "Invalid: Contains space",
			username: "user name",
			want:     false,
		},
		{
			name:     "Invalid: Contains special character",
			username: "user!name",
			want:     false,
		},
		{
			name:     "Valid: Only numbers",
			username: "12345",
			want:     true,
		},
		{
			name:     "Valid: Only letters",
			username: "myusername",
			want:     true,
		},
		{
			name:     "Invalid: Single character", // Per pattern, requires at least 3 chars or complex middle part
			username: "a",
			want:     false,
		},
		{
			name:     "Invalid: Two characters", // Per pattern, requires at least 3 chars
			username: "ab",
			want:     false,
		},
		{
			name:     "Invalid: Contains multiple hyphens",
			username: "user--name",
			want:     true,
		},
		{
			name:     "Invalid: Contains multiple underscores",
			username: "user__name",
			want:     true,
		},
		{
			name:     "Invalid: Starts with number and has special char in middle",
			username: "1user!name",
			want:     false,
		},
		{
			name:     "Valid: Exactly 3 characters",
			username: "usr",
			want:     true,
		},
		{
			name:     "Valid: Exactly 20 characters",
			username: "username_test_123456",
			want:     true,
		},
		{
			name:     "Invalid: Contains only numbers (less than 3)",
			username: "12",
			want:     false,
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validation.ValidateUsername(tt.username); got != tt.want { // Call from the imported package
				t.Errorf("ValidateUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}

// TestValidateEmail tests the ValidateEmail function.
func TestValidateEmail(t *testing.T) {
	// Define test cases.
	// Each test case contains a name, an email address, and the expected result.
	// The expected result is true if the email address is valid, false otherwise.
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "Valid: Standard email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "Valid: With subdomain",
			email: "user@sub.domain.co.uk",
			want:  true,
		},
		{
			name:  "Valid: With numbers and special chars in local part",
			email: "john.doe+123@example.org",
			want:  true,
		},
		{
			name:  "Valid: Short domain",
			email: "test@a.co",
			want:  true,
		},
		{
			name:  "Invalid: No at symbol",
			email: "testexample.com",
			want:  false,
		},
		{
			name:  "Invalid: No domain",
			email: "test@",
			want:  false,
		},
		{
			name:  "Invalid: No local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "Invalid: No top-level domain",
			email: "test@example",
			want:  false,
		},
		{
			name:  "Invalid: Domain with invalid characters",
			email: "test@exa_mple.com",
			want:  false,
		},
		{
			name:  "Invalid: Empty string",
			email: "",
			want:  false,
		},
		{
			name:  "Invalid: Contains space",
			email: "test @example.com",
			want:  false,
		},
		{
			name:  "Valid: Email with single character local part",
			email: "a@b.com",
			want:  true,
		},
		{
			name:  "Valid: Email with long TLD",
			email: "user@example.travel",
			want:  true,
		},
		{
			name:  "Invalid: Missing dot in domain",
			email: "test@examplecom",
			want:  false,
		},
		{
			name:  "Invalid: Multiple '@' symbols",
			email: "test@ex@ample.com",
			want:  false,
		},
		{
			name:  "Invalid: Invalid characters in local part (e.g., space)",
			email: "first last@example.com",
			want:  false,
		},
		{
			name:  "Invalid: Leading dot in local part",
			email: ".test@example.com",
			want:  false,
		},
		{
			name:  "Invalid: Trailing dot in local part",
			email: "test.@example.com",
			want:  false,
		},
		{
			name:  "Invalid: Double dot in local part",
			email: "test..user@example.com",
			want:  false,
		},
		{
			name:  "Invalid: IP address as domain (not covered by this regex)",
			email: "test@[192.168.1.1]",
			want:  false,
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validation.ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

// TestValidatePassword tests the ValidatePassword function against all criteria.
func TestValidatePassword(t *testing.T) {
	// Define test cases
	// Each test case contains a name, a password, and the expected result
	// The expected result is true if the password is valid, false otherwise
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "Valid: All criteria met",
			password: "Passw0rd!",
			want:     true,
		},
		{
			name:     "Valid: All criteria met (different chars)",
			password: "MySuperStrongP@ssw0rd1",
			want:     true,
		},
		{
			name:     "Invalid: Too short (less than 8 chars)",
			password: "P@ssw1",
			want:     false,
		},
		{
			name:     "Invalid: No uppercase",
			password: "password1!",
			want:     false,
		},
		{
			name:     "Invalid: No lowercase",
			password: "PASSWORD1!",
			want:     false,
		},
		{
			name:     "Invalid: No number",
			password: "Password!!",
			want:     false,
		},
		{
			name:     "Invalid: No special character",
			password: "Password123",
			want:     false,
		},
		{
			name:     "Invalid: Empty string",
			password: "",
			want:     false,
		},
		{
			name:     "Invalid: Only letters",
			password: "Longpassword",
			want:     false,
		},
		{
			name:     "Invalid: Only numbers",
			password: "123456789",
			want:     false,
		},
		{
			name:     "Invalid: Only special characters",
			password: "!!!!!!!!",
			want:     false,
		},
		{
			name:     "Invalid: Missing two criteria (no upper, no special)",
			password: "password123",
			want:     false,
		},
		{
			name:     "Invalid: Missing three criteria (no upper, no lower, no special)",
			password: "12345678",
			want:     false,
		},
		{
			name:     "Valid: Password with all special characters",
			password: "P@ssw0rd!#$%",
			want:     true,
		},
		{
			name:     "Valid: Long password with mixed cases and numbers",
			password: "AnEvenLongerPassword123ABC!@#",
			want:     true,
		},
		{
			name:     "Invalid: Missing number but otherwise valid",
			password: "Password!!A",
			want:     false,
		},
		{
			name:     "Invalid: Missing special character but otherwise valid",
			password: "Password12A",
			want:     false,
		},
		{
			name:     "Invalid: Missing uppercase but otherwise valid",
			password: "password12!",
			want:     false,
		},
		{
			name:     "Invalid: Missing lowercase but otherwise valid",
			password: "PASSWORD12!",
			want:     false,
		},
		{
			name:     "Invalid: Exactly 8 characters, but missing number",
			password: "PasswOrd!",
			want:     false,
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validation.ValidatePassword(tt.password); got != tt.want {
				t.Errorf("ValidatePassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

// TestValidateName tests the ValidateName function.
func TestValidateName(t *testing.T) {
	// Define test cases
	// Each test case contains a name, a value, and the expected result
	// The expected result is true if the value is valid, false otherwise
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{
			name: "Valid: Standard first name",
			val:  "John",
			want: true,
		},
		{
			name: "Valid: First name with hyphen",
			val:  "Mary-Jane",
			want: true,
		},
		{
			name: "Valid: Name with apostrophe",
			val:  "O'Malley",
			want: true,
		},
		{
			name: "Valid: Name with space (e.g., middle name)",
			val:  "Ann Marie",
			want: true,
		},
		{
			name: "Valid: Minimum length (3 chars)",
			val:  "Ali",
			want: true,
		},
		{
			name: "Valid: Maximum length (20 chars)",
			val:  "abcdefghijklmnopqrsu", // 20 chars
			want: true,
		},
		{
			name: "Invalid: Too short (2 chars)",
			val:  "Jo",
			want: false,
		},
		{
			name: "Invalid: Too long (21 chars)",
			val:  "abcdefghijklmnopqrsut", // 21 chars
			want: false,
		},
		{
			name: "Invalid: Starts with hyphen",
			val:  "-John",
			want: false,
		},
		{
			name: "Invalid: Starts with apostrophe",
			val:  "'Malley",
			want: false,
		},
		{
			name: "Invalid: Ends with hyphen",
			val:  "John-",
			want: false,
		},
		{
			name: "Invalid: Ends with apostrophe",
			val:  "O'Malle'",
			want: false,
		},
		{
			name: "Invalid: Contains number",
			val:  "John123",
			want: false,
		},
		{
			name: "Invalid: Contains invalid special character",
			val:  "John!",
			want: false,
		},
		{
			name: "Invalid: Empty string",
			val:  "",
			want: false,
		},
		{
			name: "Invalid: Only spaces",
			val:  "   ",
			want: false,
		},
		{
			name: "Invalid: Starts with space",
			val:  " John",
			want: false,
		},
		{
			name: "Valid: Name with mixed case",
			val:  "Alice Wonderland",
			want: true,
		},
		{
			name: "Invalid: Name with leading space",
			val:  " John Doe",
			want: false,
		},
		{
			name: "Invalid: Name with trailing space",
			val:  "Jane Doe ",
			want: false,
		},
		{
			name: "Valid: Name with multiple internal spaces",
			val:  "Ann   Marie",
			want: true,
		},
		{
			name: "Invalid: Name with numbers",
			val:  "Test1Name",
			want: false,
		},
		{
			name: "Valid: Very long valid name (exactly 20 chars)",
			val:  "ABCDEFGHIJKLMNOPQRST",
			want: true,
		},
		{
			name: "Invalid: Contains unsupported special char",
			val:  "John.Doe",
			want: false,
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validation.ValidateName(tt.val); got != tt.want {
				t.Errorf("ValidateName(%q) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}
