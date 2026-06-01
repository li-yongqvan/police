ALTER TABLE schema_forum.comments
    ADD COLUMN IF NOT EXISTS parent_id BIGINT REFERENCES schema_forum.comments(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS depth INT NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON schema_forum.comments(parent_id);
