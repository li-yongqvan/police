CREATE TABLE schema_forum.comments (
    id SERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES schema_forum.posts(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_comments_post_id ON schema_forum.comments(post_id);

CREATE TABLE schema_forum.likes (
    id SERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES schema_forum.posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(post_id, user_id)
);

CREATE TABLE schema_forum.collections (
    id SERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES schema_forum.posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(post_id, user_id)
);
