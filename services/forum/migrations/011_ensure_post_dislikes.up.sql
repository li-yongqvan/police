CREATE TABLE IF NOT EXISTS schema_forum.dislikes (
    id SERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES schema_forum.posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(post_id, user_id)
);

ALTER TABLE schema_forum.posts
    ADD COLUMN IF NOT EXISTS dislike_count INT NOT NULL DEFAULT 0;

UPDATE schema_forum.posts p
SET dislike_count = COALESCE(d.count, 0)
FROM (
    SELECT post_id, COUNT(*)::INT AS count
    FROM schema_forum.dislikes
    GROUP BY post_id
) d
WHERE p.id = d.post_id;

UPDATE schema_forum.posts
SET dislike_count = 0
WHERE dislike_count IS NULL;
