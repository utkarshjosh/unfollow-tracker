package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/utkarsh/unfollow-tracker/internal/domain"
)

type SnapshotRepository interface {
	Create(ctx context.Context, snapshot *domain.Snapshot) error
	FindLatestByAccountAndChunk(ctx context.Context, accountID uuid.UUID, chunkID int) (*domain.Snapshot, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error)
	DeleteOldSnapshots(ctx context.Context, accountID uuid.UUID, chunkID int, keepCount int) error
}

type PostgresSnapshotRepository struct {
	db *sql.DB
}

func NewPostgresSnapshotRepository(db *sql.DB) *PostgresSnapshotRepository {
	return &PostgresSnapshotRepository{db: db}
}

func (r *PostgresSnapshotRepository) Create(ctx context.Context, snapshot *domain.Snapshot) error {
	query := `
		INSERT INTO snapshots (id, account_id, chunk_id, follower_hashes, scanned_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		snapshot.ID,
		snapshot.AccountID,
		snapshot.ChunkID,
		pq.Array(snapshot.FollowerHashes),
		snapshot.ScannedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresSnapshotRepository) FindLatestByAccountAndChunk(ctx context.Context, accountID uuid.UUID, chunkID int) (*domain.Snapshot, error) {
	query := `
		SELECT id, account_id, chunk_id, follower_hashes, scanned_at
		FROM snapshots
		WHERE account_id = $1 AND chunk_id = $2
		ORDER BY scanned_at DESC
		LIMIT 1
	`

	snapshot := &domain.Snapshot{}
	err := r.db.QueryRowContext(ctx, query, accountID, chunkID).Scan(
		&snapshot.ID,
		&snapshot.AccountID,
		&snapshot.ChunkID,
		pq.Array(&snapshot.FollowerHashes),
		&snapshot.ScannedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // No previous snapshot is not an error
	}

	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (r *PostgresSnapshotRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error) {
	query := `
		SELECT id, account_id, chunk_id, follower_hashes, scanned_at
		FROM snapshots
		WHERE id = $1
	`

	snapshot := &domain.Snapshot{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&snapshot.ID,
		&snapshot.AccountID,
		&snapshot.ChunkID,
		pq.Array(&snapshot.FollowerHashes),
		&snapshot.ScannedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (r *PostgresSnapshotRepository) DeleteOldSnapshots(ctx context.Context, accountID uuid.UUID, chunkID int, keepCount int) error {
	query := `
		DELETE FROM snapshots
		WHERE id IN (
			SELECT id
			FROM snapshots
			WHERE account_id = $1 AND chunk_id = $2
			ORDER BY scanned_at DESC
			OFFSET $3
		)
	`

	_, err := r.db.ExecContext(ctx, query, accountID, chunkID, keepCount)
	if err != nil {
		return err
	}

	return nil
}
