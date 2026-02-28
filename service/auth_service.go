package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles JWT token generation and validation
type AuthService struct {
	jwtSecret  []byte
	expiration time.Duration
}

// NewAuthService creates a new AuthService instance
func NewAuthService(secret string, expiration time.Duration) *AuthService {
	return &AuthService{
		jwtSecret:  []byte(secret),
		expiration: expiration,
	}
}

// GenerateToken creates a new JWT token for the given username
// The token is signed using HS256 and includes the username and expiration time
func (s *AuthService) GenerateToken(username string) (string, error) {
	if username == "" {
		return "", errors.New("username cannot be empty")
	}

	// Create claims with username and expiration
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(s.expiration).Unix(),
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return the token string
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken verifies a JWT token string and returns the parsed token
// Returns an error if the token is invalid, expired, or uses wrong signing method
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, errors.New("token cannot be empty")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
