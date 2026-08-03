-- Migration 013 down: Remove seed posts
DELETE FROM schema_forum.posts WHERE id NOT IN (
    SELECT id FROM schema_forum.posts ORDER BY id LIMIT 4
);
