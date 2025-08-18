package jwt_token

import (
	// For formatted error messages
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWTToken generates a new JWT token for the given user details.
// params:
// - username: The username of the user.
// - id: The ID of the user.
// - profilePictureAssetId: The ID of the user's profile picture asset.
// returns:
// - string: The generated JWT token.
// - error: An error if the token generation fails.
func GenerateJWTToken(username string, id string, profilePictureAssetId string) (string, *models.CustomError) {
	// Define JWT expiration duration as a constant.
	const tokenDuration = 24 * time.Hour
	// Calculate the expiration time for the JWT token and the current time.
	expirationTime := time.Now().Add(tokenDuration)
	now := time.Now()

	// Retrieve JWT_SECRET_KEY more securely from an environment variable
	jwtSecretKeyStr := os.Getenv("JWT_SECRET_KEY")
	// If the environment variable is not set, return an error
	if jwtSecretKeyStr == "" {
		return "", models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "JWT_SECRET_KEY environment variable not found", nil, nil)
	}
	// Convert the JWT_SECRET_KEY to a byte slice for JWT signing
	jwtSecretKey := []byte(jwtSecretKeyStr)

	// Create a new MyClaims struct with user details to embed in the token.
	claims := &models.MyClaims{
		Username:              username,
		Id:                    id,
		ProfilePictureAssetId: profilePictureAssetId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "BlueFox",
			Subject:   username,
			Audience:  []string{"users"},
		},
	}
	// Create a new JWT token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key.
	tokenString, err := token.SignedString(jwtSecretKey)
	// If token generation fails, return an error.
	if err != nil {
		return "", models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Failed to generate JWT token", &err, nil)
	}

	// Return the generated JWT token.
	return tokenString, nil
}

// VerifyJWTToken verifies and returns the claims of a JWT token.
//
// 1. secretKey: Get the JWT secret key from the environment variable.
// 2. parseToken: Parse the JWT token using the secret key.
// 3. claims: Extract the claims from the parsed token.
// 4. return claims: Return the claims of the verified JWT token.
//
// params:
// - tokenString: The JWT token string to verify.
// returns:
// - *models.MyClaims: The claims of the verified JWT token.
// - *models.CustomError: An error if the token verification fails.
//
// Error conditions:
// - if the secret key is not found: JWT_SECRET_KEY environment variable not found.
// - if the token parsing fails: Failed to parse JWT token.
//
// Example usage:
//
//	claims, err := VerifyJWTToken(tokenString)
//	if err != nil {
//		return err
//	}
func VerifyJWTToken(tokenString string) (*models.MyClaims, *models.CustomError) {
	// Retrieve JWT_SECRET_KEY more securely from an environment variable
	jwtSecretKeyStr := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKeyStr == "" {
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "JWT_SECRET_KEY environment variable not found", nil, nil)
	}
	// Convert the JWT_SECRET_KEY to a byte slice for JWT parsing
	jwtSecretKey := []byte(jwtSecretKeyStr)

	// Define the expected audience for the token.
	expectedAudience := "users"

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &models.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, fmt.Sprintf("Unexpected signing method: %v\n", token.Header["alg"]), nil, nil)
		}
		return jwtSecretKey, nil
	}, jwt.WithAudience(expectedAudience))

	// If token parsing fails, return an error
	if err != nil {
		// First, check if the error is our custom error from the signing method validation
		var customErr *models.CustomError
		if errors.As(err, &customErr) {
			// It's your custom error from the callback, return it directly
			return nil, customErr
		}

		return nil, JwtErrorHelper(err)

	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(*models.MyClaims); ok {
		return claims, nil
	}

	// If token is not valid or claims type assertion fails, return an error
	return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token verification failed", nil, nil)
}

// JwtErrorHelper returns a custom error based on the given JWT error.
// params:
// - err: The JWT error to convert to a custom error.
// returns:
// - *models.CustomError: The custom error corresponding to the JWT error.
//
// Error conditions:
// - jwt.ErrTokenExpired: Token expired
// - jwt.ErrTokenMalformed: Malformed token
// - jwt.ErrTokenNotValidYet: Token not valid yet
// - jwt.ErrTokenSignatureInvalid: Token signature invalid
// - jwt.ErrTokenInvalidAudience: Token audience invalid
// - jwt.ErrTokenInvalidIssuer: Token issuer invalid
// - jwt.ErrTokenInvalidId: Token ID invalid
// - default: Generic error for any other JWT error
// Example usage:
//
//	if err != nil {
//		return jwtErrorHelper(err)
//	}
func JwtErrorHelper(err error) *models.CustomError {
	switch err {
	case jwt.ErrTokenExpired:
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token expired", &err, nil)
	case jwt.ErrTokenMalformed:
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Malformed token", &err, nil)
	case jwt.ErrTokenNotValidYet:
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token not valid yet", &err, nil)
	case jwt.ErrTokenSignatureInvalid:
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token signature invalid", &err, nil)
	case jwt.ErrTokenInvalidAudience:
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token audience invalid", &err, nil)
	case jwt.ErrTokenInvalidIssuer:
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token issuer invalid", &err, nil)
	case jwt.ErrTokenInvalidId:
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token ID invalid", &err, nil)
	default:
		return models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token verification failed unexpectedly", &err, nil)
	}
}
