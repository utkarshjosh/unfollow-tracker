package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/api/handlers"
	"github.com/utkarsh/unfollow-tracker/internal/api/middleware"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
)

// ListAccounts returns all tracked accounts for the user
func (s *Server) ListAccounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handlers.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	accounts, err := s.accountSvc.GetAccounts(r.Context(), userID)
	if err != nil {
		handlers.Error(w, http.StatusInternalServerError, "failed to get accounts")
		return
	}

	// Convert to response format
	accountResponses := make([]handlers.AccountResponse, len(accounts))
	for i, acc := range accounts {
		var lastScanAt *string
		if acc.LastScanAt != nil {
			formatted := acc.LastScanAt.Format(time.RFC3339)
			lastScanAt = &formatted
		}

		accountResponses[i] = handlers.AccountResponse{
			ID:            acc.ID.String(),
			Username:      acc.Username,
			Platform:      string(acc.Platform),
			FollowerCount: acc.FollowerCount,
			ScanStatus:    string(acc.ScanStatus),
			LastScanAt:    lastScanAt,
			CreatedAt:     acc.CreatedAt.Format(time.RFC3339),
		}
	}

	handlers.Success(w, map[string]interface{}{
		"accounts": accountResponses,
		"total":    len(accountResponses),
	})
}

// CreateAccount adds a new Instagram account to track
func (s *Server) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handlers.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req handlers.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate username
	if req.Username == "" {
		handlers.Error(w, http.StatusBadRequest, "username is required")
		return
	}

	// Default to Instagram
	if req.Platform == "" {
		req.Platform = "instagram"
	}

	// Validate platform
	if req.Platform != "instagram" {
		handlers.Error(w, http.StatusBadRequest, "only instagram is supported currently")
		return
	}

	// Call service
	account, err := s.accountSvc.CreateAccount(r.Context(), userID, req.Platform, req.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handlers.Error(w, http.StatusBadRequest, "user profile not found")
			return
		}
		if errors.Is(err, domain.ErrAccountLimitReached) {
			handlers.Error(w, http.StatusBadRequest, "account limit reached for your plan")
			return
		}
		if errors.Is(err, domain.ErrAccountAlreadyExists) {
			handlers.Error(w, http.StatusBadRequest, "account is already being tracked")
			return
		}
		if errors.Is(err, domain.ErrInvalidUsername) {
			handlers.Error(w, http.StatusBadRequest, "invalid username")
			return
		}
		if errors.Is(err, domain.ErrInvalidPlatform) {
			handlers.Error(w, http.StatusBadRequest, "invalid platform")
			return
		}
		handlers.Error(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	var lastScanAt *string
	if account.LastScanAt != nil {
		formatted := account.LastScanAt.Format(time.RFC3339)
		lastScanAt = &formatted
	}

	handlers.Created(w, handlers.AccountResponse{
		ID:            account.ID.String(),
		Username:      account.Username,
		Platform:      string(account.Platform),
		FollowerCount: account.FollowerCount,
		ScanStatus:    string(account.ScanStatus),
		LastScanAt:    lastScanAt,
		CreatedAt:     account.CreatedAt.Format(time.RFC3339),
	})
}

// GetAccount returns a specific tracked account
func (s *Server) GetAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handlers.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	accountIDStr := chi.URLParam(r, "accountID")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		handlers.Error(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := s.accountSvc.GetAccount(r.Context(), accountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			handlers.Error(w, http.StatusNotFound, "account not found")
			return
		}
		handlers.Error(w, http.StatusInternalServerError, "failed to get account")
		return
	}

	// Verify ownership
	if account.UserID != userID {
		handlers.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	var lastScanAt *string
	if account.LastScanAt != nil {
		formatted := account.LastScanAt.Format(time.RFC3339)
		lastScanAt = &formatted
	}

	handlers.Success(w, handlers.AccountResponse{
		ID:            account.ID.String(),
		Username:      account.Username,
		Platform:      string(account.Platform),
		FollowerCount: account.FollowerCount,
		ScanStatus:    string(account.ScanStatus),
		LastScanAt:    lastScanAt,
		CreatedAt:     account.CreatedAt.Format(time.RFC3339),
	})
}

// DeleteAccount removes a tracked account
func (s *Server) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handlers.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	accountIDStr := chi.URLParam(r, "accountID")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		handlers.Error(w, http.StatusBadRequest, "invalid account id")
		return
	}

	// Call service
	if err := s.accountSvc.DeleteAccount(r.Context(), accountID, userID); err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			handlers.Error(w, http.StatusNotFound, "account not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			handlers.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		handlers.Error(w, http.StatusInternalServerError, "failed to delete account")
		return
	}

	handlers.NoContent(w)
}

// GetAccountStats returns statistics for a specific account
func (s *Server) GetAccountStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handlers.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	accountIDStr := chi.URLParam(r, "accountID")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		handlers.Error(w, http.StatusBadRequest, "invalid account id")
		return
	}

	// First verify ownership
	account, err := s.accountSvc.GetAccount(r.Context(), accountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			handlers.Error(w, http.StatusNotFound, "account not found")
			return
		}
		handlers.Error(w, http.StatusInternalServerError, "failed to get account")
		return
	}

	if account.UserID != userID {
		handlers.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	// Get stats
	stats, err := s.accountSvc.GetStats(r.Context(), accountID)
	if err != nil {
		handlers.Error(w, http.StatusInternalServerError, "failed to get account stats")
		return
	}

	handlers.Success(w, handlers.AccountStatsResponse{
		TotalUnfollows: stats.TotalUnfollows,
		Last24Hours:    stats.Last24Hours,
		Last7Days:      stats.Last7Days,
		Last30Days:     stats.Last30Days,
		HealthScore:    stats.HealthScore,
	})
}
