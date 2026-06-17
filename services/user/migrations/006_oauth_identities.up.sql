CREATE TABLE schema_auth.oauth_identities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES schema_auth.users(id) ON DELETE CASCADE,
    provider VARCHAR(20) NOT NULL,
    provider_user_id VARCHAR(64) NOT NULL,
    union_id VARCHAR(64),
    raw_profile JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (provider, provider_user_id)
);

CREATE INDEX idx_oauth_identities_user_id ON schema_auth.oauth_identities(user_id);
CREATE INDEX idx_oauth_identities_provider ON schema_auth.oauth_identities(provider);
