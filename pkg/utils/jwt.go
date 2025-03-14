package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key for JWT
var jwtKey = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret - CHANGE THIS IN PRODUCTION!
		return "your-secret-key-change-this-in-production"
	}
	return secret
}

// Claims structure for JWT
type Claims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	IsAdmin    bool   `json:"is_admin"`
	IsEmployee bool   `json:"is_employee"` // Add this line
	jwt.RegisteredClaims
}

// Generate a new JWT token
func GenerateToken(userID uint, email string, isAdmin bool, isEmployee bool) (string, error) {
	// Set expiration time - 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:     userID,
		Email:      email,
		IsAdmin:    isAdmin,
		IsEmployee: isEmployee, // Add this line
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Parse and validate JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
