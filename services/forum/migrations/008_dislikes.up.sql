CREATE TABLE schema_forum.dislikes (
    id SERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES schema_forum.posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(post_id, user_id)
);

ALTER TABLE schema_forum.posts ADD COLUMN dislike_count INT DEFAULT 0;
