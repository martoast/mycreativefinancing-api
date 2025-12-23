package controllers

import (
	"api/pkg/middleware"
	"api/pkg/models"
	"api/pkg/utils"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

// Request structure for registration
type RegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsEmployee bool   `json:"is_employee"` // Add this line
}

// Request structure for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response structure for auth
type AuthResponse struct {
	Token string `json:"token"`
}

// Request structure for password change
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// Response structure for password change
type ChangePasswordResponse struct {
	Message string `json:"message"`
}

// Register a new user
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	utils.ParseBody(r, &req)

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Create new user
	user := &models.User{
		Email:      req.Email,
		IsEmployee: req.IsEmployee, // Add this line
	}

	// Hash the password
	if err := user.HashPassword(req.Password); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Save user to database
	createdUser, err := models.CreateUser(user)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			http.Error(w, "User with this email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(createdUser.ID, createdUser.Email, createdUser.IsAdmin, createdUser.IsEmployee)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

// Login user
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	utils.ParseBody(r, &req)

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.IsAdmin, user.IsEmployee)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

// ChangeAdminPassword allows an admin user to change their password
func ChangeAdminPassword(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by AuthMiddleware)
	userContext, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is an admin
	if !userContext.IsAdmin {
		http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
		return
	}

	// Parse request body
	var req ChangePasswordRequest
	utils.ParseBody(r, &req)

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "Current password and new password are required", http.StatusBadRequest)
		return
	}

	// Validate new password length
	if len(req.NewPassword) < 8 {
		http.Error(w, "New password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Fetch the full user record from database
	user, err := models.GetUserByID(userContext.ID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash new password
	if err := user.HashPassword(req.NewPassword); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Save to database
	if err := models.UpdateUserPassword(user); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChangePasswordResponse{
		Message: "Password changed successfully",
	})
}
