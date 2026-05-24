CREATE TABLE schema_forum.boards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    enabled BOOLEAN DEFAULT true,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE schema_forum.posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    author_id BIGINT NOT NULL,
    board_id BIGINT NOT NULL REFERENCES schema_forum.boards(id),
    status VARCHAR(20) NOT NULL DEFAULT 'published',
    is_pinned BOOLEAN DEFAULT false,
    is_featured BOOLEAN DEFAULT false,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_posts_board_id ON schema_forum.posts(board_id);
CREATE INDEX idx_posts_author_id ON schema_forum.posts(author_id);
CREATE INDEX idx_posts_status ON schema_forum.posts(status);
CREATE INDEX idx_posts_created_at ON schema_forum.posts(created_at DESC);

-- Seed three core boards
INSERT INTO schema_forum.boards (name, slug, description, sort_order) VALUES
    ('AI学习交流区', 'ai-learning', '分享学习心得、讨论AI技术、交流学习经验', 1),
    ('协会公告&活动区', 'announcements', '查看协会最新公告、活动通知、招新信息', 2),
    ('技术问答求助区', 'tech-help', '提出技术问题、寻求帮助、解答疑惑', 3);
