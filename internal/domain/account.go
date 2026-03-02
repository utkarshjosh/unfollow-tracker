package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Platform represents the social media platform
type Platform string

const (
	PlatformInstagram Platform = "instagram"
	// Future platforms can be added here
	// PlatformTwitter   Platform = "twitter"
	// PlatformTikTok    Platform = "tiktok"
)

// ScanStatus represents the current scanning state
type ScanStatus string

const (
	ScanStatusPending     ScanStatus = "pending"
	ScanStatusScanning    ScanStatus = "scanning"
	ScanStatusCompleted   ScanStatus = "completed"
	ScanStatusFailed      ScanStatus = "failed"
	ScanStatusRateLimited ScanStatus = "rate_limited"
)

// Account represents an Instagram account being tracked
type Account struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	Platform      Platform   `json:"platform" db:"platform"`
	Username      string     `json:"username" db:"username"`
	FollowerCount int        `json:"follower_count" db:"follower_count"`
	LastScanAt    *time.Time `json:"last_scan_at" db:"last_scan_at"`
	ScanStatus    ScanStatus `json:"scan_status" db:"scan_status"`
	ChunkCount    int        `json:"chunk_count" db:"chunk_count"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// NormalizeUsername converts user input into a canonical username format.
func NormalizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.TrimPrefix(username, "@")
	return strings.ToLower(username)
}

// NewAccount creates a new account to track
func NewAccount(userID uuid.UUID, username string, platform Platform) *Account {
	return &Account{
		ID:         uuid.New(),
		UserID:     userID,
		Platform:   platform,
		Username:   username,
		ScanStatus: ScanStatusPending,
		ChunkCount: 1,
		CreatedAt:  time.Now(),
	}
}

// CalculateChunkCount determines optimal chunk count based on follower count
// Each chunk should ideally contain ~1000 followers
func (a *Account) CalculateChunkCount() int {
	if a.FollowerCount <= 1000 {
		return 1
	}
	chunks := a.FollowerCount / 1000
	if a.FollowerCount%1000 > 0 {
		chunks++
	}
	// Cap at 100 chunks max (100k followers per chunk iteration)
	if chunks > 100 {
		return 100
	}
	return chunks
}

// NeedsScan returns true if the account is due for scanning
func (a *Account) NeedsScan(intervalHours int) bool {
	if a.LastScanAt == nil {
		return true
	}
	return time.Since(*a.LastScanAt) > time.Duration(intervalHours)*time.Hour
}

// AccountStats holds aggregated statistics for an account
type AccountStats struct {
	AccountID      uuid.UUID `json:"account_id"`
	TotalUnfollows int       `json:"total_unfollows"`
	Last24Hours    int       `json:"last_24_hours"`
	Last7Days      int       `json:"last_7_days"`
	Last30Days     int       `json:"last_30_days"`
	HealthScore    float64   `json:"health_score"` // 0-100, higher is better
}
