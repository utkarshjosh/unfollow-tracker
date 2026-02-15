package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStats(ctx context.Context, id uuid.UUID, followerCount int, scanStatus domain.ScanStatus, chunkCount int) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	FindAccountsNeedingScan(ctx context.Context, intervalHours int, limit int) ([]*domain.Account, error)
	UpdateLastScanned(ctx context.Context, id uuid.UUID) error
}

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (id, user_id, platform, username, follower_count, last_scan_at, scan_status, chunk_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		account.ID,
		account.UserID,
		account.Platform,
		account.Username,
		account.FollowerCount,
		account.LastScanAt,
		account.ScanStatus,
		account.ChunkCount,
		account.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresAccountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	query := `
		SELECT id, user_id, platform, username, follower_count, last_scan_at, scan_status, chunk_count, created_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*domain.Account{}
	for rows.Next() {
		account := &domain.Account{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Platform,
			&account.Username,
			&account.FollowerCount,
			&account.LastScanAt,
			&account.ScanStatus,
			&account.ChunkCount,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *PostgresAccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	query := `
		SELECT id, user_id, platform, username, follower_count, last_scan_at, scan_status, chunk_count, created_at
		FROM accounts
		WHERE id = $1
	`

	account := &domain.Account{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Platform,
		&account.Username,
		&account.FollowerCount,
		&account.LastScanAt,
		&account.ScanStatus,
		&account.ChunkCount,
		&account.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *PostgresAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM accounts
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}

func (r *PostgresAccountRepository) UpdateStats(ctx context.Context, id uuid.UUID, followerCount int, scanStatus domain.ScanStatus, chunkCount int) error {
	query := `
		UPDATE accounts
		SET follower_count = $1, scan_status = $2, chunk_count = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, followerCount, scanStatus, chunkCount, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}

func (r *PostgresAccountRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM accounts
		WHERE user_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresAccountRepository) FindAccountsNeedingScan(ctx context.Context, intervalHours int, limit int) ([]*domain.Account, error) {
	query := `
		SELECT id, user_id, platform, username, follower_count, last_scan_at, scan_status, chunk_count, created_at
		FROM accounts
		WHERE (
			last_scan_at IS NULL
			OR last_scan_at + ($1 || ' hours')::INTERVAL < NOW()
		)
		AND scan_status NOT IN ('scanning', 'rate_limited')
		ORDER BY last_scan_at ASC NULLS FIRST
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, intervalHours, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*domain.Account{}
	for rows.Next() {
		account := &domain.Account{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Platform,
			&account.Username,
			&account.FollowerCount,
			&account.LastScanAt,
			&account.ScanStatus,
			&account.ChunkCount,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *PostgresAccountRepository) UpdateLastScanned(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE accounts
		SET last_scan_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}
