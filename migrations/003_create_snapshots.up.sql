-- Snapshots table (chunked follower data)
CREATE TABLE IF NOT EXISTS snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    chunk_id INTEGER NOT NULL,
    follower_hashes TEXT[] NOT NULL, -- Array of hashed follower IDs
    scanned_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for chunk lookups
CREATE INDEX IF NOT EXISTS idx_snapshots_account_chunk ON snapshots(account_id, chunk_id);
CREATE INDEX IF NOT EXISTS idx_snapshots_scanned_at ON snapshots(scanned_at);

-- Partial index for latest snapshot per account/chunk
CREATE INDEX IF NOT EXISTS idx_snapshots_latest ON snapshots(account_id, chunk_id, scanned_at DESC);
