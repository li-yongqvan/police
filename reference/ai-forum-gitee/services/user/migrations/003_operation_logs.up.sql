CREATE TABLE schema_auth.operation_logs (
    id SERIAL PRIMARY KEY,
    operator_id BIGINT NOT NULL,
    operator_username VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    target_id BIGINT,
    detail JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_operation_logs_operator ON schema_auth.operation_logs(operator_id);
CREATE INDEX idx_operation_logs_action ON schema_auth.operation_logs(action);
CREATE INDEX idx_operation_logs_created ON schema_auth.operation_logs(created_at DESC);
