package token_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/413ksz/BlueFox/backEnd/internal/apierrors"
	jwt_token "github.com/413ksz/BlueFox/backEnd/internal/util/token"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Test data constants
const (
	testSecretKey      = "super-secret-test-key-123"
	testUsername       = "testuser"
	testUserID         = "user-123"
	testProfilePicture = "asset-456"
	alternateSecretKey = "another-secret-key-456"
)

// setupEnv sets the JWT_SECRET_KEY environment variable.
func setupEnv(secretKey string) {
	os.Setenv("JWT_SECRET_KEY", secretKey)
}

// clearEnv unsets the JWT_SECRET_KEY environment variable.
func clearEnv() {
	os.Unsetenv("JWT_SECRET_KEY")
}

// TestGenerateJWTToken_Success tests successful token generation.
func TestGenerateJWTToken_Success(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	// Corrected line: 'jwt_token' (lowercase 'j') instead of 'JwtToken' (uppercase 'J')
	tokenString, err := jwt_token.GenerateJWTToken(testUsername, testUserID, testProfilePicture)

	// Assert no error occurred
	assert.NoError(t, err, "GenerateJWTToken should not return an error on success")
	// Assert token string is not empty
	assert.NotEmpty(t, tokenString, "Generated token string should not be empty")

	// Optionally, verify the token to ensure it's valid with the correct claims
	claims, err := jwt_token.VerifyJWTToken(tokenString)
	assert.NoError(t, err, "VerifyJWTToken should not return an error for a valid token")
	assert.NotNil(t, claims, "Claims should not be nil for a valid token")
	assert.Equal(t, testUsername, claims.Username, "Username in claims should match")
	assert.Equal(t, testUserID, claims.Id, "ID in claims should match")
	assert.Equal(t, testProfilePicture, claims.ProfilePictureAssetId, "ProfilePictureAssetId in claims should match")
	assert.Equal(t, "BlueFox", claims.Issuer, "Issuer in claims should match")
	assert.Contains(t, claims.Audience, "users", "Audience should contain 'users'")
}

// TestGenerateJWTToken_EnvironmentVariableNotFound tests token generation when JWT_SECRET_KEY is missing.
func TestGenerateJWTToken_EnvironmentVariableNotFound(t *testing.T) {
	clearEnv() // Ensure the environment variable is not set

	tokenString, err := jwt_token.GenerateJWTToken(testUsername, testUserID, testProfilePicture)

	// Assert an error occurred
	assert.Error(t, err, "GenerateJWTToken should return an error when JWT_SECRET_KEY is not set")
	// Assert the specific error code
	assert.Equal(t, apierrors.ERROR_CODE_ENVIREMENT_VARIABLE_NOT_FOUND, err, "Error should be ERROR_CODE_ENVIREMENT_VARIABLE_NOT_FOUND")
	// Assert token string is empty
	assert.Empty(t, tokenString, "Generated token string should be empty on error")
}

// TestGenerateJWTToken_EmptyInputs tests token generation with empty username, id, or profilePictureAssetId.
// The current GenerateJWTToken implementation doesn't explicitly return an error for empty strings
// but generates a token with empty values. This test confirms that behavior.
func TestGenerateJWTToken_EmptyInputs(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	// Test with empty username
	tokenString, err := jwt_token.GenerateJWTToken("", testUserID, testProfilePicture)
	assert.NoError(t, err, "GenerateJWTToken should not return an error with empty username")
	assert.NotEmpty(t, tokenString, "Generated token string should not be empty with empty username")
	claims, err := jwt_token.VerifyJWTToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, "", claims.Username, "Username in claims should be empty")

	// Test with empty user ID
	tokenString, err = jwt_token.GenerateJWTToken(testUsername, "", testProfilePicture)
	assert.NoError(t, err, "GenerateJWTToken should not return an error with empty ID")
	assert.NotEmpty(t, tokenString, "Generated token string should not be empty with empty ID")
	claims, err = jwt_token.VerifyJWTToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, "", claims.Id, "ID in claims should be empty")

	// Test with empty profile picture asset ID
	tokenString, err = jwt_token.GenerateJWTToken(testUsername, testUserID, "")
	assert.NoError(t, err, "GenerateJWTToken should not return an error with empty profile picture ID")
	assert.NotEmpty(t, tokenString, "Generated token string should not be empty with empty profile picture ID")
	claims, err = jwt_token.VerifyJWTToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, "", claims.ProfilePictureAssetId, "ProfilePictureAssetId in claims should be empty")
}

