package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token operations
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new JWT service instance
func NewJWTService() *JWTService {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default-secret-key" // fallback for development
	}
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// Claims represents the JWT claims structure
type Claims struct {
	Email    string `json:"email"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token with the specified claims
func (s *JWTService) GenerateToken(email string, userName string) (string, error) {
	// Set expiration to 30 days from now
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	claims := &Claims{
		Email:    email,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
