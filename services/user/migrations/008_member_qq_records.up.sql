CREATE TABLE schema_auth.member_qq_records (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    qq_number VARCHAR(20) NOT NULL,
    qq_nickname VARCHAR(100) DEFAULT '',
    submitted_at TIMESTAMPTZ,
    submitted_by VARCHAR(50) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_member_qq_records_qq_number ON schema_auth.member_qq_records(qq_number);
CREATE INDEX idx_member_qq_records_name ON schema_auth.member_qq_records(name);
