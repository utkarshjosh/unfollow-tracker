package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
)

type stubAccountRepo struct{}

func (s *stubAccountRepo) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (s *stubAccountRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	return nil, nil
}

func (s *stubAccountRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, domain.ErrAccountNotFound
}

func (s *stubAccountRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *stubAccountRepo) UpdateStats(ctx context.Context, id uuid.UUID, followerCount int, scanStatus domain.ScanStatus, chunkCount int) error {
	return nil
}

func (s *stubAccountRepo) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	return 0, nil
}

func (s *stubAccountRepo) FindAccountsNeedingScan(ctx context.Context, intervalHours int, limit int) ([]*domain.Account, error) {
	return nil, nil
}

func (s *stubAccountRepo) UpdateLastScanned(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *stubAccountRepo) ExistsByPlatformAndUsername(ctx context.Context, platform domain.Platform, username string) (bool, error) {
	return false, nil
}

type stubUserRepo struct {
	findByID func(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

func (s *stubUserRepo) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (s *stubUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, domain.ErrUserNotFound
}

func (s *stubUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if s.findByID != nil {
		return s.findByID(ctx, id)
	}
	return nil, domain.ErrUserNotFound
}

func (s *stubUserRepo) UpdatePlan(ctx context.Context, userID uuid.UUID, plan domain.Plan) error {
	return nil
}

func TestCreateAccount_ReturnsUserNotFound_WhenUserProfileDoesNotExist(t *testing.T) {
	svc := NewAccountService(
		&stubAccountRepo{},
		&stubUserRepo{
			findByID: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
		},
	)

	_, err := svc.CreateAccount(context.Background(), uuid.New(), "instagram", "someuser")
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
