package handlers

import (
	"net/http"
	"strconv"
)

// Response types

type UnfollowResponse struct {
	ID           string `json:"id"`
	AccountID    string `json:"account_id"`
	DetectedDate string `json:"detected_date"` // Date only, ethical choice
}

type UnfollowSummaryResponse struct {
	AccountID   string  `json:"account_id"`
	Username    string  `json:"username"`
	Period      string  `json:"period"`
	Count       int     `json:"count"`
	TrendChange float64 `json:"trend_change"`
}

// ListUnfollows returns detected unfollows for the user's accounts
func ListUnfollows(w http.ResponseWriter, r *http.Request) {
	// Query params
	accountID := r.URL.Query().Get("account_id") // optional, filter by account
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

	// TODO: Fetch unfollows from database
	// TODO: Apply filters
	_ = accountID
	_ = limit
	_ = offset

	Success(w, map[string]interface{}{
		"unfollows": []UnfollowResponse{},
		"total":     0,
		"limit":     limit,
		"offset":    offset,
	})
}

// GetUnfollowSummary returns aggregated unfollow data
// This is the ethical way to present data - focusing on trends, not individuals
func GetUnfollowSummary(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period") // "day", "week", "month"
	if period == "" {
		period = "week"
	}

	// Validate period
	if period != "day" && period != "week" && period != "month" {
		Error(w, http.StatusBadRequest, "period must be day, week, or month")
		return
	}

	// TODO: Calculate summaries from database
	// Focus on:
	// - Aggregated counts (not individual unfollowers)
	// - Trend changes (is it getting better or worse?)
	// - Health scores

	Success(w, map[string]interface{}{
		"period":    period,
		"summaries": []UnfollowSummaryResponse{},
		"overall_health": map[string]interface{}{
			"score":   100.0,
			"trend":   "stable",
			"message": "Your audience is healthy!",
		},
	})
}
