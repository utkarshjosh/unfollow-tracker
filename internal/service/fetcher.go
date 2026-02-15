package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
	"github.com/utkarsh/unfollow-tracker/internal/repository"
)

type FetcherService struct {
	accountRepo  repository.AccountRepository
	snapshotRepo repository.SnapshotRepository
	unfollowRepo repository.UnfollowRepository
}

func NewFetcherService(
	accountRepo repository.AccountRepository,
	snapshotRepo repository.SnapshotRepository,
	unfollowRepo repository.UnfollowRepository,
) *FetcherService {
	return &FetcherService{
		accountRepo:  accountRepo,
		snapshotRepo: snapshotRepo,
		unfollowRepo: unfollowRepo,
	}
}

func (s *FetcherService) ProcessFetchJob(ctx context.Context, accountID uuid.UUID, chunkIndex int, newHashes []string) error {
	// Find account
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to find account: %w", err)
	}

	// Find latest snapshot for chunk
	previousSnapshot, err := s.snapshotRepo.FindLatestByAccountAndChunk(ctx, accountID, chunkIndex)
	if err != nil {
		return fmt.Errorf("failed to find previous snapshot: %w", err)
	}

	// If no previous snapshot (first scan)
	if previousSnapshot == nil {
		log.Printf("First scan for account %s, chunk %d. Creating initial snapshot with %d followers", account.Username, chunkIndex, len(newHashes))

		// Create first snapshot
		snapshot := domain.NewSnapshot(accountID, chunkIndex, newHashes)

		// Insert snapshot
		if err := s.snapshotRepo.Create(ctx, snapshot); err != nil {
			return fmt.Errorf("failed to create initial snapshot: %w", err)
		}

		// Update account scan time
		if err := s.accountRepo.UpdateLastScanned(ctx, accountID); err != nil {
			return fmt.Errorf("failed to update last scanned: %w", err)
		}

		return nil
	}

	// If previous snapshot exists, compute diff
	diff := domain.ComputeDiff(accountID, chunkIndex, previousSnapshot.FollowerHashes, newHashes)

	// Log stats
	unchangedCount := len(newHashes) - len(diff.NewFollows)
	log.Printf("Chunk %d diff for account %s: %d new, %d lost, %d unchanged",
		chunkIndex, account.Username, len(diff.NewFollows), len(diff.Unfollowed), unchangedCount)

	// Create unfollow events for lost followers
	if len(diff.Unfollowed) > 0 {
		var unfollows []*domain.Unfollow
		for _, hash := range diff.Unfollowed {
			unfollows = append(unfollows, domain.NewUnfollow(accountID, hash))
		}

		// Batch insert unfollows
		if err := s.unfollowRepo.BatchCreate(ctx, unfollows); err != nil {
			return fmt.Errorf("failed to create unfollows: %w", err)
		}

		log.Printf("Created %d unfollow records for account %s, chunk %d", len(unfollows), account.Username, chunkIndex)
	}

	// Create new snapshot
	newSnapshot := domain.NewSnapshot(accountID, chunkIndex, newHashes)
	if err := s.snapshotRepo.Create(ctx, newSnapshot); err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Update account stats
	if err := s.accountRepo.UpdateStats(ctx, accountID, len(newHashes), domain.ScanStatusCompleted, account.ChunkCount); err != nil {
		return fmt.Errorf("failed to update account stats: %w", err)
	}

	// Update last scanned time
	if err := s.accountRepo.UpdateLastScanned(ctx, accountID); err != nil {
		return fmt.Errorf("failed to update last scanned: %w", err)
	}

	// Delete old snapshots (keep last 3)
	if err := s.snapshotRepo.DeleteOldSnapshots(ctx, accountID, chunkIndex, 3); err != nil {
		log.Printf("Warning: failed to delete old snapshots for account %s, chunk %d: %v", account.Username, chunkIndex, err)
		// Don't fail the whole operation if cleanup fails
	}

	return nil
}
