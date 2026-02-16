package api

import (
	"errors"
	"net/http"
	"sort"
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

	var allUnfollows []handlers.UnfollowResponse

	// Case 1: account_id is provided - return unfollows for that specific account
	if accountIDStr != "" {
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

		// Get unfollows for this specific account
		unfollows, err := s.unfollowSvc.GetUnfollows(r.Context(), accountID, limit, offset)
		if err != nil {
			handlers.Error(w, http.StatusInternalServerError, "failed to get unfollows")
			return
		}

		// Convert to response format (date only, not exact timestamp)
		allUnfollows = make([]handlers.UnfollowResponse, len(unfollows))
		for i, unfollow := range unfollows {
			allUnfollows[i] = handlers.UnfollowResponse{
				ID:           unfollow.ID.String(),
				AccountID:    unfollow.AccountID.String(),
				DetectedDate: unfollow.DetectedAt.Format("2006-01-02"), // Date only
			}
		}

		handlers.Success(w, map[string]interface{}{
			"unfollows": allUnfollows,
			"total":     len(allUnfollows),
			"limit":     limit,
			"offset":    offset,
		})
		return
	}

	// Case 2: account_id is NOT provided - return unfollows for ALL user's accounts
	accounts, err := s.accountSvc.GetAccounts(r.Context(), userID)
	if err != nil {
		handlers.Error(w, http.StatusInternalServerError, "failed to get accounts")
		return
	}

	// Handle empty account list gracefully
	if len(accounts) == 0 {
		handlers.Success(w, map[string]interface{}{
			"unfollows": []handlers.UnfollowResponse{},
			"total":     0,
			"limit":     limit,
			"offset":    offset,
		})
		return
	}

	// Fetch unfollows for each account (without pagination, we'll aggregate and paginate later)
	// We fetch all to properly sort and paginate across accounts
	for _, account := range accounts {
		unfollows, err := s.unfollowSvc.GetUnfollows(r.Context(), account.ID, 1000, 0) // Fetch large batch
		if err != nil {
			// Log error but continue with other accounts
			continue
		}

		// Convert to response format
		for _, unfollow := range unfollows {
			allUnfollows = append(allUnfollows, handlers.UnfollowResponse{
				ID:           unfollow.ID.String(),
				AccountID:    unfollow.AccountID.String(),
				DetectedDate: unfollow.DetectedAt.Format("2006-01-02"), // Date only
			})
		}
	}

	// Sort by detected date descending (most recent first)
	sort.Slice(allUnfollows, func(i, j int) bool {
		return allUnfollows[i].DetectedDate > allUnfollows[j].DetectedDate
	})

	// Apply pagination to the aggregated results
	total := len(allUnfollows)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedUnfollows := allUnfollows[start:end]

	handlers.Success(w, map[string]interface{}{
		"unfollows": paginatedUnfollows,
		"total":     total,
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

	// Case 1: account_id is provided - return summary for that specific account
	if accountIDStr != "" {
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
			"period":  period,
			"summary": summaryResponse,
			"overall_health": map[string]interface{}{
				"score":   healthScore,
				"trend":   trend,
				"message": getHealthMessage(healthScore),
			},
		})
		return
	}

	// Case 2: account_id is NOT provided - return aggregated summary across ALL user's accounts
	accounts, err := s.accountSvc.GetAccounts(r.Context(), userID)
	if err != nil {
		handlers.Error(w, http.StatusInternalServerError, "failed to get accounts")
		return
	}

	// Handle empty account list gracefully
	if len(accounts) == 0 {
		handlers.Success(w, map[string]interface{}{
			"period":    period,
			"summaries": []handlers.UnfollowSummaryResponse{},
			"overall_health": map[string]interface{}{
				"score":   100.0,
				"trend":   "stable",
				"message": "No accounts connected",
			},
		})
		return
	}

	// Fetch summary for each account
	var summaries []handlers.UnfollowSummaryResponse
	totalUnfollows := 0
	totalTrendChange := 0.0

	for _, account := range accounts {
		summary, err := s.unfollowSvc.GetSummary(r.Context(), account.ID, since)
		if err != nil {
			// Log error but continue with other accounts
			continue
		}

		summaryResponse := handlers.UnfollowSummaryResponse{
			AccountID:   summary.AccountID.String(),
			Username:    summary.Username,
			Period:      summary.Period,
			Count:       summary.Count,
			TrendChange: summary.TrendChange,
		}

		summaries = append(summaries, summaryResponse)
		totalUnfollows += summary.Count
		totalTrendChange += summary.TrendChange
	}

	// Calculate overall health based on total unfollows across all accounts
	healthScore := 100.0
	if totalUnfollows > 0 {
		// Simple formula: reduce health by 10 points per unfollow, min 0
		healthScore = 100.0 - float64(totalUnfollows)*10.0
		if healthScore < 0 {
			healthScore = 0
		}
	}

	// Calculate average trend change
	avgTrendChange := 0.0
	if len(summaries) > 0 {
		avgTrendChange = totalTrendChange / float64(len(summaries))
	}

	trend := "stable"
	if avgTrendChange > 10 {
		trend = "worsening"
	} else if avgTrendChange < -10 {
		trend = "improving"
	}

	handlers.Success(w, map[string]interface{}{
		"period":    period,
		"summaries": summaries,
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
