package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/api/handlers"
	"github.com/utkarsh/unfollow-tracker/internal/api/middleware"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
)

// ListUnfollows returns detected unfollows for the user's accounts
func (s *Server) ListUnfollows(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handlers.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Query params
	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		handlers.Error(w, http.StatusBadRequest, "account_id is required")
		return
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		handlers.Error(w, http.StatusBadRequest, "invalid account_id")
		return
	}

	// Verify ownership
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

	// Parse pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get unfollows
	unfollows, err := s.unfollowSvc.GetUnfollows(r.Context(), accountID, limit, offset)
	if err != nil {
		handlers.Error(w, http.StatusInternalServerError, "failed to get unfollows")
		return
	}

	// Convert to response format (date only, not exact timestamp)
	unfollowResponses := make([]handlers.UnfollowResponse, len(unfollows))
	for i, unfollow := range unfollows {
		unfollowResponses[i] = handlers.UnfollowResponse{
			ID:           unfollow.ID.String(),
			AccountID:    unfollow.AccountID.String(),
			DetectedDate: unfollow.DetectedAt.Format("2006-01-02"), // Date only
		}
	}

	handlers.Success(w, map[string]interface{}{
		"unfollows": unfollowResponses,
		"total":     len(unfollowResponses),
		"limit":     limit,
		"offset":    offset,
	})
}

// GetUnfollowSummary returns aggregated unfollow data
func (s *Server) GetUnfollowSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handlers.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Query params
	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		handlers.Error(w, http.StatusBadRequest, "account_id is required")
		return
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		handlers.Error(w, http.StatusBadRequest, "invalid account_id")
		return
	}

	// Verify ownership
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

	// Parse period
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "week"
	}

	// Validate period
	if period != "day" && period != "week" && period != "month" {
		handlers.Error(w, http.StatusBadRequest, "period must be day, week, or month")
		return
	}

	// Calculate since timestamp based on period
	var since time.Time
	now := time.Now()
	switch period {
	case "day":
		since = now.AddDate(0, 0, -1)
	case "week":
		since = now.AddDate(0, 0, -7)
	case "month":
		since = now.AddDate(0, -1, 0)
	}

	// Get summary
	summary, err := s.unfollowSvc.GetSummary(r.Context(), accountID, since)
	if err != nil {
		handlers.Error(w, http.StatusInternalServerError, "failed to get unfollow summary")
		return
	}

	// Convert to response format
	summaryResponse := handlers.UnfollowSummaryResponse{
		AccountID:   summary.AccountID.String(),
		Username:    summary.Username,
		Period:      summary.Period,
		Count:       summary.Count,
		TrendChange: summary.TrendChange,
	}

	// Calculate health score based on unfollow count
	healthScore := 100.0
	if summary.Count > 0 {
		// Simple formula: reduce health by 10 points per unfollow, min 0
		healthScore = 100.0 - float64(summary.Count)*10.0
		if healthScore < 0 {
			healthScore = 0
		}
	}

	trend := "stable"
	if summary.TrendChange > 10 {
		trend = "worsening"
	} else if summary.TrendChange < -10 {
		trend = "improving"
	}

	handlers.Success(w, map[string]interface{}{
		"period":   period,
		"summary":  summaryResponse,
		"overall_health": map[string]interface{}{
			"score":   healthScore,
			"trend":   trend,
			"message": getHealthMessage(healthScore),
		},
	})
}

func getHealthMessage(score float64) string {
	if score >= 90 {
		return "Your audience is healthy!"
	} else if score >= 70 {
		return "Your audience is doing well"
	} else if score >= 50 {
		return "Some audience changes detected"
	} else if score >= 30 {
		return "Noticeable audience changes"
	}
	return "Significant audience changes"
}
