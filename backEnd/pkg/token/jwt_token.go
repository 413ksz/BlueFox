package jwt_token

import (
	// For formatted error messages
	"fmt"
	"os"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
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
func GenerateJWTToken(username string, id string, profilePictureAssetId string) (string, error) {
	// Define JWT expiration duration as a constant.
	const tokenDuration = 24 * time.Hour
	// Calculate the expiration time for the JWT token and the current time.
	expirationTime := time.Now().Add(tokenDuration)
	now := time.Now()

	// Retrieve JWT_SECRET_KEY more securely from an environment variable
	jwtSecretKeyStr := os.Getenv("JWT_SECRET_KEY")
	// If the environment variable is not set, return an error
	if jwtSecretKeyStr == "" {
		return "", apierrors.ERROR_CODE_ENVIREMENT_VARIABLE_NOT_FOUND
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
		return "", apierrors.ERROR_CODE_INTERNAL_SERVER
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
func VerifyJWTToken(tokenString string) (*models.MyClaims, error) {
	// Retrieve JWT_SECRET_KEY more securely from an environment variable
	jwtSecretKeyStr := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKeyStr == "" {
		fmt.Println("Error: JWT_SECRET_KEY environment variable not found.")
		return nil, apierrors.ERROR_CODE_INTERNAL_SERVER
	}
	// Convert the JWT_SECRET_KEY to a byte slice for JWT parsing
	jwtSecretKey := []byte(jwtSecretKeyStr)

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &models.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, apierrors.ERROR_CODE_UNAUTHORIZED
		}
		return jwtSecretKey, nil
	})

	// If token parsing fails, return an error
	if err != nil {
		fmt.Printf("Token verification failed: %v\n", err)
		return nil, apierrors.ERROR_CODE_UNAUTHORIZED
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(*models.MyClaims); ok && token.Valid {
		return claims, nil
	}

	// If token is not valid or claims type assertion fails, return an error
	fmt.Println("Token is not valid or claims type assertion failed.")
	return nil, apierrors.ERROR_CODE_UNAUTHORIZED
}
