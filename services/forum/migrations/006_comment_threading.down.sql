DROP INDEX IF EXISTS schema_forum.idx_comments_parent_id;
ALTER TABLE schema_forum.comments DROP COLUMN IF EXISTS parent_id;
ALTER TABLE schema_forum.comments DROP COLUMN IF EXISTS depth;
