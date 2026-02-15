package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// HashFollowerID creates a privacy-safe hash of a follower ID
// This ensures we never store raw Instagram IDs
func HashFollowerID(followerID string, salt string) string {
	data := followerID + salt
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// HashFollowerIDs hashes a list of follower IDs
func HashFollowerIDs(followerIDs []string, salt string) []string {
	hashes := make([]string, len(followerIDs))
	for i, id := range followerIDs {
		hashes[i] = HashFollowerID(id, salt)
	}
	// Sort for consistent comparison
	sort.Strings(hashes)
	return hashes
}

// QuickHash creates a simple hash without salt (for internal use only)
func QuickHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:16]) // Shorter hash for internal keys
}
