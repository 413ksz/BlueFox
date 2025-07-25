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

// VerifyJWTToken verifies the authenticity and validity of a JWT token.
// params:
// - tokenString: The JWT token string to verify.
// returns:
// - *models.MyClaims: The claims extracted from the token if verification is successful.
// - error: A generic apierrors.ERROR_CODE_UNAUTHORIZED error if the token is invalid or any other error occurs during verification.
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

		// Now, handle other specific JWT errors using errors.Is
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token expired", &err, nil)
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Malformed token", &err, nil)
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token not valid yet", &err, nil)
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token signature invalid", &err, nil)
		} else if errors.Is(err, jwt.ErrTokenInvalidAudience) {
			return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token audience invalid", &err, nil)
		} else if errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token issuer invalid", &err, nil)
		} else if errors.Is(err, jwt.ErrTokenInvalidId) {
			return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token ID invalid", &err, nil)
		}

		// If it's none of the above specific JWT errors, or our custom error,
		// then it's a general token verification failure or an unexpected error.
		return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token verification failed unexpectedly", &err, nil)
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(*models.MyClaims); ok {
		return claims, nil
	}

	// If token is not valid or claims type assertion fails, return an error
	return nil, models.NewCustomError(models.ERROR_CODE_UNAUTHORIZED, "Token verification failed", nil, nil)
}
