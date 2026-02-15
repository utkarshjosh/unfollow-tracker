package domain

import (
	"time"

	"github.com/google/uuid"
)

// Snapshot represents a point-in-time capture of followers in a chunk
type Snapshot struct {
	ID             uuid.UUID `json:"id" db:"id"`
	AccountID      uuid.UUID `json:"account_id" db:"account_id"`
	ChunkID        int       `json:"chunk_id" db:"chunk_id"`
	FollowerHashes []string  `json:"follower_hashes" db:"follower_hashes"`
	ScannedAt      time.Time `json:"scanned_at" db:"scanned_at"`
}

// NewSnapshot creates a new snapshot
func NewSnapshot(accountID uuid.UUID, chunkID int, hashes []string) *Snapshot {
	return &Snapshot{
		ID:             uuid.New(),
		AccountID:      accountID,
		ChunkID:        chunkID,
		FollowerHashes: hashes,
		ScannedAt:      time.Now(),
	}
}

// FollowerChunk represents a chunk of followers for processing
type FollowerChunk struct {
	AccountID uuid.UUID
	ChunkID   int
	Hashes    []string
}

// ChunkDiff represents the difference between two snapshots
type ChunkDiff struct {
	AccountID  uuid.UUID
	ChunkID    int
	Unfollowed []string // Hashes of followers who unfollowed
	NewFollows []string // Hashes of new followers
	ScannedAt  time.Time
}

// ComputeDiff calculates the difference between previous and current follower hashes
func ComputeDiff(accountID uuid.UUID, chunkID int, previous, current []string) *ChunkDiff {
	prevSet := make(map[string]bool)
	for _, h := range previous {
		prevSet[h] = true
	}

	currSet := make(map[string]bool)
	for _, h := range current {
		currSet[h] = true
	}

	var unfollowed, newFollows []string

	// Find unfollows: in previous but not in current
	for hash := range prevSet {
		if !currSet[hash] {
			unfollowed = append(unfollowed, hash)
		}
	}

	// Find new follows: in current but not in previous
	for hash := range currSet {
		if !prevSet[hash] {
			newFollows = append(newFollows, hash)
		}
	}

	return &ChunkDiff{
		AccountID:  accountID,
		ChunkID:    chunkID,
		Unfollowed: unfollowed,
		NewFollows: newFollows,
		ScannedAt:  time.Now(),
	}
}
