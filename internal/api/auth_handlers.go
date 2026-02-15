package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/utkarsh/unfollow-tracker/internal/api/handlers"
	"github.com/utkarsh/unfollow-tracker/internal/api/middleware"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
)

// Register handles user registration
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var req handlers.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Email == "" {
		handlers.Error(w, http.StatusBadRequest, "email is required")
		return
	}
	if len(req.Password) < 8 {
		handlers.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	// Call service
	user, token, err := s.authSvc.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			handlers.Error(w, http.StatusConflict, "user already exists")
			return
		}
		handlers.Error(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	// Return response
	handlers.Success(w, handlers.AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(s.config.JWT.Expiry).Unix(),
		User: &handlers.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Plan:      string(user.Plan),
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	})
}

// Login handles user authentication
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req handlers.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Email == "" || req.Password == "" {
		handlers.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	// Call service
	user, token, err := s.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			handlers.Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		handlers.Error(w, http.StatusInternalServerError, "failed to login")
		return
	}

	// Return response
	handlers.Success(w, handlers.AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(s.config.JWT.Expiry).Unix(),
		User: &handlers.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Plan:      string(user.Plan),
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	})
}

// GetCurrentUser returns the authenticated user's profile
func (s *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handlers.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := s.authSvc.GetUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handlers.Error(w, http.StatusNotFound, "user not found")
			return
		}
		handlers.Error(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	handlers.Success(w, handlers.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Plan:      string(user.Plan),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}
