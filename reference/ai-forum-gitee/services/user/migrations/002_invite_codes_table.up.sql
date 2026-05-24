CREATE TABLE schema_auth.invite_codes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(64) UNIQUE NOT NULL,
    created_by BIGINT DEFAULT NULL,
    used_by BIGINT DEFAULT NULL,
    used_at TIMESTAMPTZ DEFAULT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'unused',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_invite_codes_code ON schema_auth.invite_codes(code);
CREATE INDEX idx_invite_codes_status ON schema_auth.invite_codes(status);
