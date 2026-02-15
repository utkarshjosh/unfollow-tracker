-- Tracked accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(50) DEFAULT 'instagram',
    username VARCHAR(255) NOT NULL,
    follower_count INTEGER DEFAULT 0,
    last_scan_at TIMESTAMPTZ,
    scan_status VARCHAR(50) DEFAULT 'pending',
    chunk_count INTEGER DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Ensure unique account per user per platform
    UNIQUE(user_id, platform, username)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_scan_status ON accounts(scan_status);
CREATE INDEX IF NOT EXISTS idx_accounts_last_scan ON accounts(last_scan_at);
CREATE INDEX IF NOT EXISTS idx_accounts_platform_username ON accounts(platform, username);
