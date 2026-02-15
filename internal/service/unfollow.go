package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
	"github.com/utkarsh/unfollow-tracker/internal/repository"
)

type UnfollowService struct {
	unfollowRepo repository.UnfollowRepository
	accountRepo  repository.AccountRepository
}

func NewUnfollowService(
	unfollowRepo repository.UnfollowRepository,
	accountRepo repository.AccountRepository,
) *UnfollowService {
	return &UnfollowService{
		unfollowRepo: unfollowRepo,
		accountRepo:  accountRepo,
	}
}

func (s *UnfollowService) GetUnfollows(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*domain.Unfollow, error) {
	// Verify account exists
	_, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify account: %w", err)
	}

	// Fetch unfollows
	unfollows, err := s.unfollowRepo.FindByAccount(ctx, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get unfollows: %w", err)
	}

	return unfollows, nil
}

func (s *UnfollowService) GetSummary(ctx context.Context, accountID uuid.UUID, since time.Time) (*domain.UnfollowSummary, error) {
	// Verify account exists
	_, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify account: %w", err)
	}

	// Get summary
	summary, err := s.unfollowRepo.GetSummary(ctx, accountID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary: %w", err)
	}

	return summary, nil
}
