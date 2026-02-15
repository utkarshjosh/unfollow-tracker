-- Unfollows table (detected unfollow events)
CREATE TABLE IF NOT EXISTS unfollows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    follower_hash VARCHAR(64) NOT NULL, -- SHA256 hash of follower ID
    detected_at TIMESTAMPTZ DEFAULT NOW(),
    notified BOOLEAN DEFAULT FALSE
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_unfollows_account_id ON unfollows(account_id);
CREATE INDEX IF NOT EXISTS idx_unfollows_detected_at ON unfollows(detected_at);
CREATE INDEX IF NOT EXISTS idx_unfollows_notified ON unfollows(notified) WHERE notified = FALSE;

-- Composite index for timeline queries
CREATE INDEX IF NOT EXISTS idx_unfollows_account_detected ON unfollows(account_id, detected_at DESC);