// TestGenerateJWTToken_LongTokenDuration verifies that a token generated with a standard long duration is valid.
func TestGenerateJWTToken_LongTokenDuration(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	tokenString, err := jwt_token.GenerateJWTToken(testUsername, testUserID, testProfilePicture)
	assert.NoError(t, err, "GenerateJWTToken should not return an error for long duration token")
	assert.NotEmpty(t, tokenString, "Generated token string should not be empty for long duration token")

	claims, err := jwt_token.VerifyJWTToken(tokenString)
	assert.NoError(t, err, "VerifyJWTToken should not return an error for a valid long duration token")
	assert.NotNil(t, claims, "Claims should not be nil for a valid long duration token")

	// Check if expiration is roughly 24 hours from now
	assert.InDelta(t, time.Now().Add(24*time.Hour).Unix(), claims.ExpiresAt.Unix(), float64(time.Minute.Seconds()), "Token expiration should be approximately 24 hours from now")
}

// TestVerifyJWTToken_Success tests successful token verification.
func TestVerifyJWTToken_Success(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	// Generate a valid token
	tokenString, err := jwt_token.GenerateJWTToken(testUsername, testUserID, testProfilePicture)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Verify the token
	claims, err := jwt_token.VerifyJWTToken(tokenString)

	// Assert no error occurred
	assert.NoError(t, err, "VerifyJWTToken should not return an error for a valid token")
	// Assert claims are not nil
	assert.NotNil(t, claims, "Claims should not be nil for a valid token")
	// Assert the extracted claims match the original data
	assert.Equal(t, testUsername, claims.Username, "Username in claims should match")
	assert.Equal(t, testUserID, claims.Id, "ID in claims should match")
	assert.Equal(t, testProfilePicture, claims.ProfilePictureAssetId, "ProfilePictureAssetId in claims should match")
}

// TestVerifyJWTToken_EmptyTokenString tests verification with an empty token string.
func TestVerifyJWTToken_EmptyTokenString(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	claims, err := jwt_token.VerifyJWTToken("")

	assert.Error(t, err, "VerifyJWTToken should return an error for an empty token string")
	assert.Equal(t, apierrors.ERROR_CODE_UNAUTHORIZED, err, "Error should be ERROR_CODE_UNAUTHORIZED for empty token string")
	assert.Nil(t, claims, "Claims should be nil for an empty token string")
}

// TestVerifyJWTToken_InvalidToken tests verification with a malformed or invalid token string.
func TestVerifyJWTToken_InvalidToken(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	invalidTokenString := "this.is.an.invalid.token" // A clearly malformed token

	claims, err := jwt_token.VerifyJWTToken(invalidTokenString)

	// Assert an error occurred
	assert.Error(t, err, "VerifyJWTToken should return an error for an invalid token")
	// Assert the specific error code
	assert.Equal(t, apierrors.ERROR_CODE_UNAUTHORIZED, err, "Error should be ERROR_CODE_UNAUTHORIZED for invalid token")
	// Assert claims are nil
	assert.Nil(t, claims, "Claims should be nil for an invalid token")
}

// TestVerifyJWTToken_ExpiredToken tests verification with an expired token.
func TestVerifyJWTToken_ExpiredToken(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	// Create a token with a very short expiration time (e.g., 100 milliseconds)
	shortDuration := 100 * time.Millisecond
	expirationTime := time.Now().Add(shortDuration)
	now := time.Now()

	claims := &models.MyClaims{
		Username:              testUsername,
		Id:                    testUserID,
		ProfilePictureAssetId: testProfilePicture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "BlueFox",
			Subject:   testUsername,
			Audience:  []string{"users"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecretKey))
	assert.NoError(t, err)

	// Wait for the token to expire
	time.Sleep(shortDuration + 50*time.Millisecond) // Add a small buffer

	// Attempt to verify the expired token
	verifiedClaims, err := jwt_token.VerifyJWTToken(tokenString)

	// Assert an error occurred
	assert.Error(t, err, "VerifyJWTToken should return an error for an expired token")
	// Assert the specific error code
	assert.Equal(t, apierrors.ERROR_CODE_UNAUTHORIZED, err, "Error should be ERROR_CODE_UNAUTHORIZED for expired token")
	// Assert claims are nil
	assert.Nil(t, verifiedClaims, "Claims should be nil for an expired token")
}

