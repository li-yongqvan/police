CREATE TABLE schema_admin.sensitive_words (
    id SERIAL PRIMARY KEY,
    word VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'general',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(word)
);

CREATE INDEX idx_sensitive_words_word ON schema_admin.sensitive_words(word);

-- Seed initial sensitive words
INSERT INTO schema_admin.sensitive_words (word, category) VALUES
    ('赌博', 'spam'),
    ('代写', 'abuse'),
    ('作弊', 'abuse'),
    ('刷单', 'spam'),
    ('加微信', 'ads'),
    ('加QQ', 'ads'),
    ('兼职', 'spam'),
    ('刷粉', 'spam'),
    ('引流', 'ads')
ON CONFLICT (word) DO NOTHING;
