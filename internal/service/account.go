package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
	"github.com/utkarsh/unfollow-tracker/internal/repository"
)

type AccountService struct {
	accountRepo repository.AccountRepository
	userRepo    repository.UserRepository
}

func NewAccountService(accountRepo repository.AccountRepository, userRepo repository.UserRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		userRepo:    userRepo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID uuid.UUID, platform, username string) (*domain.Account, error) {
	normalizedUsername := domain.NormalizeUsername(username)
	if normalizedUsername == "" {
		return nil, domain.ErrInvalidUsername
	}

	// Get user to check plan limits
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Count existing accounts
	count, err := s.accountRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count accounts: %w", err)
	}

	// Get plan limits
	limits := user.Plan.Limits()

	// Enforce limit
	if count >= limits.MaxAccounts {
		return nil, domain.ErrAccountLimitReached
	}

	// Validate platform (only Instagram for now)
	if platform != string(domain.PlatformInstagram) {
		return nil, domain.ErrInvalidPlatform
	}

	alreadyTracked, err := s.accountRepo.ExistsByPlatformAndUsername(ctx, domain.Platform(platform), normalizedUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing account: %w", err)
	}
	if alreadyTracked {
		return nil, domain.ErrAccountAlreadyExists
	}

	// Create account
	account := domain.NewAccount(userID, normalizedUsername, domain.Platform(platform))

	// Insert to database
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s *AccountService) GetAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	accounts, err := s.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}

func (s *AccountService) GetAccount(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, accountID, userID uuid.UUID) error {
	// Find account
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to find account: %w", err)
	}

	// Verify ownership
	if account.UserID != userID {
		return domain.ErrForbidden
	}

	// Delete account
	if err := s.accountRepo.Delete(ctx, accountID); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}

func (s *AccountService) GetStats(ctx context.Context, accountID uuid.UUID) (*domain.AccountStats, error) {
	// Find account
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Return stats (placeholder structure - actual stats would be computed)
	stats := &domain.AccountStats{
		AccountID: account.ID,
	}

	return stats, nil
}
