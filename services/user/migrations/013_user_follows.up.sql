-- Migration 013: user_follows table for social following feature
-- Adds a unidirectional follow relationship between users.

CREATE TABLE IF NOT EXISTS schema_auth.user_follows (
    id BIGSERIAL PRIMARY KEY,
    follower_id INTEGER NOT NULL REFERENCES schema_auth.users(id) ON DELETE CASCADE,
    followee_id INTEGER NOT NULL REFERENCES schema_auth.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(follower_id, followee_id),
    CHECK(follower_id <> followee_id)
);

CREATE INDEX IF NOT EXISTS idx_user_follows_follower ON schema_auth.user_follows(follower_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_follows_followee ON schema_auth.user_follows(followee_id, created_at DESC);
