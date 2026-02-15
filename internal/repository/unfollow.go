package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
)

type UnfollowRepository interface {
	Create(ctx context.Context, unfollow *domain.Unfollow) error
	BatchCreate(ctx context.Context, unfollows []*domain.Unfollow) error
	FindByAccount(ctx context.Context, accountID uuid.UUID, limit int, offset int) ([]*domain.Unfollow, error)
	GetSummary(ctx context.Context, accountID uuid.UUID, since time.Time) (*domain.UnfollowSummary, error)
	CountByAccountSince(ctx context.Context, accountID uuid.UUID, since time.Time) (int, error)
}

type PostgresUnfollowRepository struct {
	db *sql.DB
}

func NewPostgresUnfollowRepository(db *sql.DB) *PostgresUnfollowRepository {
	return &PostgresUnfollowRepository{db: db}
}

func (r *PostgresUnfollowRepository) Create(ctx context.Context, unfollow *domain.Unfollow) error {
	query := `
		INSERT INTO unfollows (id, account_id, follower_hash, detected_at, notified)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		unfollow.ID,
		unfollow.AccountID,
		unfollow.FollowerHash,
		unfollow.DetectedAt,
		unfollow.Notified,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUnfollowRepository) BatchCreate(ctx context.Context, unfollows []*domain.Unfollow) error {
	if len(unfollows) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO unfollows (id, account_id, follower_hash, detected_at, notified)
		VALUES ($1, $2, $3, $4, $5)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, unfollow := range unfollows {
		_, err := stmt.ExecContext(
			ctx,
			unfollow.ID,
			unfollow.AccountID,
			unfollow.FollowerHash,
			unfollow.DetectedAt,
			unfollow.Notified,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PostgresUnfollowRepository) FindByAccount(ctx context.Context, accountID uuid.UUID, limit int, offset int) ([]*domain.Unfollow, error) {
	query := `
		SELECT id, account_id, follower_hash, detected_at, notified
		FROM unfollows
		WHERE account_id = $1
		ORDER BY detected_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	unfollows := []*domain.Unfollow{}
	for rows.Next() {
		unfollow := &domain.Unfollow{}
		err := rows.Scan(
			&unfollow.ID,
			&unfollow.AccountID,
			&unfollow.FollowerHash,
			&unfollow.DetectedAt,
			&unfollow.Notified,
		)
		if err != nil {
			return nil, err
		}
		unfollows = append(unfollows, unfollow)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return unfollows, nil
}

func (r *PostgresUnfollowRepository) GetSummary(ctx context.Context, accountID uuid.UUID, since time.Time) (*domain.UnfollowSummary, error) {
	query := `
		SELECT
			COUNT(*) as total_count
		FROM unfollows
		WHERE account_id = $1 AND detected_at >= $2
	`

	summary := &domain.UnfollowSummary{
		AccountID: accountID,
	}

	err := r.db.QueryRowContext(ctx, query, accountID, since).Scan(&summary.Count)
	if errors.Is(err, sql.ErrNoRows) {
		summary.Count = 0
		return summary, nil
	}

	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (r *PostgresUnfollowRepository) CountByAccountSince(ctx context.Context, accountID uuid.UUID, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM unfollows
		WHERE account_id = $1 AND detected_at >= $2
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, accountID, since).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
