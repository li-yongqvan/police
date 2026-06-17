DROP INDEX IF EXISTS schema_forum.idx_post_reports_status;
ALTER TABLE schema_forum.post_reports
    DROP COLUMN IF EXISTS resolved_at,
    DROP COLUMN IF EXISTS admin_note,
    DROP COLUMN IF EXISTS status;
