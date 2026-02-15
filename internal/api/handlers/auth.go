package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Request/Response types

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string        `json:"token"`
	ExpiresAt int64         `json:"expires_at"`
	User      *UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Plan      string `json:"plan"`
	CreatedAt string `json:"created_at"`
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Email == "" {
		Error(w, http.StatusBadRequest, "email is required")
		return
	}
	if len(req.Password) < 8 {
		Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	// TODO: Implement actual registration logic
	// 1. Check if user exists
	// 2. Hash password
	// 3. Create user in database
	// 4. Generate JWT

	// Placeholder response
	Created(w, map[string]string{
		"message": "registration endpoint - implement with actual service",
	})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Email == "" || req.Password == "" {
		Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	// TODO: Implement actual login logic
	// 1. Find user by email
	// 2. Verify password
	// 3. Generate JWT

	// Placeholder - generate a demo token
	expiry := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "demo-user-id",
		"exp": expiry.Unix(),
		"iat": time.Now().Unix(),
	})

	// Sign with a demo secret (replace with config secret)
	tokenString, _ := token.SignedString([]byte("demo-secret"))

	Success(w, AuthResponse{
		Token:     tokenString,
		ExpiresAt: expiry.Unix(),
		User: &UserResponse{
			ID:        "demo-user-id",
			Email:     req.Email,
			Plan:      "free",
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	})
}

// GetCurrentUser returns the authenticated user's profile
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user from context and fetch from database

	Success(w, map[string]string{
		"message": "get current user endpoint - implement with actual service",
	})
}
