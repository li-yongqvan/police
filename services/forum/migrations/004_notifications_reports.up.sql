CREATE TABLE IF NOT EXISTS schema_forum.notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    type VARCHAR(30) NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    related_post_id INT,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON schema_forum.notifications(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS schema_forum.post_reports (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    reporter_id INT NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_post_reports_post ON schema_forum.post_reports(post_id);
