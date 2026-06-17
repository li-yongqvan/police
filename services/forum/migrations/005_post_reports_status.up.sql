ALTER TABLE schema_forum.post_reports
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS admin_note TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS resolved_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_post_reports_status ON schema_forum.post_reports(status, created_at DESC);

UPDATE schema_forum.post_reports SET status = 'pending' WHERE status IS NULL OR status = '';
