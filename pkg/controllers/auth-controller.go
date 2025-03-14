package controllers

import (
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
