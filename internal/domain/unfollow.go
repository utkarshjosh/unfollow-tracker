package domain

import (
	"time"

	"github.com/google/uuid"
)

// Unfollow represents a detected unfollow event
type Unfollow struct {
	ID           uuid.UUID `json:"id" db:"id"`
	AccountID    uuid.UUID `json:"account_id" db:"account_id"`
	FollowerHash string    `json:"follower_hash" db:"follower_hash"`
	DetectedAt   time.Time `json:"detected_at" db:"detected_at"`
	Notified     bool      `json:"notified" db:"notified"`
}

// NewUnfollow creates a new unfollow record
func NewUnfollow(accountID uuid.UUID, followerHash string) *Unfollow {
	return &Unfollow{
		ID:           uuid.New(),
		AccountID:    accountID,
		FollowerHash: followerHash,
		DetectedAt:   time.Now(),
		Notified:     false,
	}
}

// UnfollowSummary provides aggregated unfollow data for display
// This is the ethical way to present data - aggregated, not individual
type UnfollowSummary struct {
	AccountID   uuid.UUID `json:"account_id"`
	Username    string    `json:"username"`
	Period      string    `json:"period"` // "day", "week", "month"
	Count       int       `json:"count"`
	TrendChange float64   `json:"trend_change"` // Percentage change from previous period
}

// UnfollowEvent represents a single unfollow for API response
// Note: We intentionally don't expose exact timestamps to reduce anxiety
type UnfollowEvent struct {
	ID           uuid.UUID `json:"id"`
	AccountID    uuid.UUID `json:"account_id"`
	DetectedDate string    `json:"detected_date"` // Only date, not exact time
}

// ToEvent converts an Unfollow to a display-safe UnfollowEvent
func (u *Unfollow) ToEvent() *UnfollowEvent {
	return &UnfollowEvent{
		ID:           u.ID,
		AccountID:    u.AccountID,
		DetectedDate: u.DetectedAt.Format("2006-01-02"), // Date only
	}
}

// UnfollowBatch represents multiple unfollows to be inserted
type UnfollowBatch struct {
	AccountID uuid.UUID
	Hashes    []string
}

// ToUnfollows converts a batch to individual unfollow records
func (b *UnfollowBatch) ToUnfollows() []*Unfollow {
	unfollows := make([]*Unfollow, len(b.Hashes))
	for i, hash := range b.Hashes {
		unfollows[i] = NewUnfollow(b.AccountID, hash)
	}
	return unfollows
}
