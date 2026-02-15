package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Request/Response types

type CreateAccountRequest struct {
	Username string `json:"username"`
	Platform string `json:"platform"` // defaults to "instagram"
}

type AccountResponse struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Platform      string  `json:"platform"`
	FollowerCount int     `json:"follower_count"`
	ScanStatus    string  `json:"scan_status"`
	LastScanAt    *string `json:"last_scan_at"`
	CreatedAt     string  `json:"created_at"`
}

type AccountStatsResponse struct {
	TotalUnfollows int     `json:"total_unfollows"`
	Last24Hours    int     `json:"last_24_hours"`
	Last7Days      int     `json:"last_7_days"`
	Last30Days     int     `json:"last_30_days"`
	HealthScore    float64 `json:"health_score"`
}

// ListAccounts returns all tracked accounts for the user
func ListAccounts(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user ID from context
	// TODO: Fetch accounts from database

	Success(w, map[string]interface{}{
		"accounts": []AccountResponse{},
		"total":    0,
	})
}

// CreateAccount adds a new Instagram account to track
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate username
	if req.Username == "" {
		Error(w, http.StatusBadRequest, "username is required")
		return
	}

	// Default to Instagram
	if req.Platform == "" {
		req.Platform = "instagram"
	}

	// Validate platform
	if req.Platform != "instagram" {
		Error(w, http.StatusBadRequest, "only instagram is supported currently")
		return
	}

	// TODO: Implement account creation
	// 1. Check user's plan limits
	// 2. Verify account is public
	// 3. Create account in database
	// 4. Queue initial scan

	Created(w, map[string]string{
		"message": "create account endpoint - implement with actual service",
	})
}

// GetAccount returns a specific tracked account
func GetAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	if accountID == "" {
		Error(w, http.StatusBadRequest, "account id is required")
		return
	}

	// TODO: Fetch account from database
	// TODO: Verify ownership

	Success(w, map[string]string{
		"account_id": accountID,
		"message":    "get account endpoint - implement with actual service",
	})
}

// DeleteAccount removes a tracked account
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	if accountID == "" {
		Error(w, http.StatusBadRequest, "account id is required")
		return
	}

	// TODO: Delete account and associated data
	// TODO: Verify ownership

	NoContent(w)
}

// GetAccountStats returns statistics for a specific account
func GetAccountStats(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	if accountID == "" {
		Error(w, http.StatusBadRequest, "account id is required")
		return
	}

	// TODO: Calculate stats from database

	Success(w, AccountStatsResponse{
		TotalUnfollows: 0,
		Last24Hours:    0,
		Last7Days:      0,
		Last30Days:     0,
		HealthScore:    100.0, // Perfect health when no unfollows
	})
}
