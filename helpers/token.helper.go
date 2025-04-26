package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define a struct for the claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Function to generate a JWT token
func GenerateToken(userID string) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(7 * time.Hour)
	secretKey := os.Getenv("AUTH_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("AUTH_SECRET_KEY environment variable not set")
	}

	// Create the claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce_api",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
