CREATE TABLE schema_forum.attachments (
    id SERIAL PRIMARY KEY,
    post_id BIGINT DEFAULT NULL REFERENCES schema_forum.posts(id) ON DELETE SET NULL,
    comment_id BIGINT DEFAULT NULL REFERENCES schema_forum.comments(id) ON DELETE SET NULL,
    user_id BIGINT NOT NULL,
    filename VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_attachments_post_id ON schema_forum.attachments(post_id);
CREATE INDEX idx_attachments_user_id ON schema_forum.attachments(user_id);