// TestVerifyJWTToken_IncorrectSecretKey tests verification with a token signed by a different key.
func TestVerifyJWTToken_IncorrectSecretKey(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	// Generate a token using the primary secret key
	tokenString, err := jwt_token.GenerateJWTToken(testUsername, testUserID, testProfilePicture)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Change the environment variable to a different secret key for verification
	setupEnv(alternateSecretKey)

	// Attempt to verify the token with the incorrect secret key
	claims, err := jwt_token.VerifyJWTToken(tokenString)

	// Assert an error occurred
	assert.Error(t, err, "VerifyJWTToken should return an error when using an incorrect secret key")
	// Assert the specific error code
	assert.Equal(t, apierrors.ERROR_CODE_UNAUTHORIZED, err, "Error should be ERROR_CODE_UNAUTHORIZED for incorrect secret key")
	// Assert claims are nil
	assert.Nil(t, claims, "Claims should be nil when using an incorrect secret key")
}

// TestVerifyJWTToken_EnvironmentVariableNotFound tests token verification when JWT_SECRET_KEY is missing.
func TestVerifyJWTToken_EnvironmentVariableNotFound(t *testing.T) {
	clearEnv() // Ensure the environment variable is not set

	// A dummy token string; the error should occur before parsing due to missing env var
	dummyToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIn0.signature"

	// Suppress the fmt.Println output during this specific test case, as it's tested here.
	// You might want to temporarily redirect os.Stderr or capture output for more robust testing.
	// For this example, we'll just let it print, as the error check is paramount.
	fmt.Println("Note: Expecting 'Error: JWT_SECRET_KEY environment variable not found.' message below from the function under test.")

	claims, err := jwt_token.VerifyJWTToken(dummyToken)

	// Assert an error occurred
	assert.Error(t, err, "VerifyJWTToken should return an error when JWT_SECRET_KEY is not set")
	// Assert the specific error code
	assert.Equal(t, apierrors.ERROR_CODE_INTERNAL_SERVER, err, "Error should be ERROR_CODE_INTERNAL_SERVER when JWT_SECRET_KEY is not set")
	// Assert claims are nil
	assert.Nil(t, claims, "Claims should be nil on error")
}

// TestVerifyJWTToken_AudienceMismatch tests verification with a token that has an incorrect audience.
func TestVerifyJWTToken_AudienceMismatch(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	// Manually create a token with a different audience
	claims := &models.MyClaims{
		Username:              testUsername,
		Id:                    testUserID,
		ProfilePictureAssetId: testProfilePicture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "BlueFox",
			Subject:   testUsername,
			Audience:  []string{"admins"}, // Incorrect audience
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecretKey))
	assert.NoError(t, err)

	verifiedClaims, err := jwt_token.VerifyJWTToken(tokenString)

	assert.Error(t, err, "VerifyJWTToken should return an error for audience mismatch")
	assert.Equal(t, apierrors.ERROR_CODE_UNAUTHORIZED, err, "Error should be ERROR_CODE_UNAUTHORIZED for audience mismatch")
	assert.Nil(t, verifiedClaims, "Claims should be nil for audience mismatch")
}

// TestVerifyJWTToken_InvalidSigningMethod tests verification with a token signed by an unsupported algorithm.
func TestVerifyJWTToken_InvalidSigningMethod(t *testing.T) {
	setupEnv(testSecretKey)
	defer clearEnv()

	// Manually create a token signed with a different algorithm (e.g., jwt.SigningMethodNone)
	// Note: In a real scenario, this would likely involve a token generated externally.
	claims := &models.MyClaims{
		Username:              testUsername,
		Id:                    testUserID,
		ProfilePictureAssetId: testProfilePicture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "BlueFox",
			Subject:   testUsername,
			Audience:  []string{"users"},
		},
	}
	// Use a different signing method
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType) // No key needed for None method
	assert.NoError(t, err)

	verifiedClaims, err := jwt_token.VerifyJWTToken(tokenString)

	assert.Error(t, err, "VerifyJWTToken should return an error for an invalid signing method")
	assert.Equal(t, apierrors.ERROR_CODE_UNAUTHORIZED, err, "Error should be ERROR_CODE_UNAUTHORIZED for invalid signing method")
	assert.Nil(t, verifiedClaims, "Claims should be nil for an invalid signing method")
}
