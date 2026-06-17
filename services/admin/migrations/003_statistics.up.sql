-- Create statistics_daily table for daily aggregation
CREATE TABLE IF NOT EXISTS schema_admin.statistics_daily (
    id SERIAL PRIMARY KEY,
    stat_date DATE NOT NULL UNIQUE,
    new_users INT DEFAULT 0,
    new_posts INT DEFAULT 0,
    new_comments INT DEFAULT 0,
    active_users INT DEFAULT 0,
    total_users INT DEFAULT 0,
    total_posts INT DEFAULT 0,
    total_comments INT DEFAULT 0,
    board_activity JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_statistics_date ON schema_admin.statistics_daily(stat_date);
