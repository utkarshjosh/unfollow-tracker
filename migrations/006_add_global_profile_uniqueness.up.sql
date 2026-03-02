-- Prevent duplicate tracking entries for the same normalized profile across all users.
CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_platform_username_unique_global
ON accounts (platform, LOWER(username));
