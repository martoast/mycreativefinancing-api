package middleware

import (
	"api/pkg/utils"
	"context"
	"net/http"
	"strings"
)

// Key for user context
type contextKey string

const UserContextKey contextKey = "user"

// User info to store in context
type UserContext struct {
	ID         uint
	Email      string
	IsAdmin    bool
	IsEmployee bool // Add this line
}

// AuthMiddleware checks for a valid JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, UserContext{
			ID:         claims.UserID,
			Email:      claims.Email,
			IsAdmin:    claims.IsAdmin,
			IsEmployee: claims.IsEmployee, // Add this line
		})

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extracts user info from context
func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKey).(UserContext)
	if !ok {
		return nil, false
	}
	return &user, true
}
