-- 006_create_rate_limits.up.sql
-- Distributed rate limiting table

CREATE TABLE IF NOT EXISTS rate_limits (
    key TEXT NOT NULL,
    window TIMESTAMPTZ NOT NULL,
    count BIGINT NOT NULL DEFAULT 1,
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (key, window)
);

CREATE INDEX IF NOT EXISTS idx_rate_limits_expires ON rate_limits(expires_at);